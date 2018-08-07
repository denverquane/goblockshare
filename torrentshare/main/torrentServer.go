package main

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/denverquane/goblockshare/common"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type torrFileSpecs struct {
	url           string
	layerByteSize int64
	name          string
}

var env common.EnvVars

var myAddress common.PersonalAddress
var myAddress2 common.PersonalAddress
var torrentPath string

var torrents []common.TorrentFile
var layers map[string]common.LayerFileMetadata

func main() {
	env = common.LoadEnvFromFile("torrentshare")

	if len(os.Args) > 1 {
		torrentPath = os.Args[1]
	}

	log.Fatal(run())
}

func run() error {
	scanner := bufio.NewScanner(os.Stdin)

	jobs := make(chan torrFileSpecs, 10)
	results := make(chan common.TorrentFile, 10)

	for w := 1; w < 2; w++ {
		go torrentWorker(w, jobs, results)
	}
	totalJobs := 0

	jobs <- torrFileSpecs{torrentPath, 1000, torrentPath}
	totalJobs++

	myAddress = common.GenerateNewPersonalAddress()
	myAddress2 = common.GenerateNewPersonalAddress()

	if scanner.Err() != nil {
		// handle error.
	}
	muxx := makeMuxRouter()

	for i := 0; i < totalJobs; i++ {
		file := <-results
		if file.Name != "" {
			trans := common.PublishTorrentTrans{file}
			origin := myAddress.ConvertToOriginInfo()
			btt := common.SignableTransaction{origin, trans, common.PUBLISH_TORRENT, nil, nil, ""}
			signed := btt.SignAndSetTxID(&myAddress.PrivateKey)
			log.Println("Gonna broadcast " + signed.TxID + " to blockchains")
			broadcastTransaction(env.BlockchainHost+":"+env.BlockchainPort+"/addTransaction", signed)
			registerTorrent(file)

			//origin2 := common.AddressToOriginInfo(myAddress2)
			//trans2 := common.TorrentRepTrans{signed.TxID,
			//	common.RepMessage{true, true, true}}
			//btt2 := common.SignableTransaction{origin2, trans2, "TORRENT_REP", nil, nil, ""}
			//signed2 := btt2.SignAndSetTxID(&myAddress2.PrivateKey)
			//log.Println("Gonna broadcast " + signed2.TxID + " to blockchains")
			//broadcastTransaction(env.BlockchainHost + ":" + env.BlockchainPort + "/addTransaction", signed2)
		}
	}

	go listenForFeedback()

	log.Println("Listening on ", env.TorrentPort)

	s := &http.Server{
		Addr:           ":" + env.TorrentPort,
		Handler:        muxx,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func getStdin(scanner *bufio.Scanner) string {
	for scanner.Scan() {
		if scanner.Text() != "" {
			return scanner.Text()
		}
	}
	return ""
}

//torrentWorker represents a worker that should process a file on the local filesystem and process it into an
//internal TorrentFile. This could potentially still be a bottleneck from the underlying I/O filesystem, but processing
//entirely separate torrentfile is inherently parallelizable
func torrentWorker(id int, jobs <-chan torrFileSpecs, results chan<- common.TorrentFile) {
	for job := range jobs {
		torr, err := common.MakeTorrentFileFromFile(job.layerByteSize, job.url, job.name)
		if err != nil {
			fmt.Println("Worker " + strconv.Itoa(id) + " has error: " + err.Error())
		} else {
			fmt.Println("Worker " + strconv.Itoa(id) + " finished " + job.name)
		}
		results <- torr
	}
}

func listenForInput() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter the IP and port of a node you wish to query for torrents and layers")
	fmt.Println("For example (don't include http://): localhost:7070")
	fmt.Println("Enter \"done\" or \"quit\" at any time to finish querying peers")

	ip := getStdin(scanner)

	for ip != "quit" && ip != "done" {
		fmt.Println("Query \"" + ip + "\" for torrents, or layers? [T/L]")
		choice := getStdin(scanner)
		var end string
		if choice == "done" || choice == "quit" {
			break
		} else if choice == "T" || choice == "t" || choice == "torrent" || choice == "tor" || choice == "torr" {
			end = "torrents"
		} else {
			end = "layers"
		}
		resp, err := http.Get("http://" + ip + "/" + end)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Please enter new IP:")
			ip = getStdin(scanner)
			continue
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(end + ": " + string(body))
			resp.Body.Close()
		}
		fmt.Println("Query same IP for more info? [y/n]")
		t := getStdin(scanner)
		if t == "done" || t == "quit" {
			break
		} else if t == "y" || t == "yes" || t == "T" || t == "Yes" {
			//nothing, keep same ip
		} else {
			fmt.Println("Please enter new IP:")
			ip = getStdin(scanner)
		}
	}

	//for scanner.Scan() {
	//	if scanner.Text() != "" {
	//		layer := scanner.Text()
	//
	//		resp, err := http.Get("http://" + ip + "/layers/" + layer)
	//		if err != nil {
	//			fmt.Println(err)
	//			break
	//		} else {
	//			body, err := ioutil.ReadAll(resp.Body)
	//			if err == nil {
	//				fmt.Println("Layer: " + string(body))
	//				meta := torrentfile.AppendLayerDataToFile(layer, body)
	//				network.AddLayer(layer, meta)
	//			} else {
	//				fmt.Println(err)
	//			}
	//			resp.Body.Close()
	//			break
	//		}
	//	}
	//	break
	//}
	log.Println("Done receiving input, will only respond to http endpoints now")
}

func listenForFeedback() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter T or L for torrent or layer feedback")

	which := getStdin(scanner)

	var id string
	var hash string
	var btt common.SignableTransaction
	origin := myAddress2.ConvertToOriginInfo()
	if which == "T" || which == "t" {
		fmt.Println("What's the torrent TXID?")
		id = getStdin(scanner)
		trans := common.TorrentRepTrans{id, hash,
			common.RepMessage{true, true, true}}
		btt = common.SignableTransaction{origin, trans, common.TORRENT_REP, nil, nil, ""}
	} else {
		fmt.Println("What's the layer TXID?")
		id = getStdin(scanner)
		trans := common.LayerRepTrans{id, hash,
			true, true}
		btt = common.SignableTransaction{origin, trans, common.LAYER_REP, nil, nil, ""}
	}
	fmt.Println("And the hash?")
	hash = getStdin(scanner)

	signed := btt.SignAndSetTxID(&myAddress2.PrivateKey)
	log.Println("Gonna broadcast " + signed.TxID + " to blockchains")
	broadcastTransaction(env.BlockchainHost+":"+env.BlockchainPort+"/addTransaction", signed)
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()

	muxRouter.HandleFunc("/", handleIndexHelp).Methods("GET")

	muxRouter.HandleFunc("/torrents", handleGetTorrents).Methods("GET")
	muxRouter.HandleFunc("/layers", handleGetLayers).Methods("GET")

	muxRouter.HandleFunc("/layers/{layer}", handleGetLayer).Methods("POST")
	muxRouter.HandleFunc("/addLayer/{layer}", handleReceiveLayer).Methods("POST")

	return muxRouter
}

func registerTorrent(file common.TorrentFile) {
	if torrents == nil {
		torrents = make([]common.TorrentFile, 0)
	}
	torrents = append(torrents, file)
	for key, hash := range file.GetLayerHashMap() {
		addLayer(key, hash)
	}
}

func addLayer(id string, metadata common.LayerFileMetadata) {
	if layers == nil {
		layers = make(map[string]common.LayerFileMetadata, 0)
	}
	layers[id] = metadata
}

func handleIndexHelp(w http.ResponseWriter, r *http.Request) {

	io.WriteString(w, "Please use the following endpoints:\n\nGET /torrents to see available torrents\n"+
		"GET /layers to see available layers\nPOST /layers/<layerid> to POST a authentication transaction requesting the layer\n"+
		"POST /addLayer/<layerid> to POST raw layer data to add to internal server records (under <layerid>)")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
}

func handleGetLayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	layerId := vars["layer"]
	if layers == nil {
		http.Error(w, "Don't have any layers", http.StatusInternalServerError)
		return
	}

	var signedRequest common.SignableTransaction //should have received a transaction we can validate

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&signedRequest); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()

	if !signedRequest.Verify() {
		respondWithJSON(w, r, http.StatusUnauthorized, "Transaction is not signed correctly")
		fmt.Println("Transaction does not verify!")
		return
	}

	//TODO check reputation and determine access here
	//TODO add rules of some sort to filter users (maybe I have more stringent rules than other nodes; easily define rules)

	for key, layer := range layers {
		if key == layerId {
			file, err := os.Open(layer.GetUrl())

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data := make([]byte, layer.Size)

			file.ReadAt(data, layer.Begin)

			h := sha256.New()
			h.Write(data)
			io.WriteString(w, string(data))

			file.Close()
		}
	}

	fmt.Println("GET layer: " + layerId)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
}

func handleGetTorrents(w http.ResponseWriter, _ *http.Request) {
	if torrents == nil {
		fmt.Println("Don't have torrents; making new array")
		torrents = make([]common.TorrentFile, 0)
	}

	data, err := json.MarshalIndent(torrents, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("GET torrents")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	io.WriteString(w, string(data))
}

func handleGetLayers(w http.ResponseWriter, _ *http.Request) {
	if layers == nil {
		fmt.Println("Don't have torrents; making new array")
		layers = make(map[string]common.LayerFileMetadata, 0)
	}

	data, err := json.MarshalIndent(layers, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("GET torrents")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	io.WriteString(w, string(data))
}

func handleReceiveLayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	layerId := vars["layer"]
	var message string

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&message); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()

	byteMsg := []byte(message)

	respondWithJSON(w, r, http.StatusCreated, message)
	h := sha256.New()
	h.Write(byteMsg)
	if hex.EncodeToString(h.Sum(nil)) == layerId {

		fmt.Println("Received valid layer entry for " + layerId)
		layerdata := common.AppendLayerDataToFile(layerId, byteMsg)
		addLayer(layerId, layerdata)
	} else {
		fmt.Println("Hash didn't match")
	}
}

func respondWithJSON(w http.ResponseWriter, _ *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Write(response)
}

func broadcastTransaction(url string, trans common.SignableTransaction) {
	data, err := json.Marshal(trans)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Access-Control-Allow-Origin", "*")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

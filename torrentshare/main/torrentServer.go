package main

import (
	"bufio"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"strconv"
	"time"
	"github.com/denverquane/GoBlockShare/common"
)

type torrFileSpecs struct {
	url           string
	layerByteSize int64
	name          string
}

var myAddress common.PersonalAddress
var blockchainPort string

var torrents []common.TorrentFile
var layers map[string]common.LayerFileMetadata

func main() {
	err := godotenv.Load("common/.env")
	blockchainPort = os.Getenv("BLOCKCHAIN_PORT")
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(run())
}

func run() error {
	var httpAddr string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter the port this app should use [Ex: 7070]:")
	httpAddr = getStdin(scanner)

	jobs := make(chan torrFileSpecs, 10)
	results := make(chan common.TorrentFile, 10)

	for w := 1; w < 2; w++ {
		go torrentWorker(w, jobs, results)
	}
	totalJobs := 0

	fmt.Println("Would you like to generate a default torrent from README.md? [y/n]")
	text := getStdin(scanner)
	if text == "y" || text == "Y" || text == "yes" || text == "Yes" || text == "YES" {
		jobs <- torrFileSpecs{"README.md", 1000, "readme.md"}
		totalJobs++
	}

	for {
		fmt.Println("Any other torrents to provide? [y/n]")
		text = getStdin(scanner)
		if text == "n" || text == "no" || text == "No" || text == "done" || text == "quit" {
			break
		}
		fmt.Println("What is the location of your file? [Ex: test.txt, C:/Users/<...>/file.txt]:")
		url := getStdin(scanner)
		if url == "done" || url == "quit" {
			break
		}
		fmt.Println("What to call this torrent? [Ex: test.txt]")
		name := getStdin(scanner)
		if name == "done" || name == "quit" {
			break
		}
		jobs <- torrFileSpecs{url, 1000, name}
		totalJobs++
		fmt.Println("Added " + name + " to job queue")
	}

	myAddress = common.GenerateNewPersonalAddress()

	if scanner.Err() != nil {
		// handle error.
	}
	muxx := makeMuxRouter()

	for i := 0; i < totalJobs; i++ {
		file := <-results
		if file.Name != "" {
			trans := common.PublishTorrentTrans{file}
			origin := common.AddressToOriginInfo(myAddress)
			btt := common.SignableTransaction{origin, trans, "PUBLISH_TORRENT", nil, nil, ""}
			signed := btt.SignAndSetTxID(&myAddress.PrivateKey)
			log.Println("Gonna broadcast " + signed.TxID + " to blockchains")
			broadcastTransaction("http://localhost:" + blockchainPort + "/addTransaction", signed)
			registerTorrent(file)

			trans2 := common.TorrentRepTrans{signed.TxID,
				common.RepMessage{true, true, true}}
			btt2 := common.SignableTransaction{origin, trans2, "TORRENT_REP", nil, nil, ""}
			signed2 := btt2.SignAndSetTxID(&myAddress.PrivateKey)
			log.Println("Gonna broadcast " + signed2.TxID + " to blockchains")
			broadcastTransaction("http://localhost:" + blockchainPort + "/addTransaction", signed2)
		}
	}

	log.Println("Listening on ", httpAddr)
	if httpAddr == "8080" {
		log.Println("(This is the same port used internally for running Docker builds - are you running within a container?)")
	}

	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        muxx,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go listenForInput()

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

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()

	muxRouter.HandleFunc("/torrents", handleGetTorrents).Methods("GET")
	muxRouter.HandleFunc("/layers", handleGetLayers).Methods("GET")

	muxRouter.HandleFunc("/layers/{layer}", handleGetLayer).Methods("POST")
	//muxRouter.HandleFunc("/addTransaction", handleWriteTransaction).Methods("POST")
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

	for key, layer := range layers {
		if key == layerId {
			file, err := os.Open(layer.GetUrl())
			defer file.Close()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			data := make([]byte, layer.Size)

			file.ReadAt(data, layer.Begin)

			h := sha256.New()
			h.Write(data)
			io.WriteString(w, string(data))
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
	data, err := json.MarshalIndent(trans, "", "  ")
	var bytee = []byte(string(data))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bytee))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Access-Control-Allow-Origin", "*")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

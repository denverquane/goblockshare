package network

import (
	"encoding/json"
	"fmt"
	"github.com/denverquane/GoBlockShare/blockchain"
	//"github.com/denverquane/GoBlockShare/blockchain/transaction"
	"github.com/denverquane/GoBlockShare/files"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"crypto/sha256"
	"encoding/hex"
)

var globalBlockchain *blockchain.BlockChain
var torrents []files.TorrentFile
var layers map[string]files.LayerFileMetadata

func MakeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()

	muxRouter.HandleFunc("/torrents", handleGetTorrents).Methods("GET")
	muxRouter.HandleFunc("/layers", handleGetLayers).Methods("GET")
	muxRouter.HandleFunc("/blockchain", handleGetBlockchain).Methods("GET")

	muxRouter.HandleFunc("/layers/{layer}", handleGetLayer).Methods("GET")
	//muxRouter.HandleFunc("/addTransaction", handleWriteTransaction).Methods("POST")
	muxRouter.HandleFunc("/addLayer/{layer}", handleReceiveLayer).Methods("POST")

	return muxRouter
}

func RegisterBlockchain(chain *blockchain.BlockChain) {
	if globalBlockchain == nil {
		globalBlockchain = chain
	} else {
		fmt.Println("chain already set!")
	}
}

func RegisterTorrent(file files.TorrentFile) {
	if torrents == nil {
		torrents = make([]files.TorrentFile, 0)
	}
	torrents = append(torrents, file)
	for key, hash := range file.GetLayerHashMap() {
		AddLayer(key, hash)
	}
}

func AddLayer(id string, metadata files.LayerFileMetadata) {
	if layers == nil {
		layers = make(map[string]files.LayerFileMetadata, 0)
	}
	layers[id] = metadata
}

func handleGetBlockchain(w http.ResponseWriter, _ *http.Request) {
	// vars := mux.Vars(r)
	if globalBlockchain == nil {
		fmt.Println("Don't have blockchain; making new one")
		temp := blockchain.MakeInitialChain()
		globalBlockchain = &temp
	}

	data, err := json.MarshalIndent(*globalBlockchain, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("GET chain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	io.WriteString(w, string(data))
}

func handleGetLayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	layerId := vars["layer"]
	if torrents == nil {
		http.Error(w, "Don't have any torrents", http.StatusInternalServerError)
		return
	}

	for _, torr := range torrents {
		for key, meta := range torr.GetLayerHashMap() {
			if key == layerId {
				file, err := os.Open(torr.GetUrl())
				defer file.Close()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				data := make([]byte, meta.Size)

				file.ReadAt(data, meta.Begin)

				h := sha256.New()
				h.Write(data)
				io.WriteString(w, string(data) + "\n\nHASH: " + hex.EncodeToString(h.Sum(nil)))
			}
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
		torrents = make([]files.TorrentFile, 0)
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
		layers = make(map[string]files.LayerFileMetadata, 0)
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
		layerdata := files.AppendLayerDataToFile(layerId, byteMsg)
		AddLayer(layerId, layerdata)
	} else {
		fmt.Println("Hash didn't match")
	}
}

/* Below is an example of the input format for writing a transaction via the REST API:

{
"origin":
{
"address":"R9UtQ3QE4NrCxGuriwbI0qWCq0u7WqvjU0Q6muEd9Vk=",
"pubkeyx":86420643971005095497364485743353327828044563134904564182951237567725951244265,
"pubkeyy":84350736413375414420184852907452573247898047974475373171004335402121461174787
},
"txref":[],
"currency": "REP",
"quantity":5.99,
"payload":"Sending!",
"r":67869825206353784434575061723707880946031772528032340694185580017437536660581,
"s":6863529193914569235297749315606845644057909902475373433228108461283191248618,
"destAddr":"R9UtQ3QE4NrCxGuriwbI0qWCq0u7WqvjU0Q6muEd9Vk="
}

*/

//func handleWriteTransaction(w http.ResponseWriter, r *http.Request) {
//	// vars := mux.Vars(r)
//
//	var m transaction.RESTWrappedFullTransaction
//
//	decoder := json.NewDecoder(r.Body)
//	if err := decoder.Decode(&m); err != nil {
//		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
//		return
//	}
//	defer r.Body.Close()
//
//	trans, _ := m.ConvertToFull()
//	fmt.Println(trans.SignedTrans.ToString())
//	if !transaction.Verify(trans.SignedTrans) {
//		respondWithJSON(w, r, http.StatusBadRequest, "Transaction provided is invalid")
//		return
//	}
//
//	message, success := globalBlockchain.AddTransaction(trans, trans.SignedTrans.GetOrigin().Address)
//	if !success {
//		respondWithJSON(w, r, http.StatusBadRequest, message)
//	} else {
//		respondWithJSON(w, r, http.StatusCreated, message)
//	}
//	// BroadcastToAllPeers([]string{"http://localhost:8050/" + vars["channel"] + "/chain"}, newChain)
//}

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

package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"io"
	"time"
	"log"
	"github.com/joho/godotenv"
	"os"
	"github.com/denverquane/GoBlockShare/blockchain"
	"github.com/denverquane/GoBlockShare/common"
	"encoding/json"
	"strconv"
)

var globalBlockchain *blockchain.BlockChain

func main() {
	err := godotenv.Load("global.env")
	if err != nil {
		err = godotenv.Load("local.env")
		if err != nil {
			err = godotenv.Load("blockchain/local.env")
			if err != nil {
				log.Fatal(err)
			} else {
				log.Println("Using local env file")
			}
		} else {
			log.Println("Using local env file")
		}
	} else {
		log.Println("Using global env file")
	}

	temp := blockchain.MakeInitialChain()
	globalBlockchain = &temp

	log.Fatal(run())
}

func run() error {
	muxx := makeMuxRouter()
	port := os.Getenv("BLOCKCHAIN_PORT")
	log.Println("Starting blockchain server on port " + port)
	s := &http.Server{
		Addr:           ":" + port,
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



func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()

	muxRouter.HandleFunc("/", handleIndexHelp).Methods("GET")
	muxRouter.HandleFunc("/blockchain", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/block/{index}", handleGetBlock).Methods("GET")
	muxRouter.HandleFunc("/addTransaction", handleWriteTransaction).Methods("POST")

	return muxRouter
}

func handleIndexHelp(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Please use the following endpoints:\n\nGET /blockchain to see the entire recorded blockchain\n" +
		"GET /block/<index> to see a specific block of the chain\nPOST /addTransaction to POST a signed transaction " +
		"onto the blockchain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
}

func handleGetBlockchain(w http.ResponseWriter, _ *http.Request) {

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

func handleGetBlock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	index := vars["index"]
	i, err := strconv.Atoi(index)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if i < 0 || i >= globalBlockchain.Len() {
		http.Error(w, "Invalid index requested", http.StatusBadRequest)
		return
	}

	data, err := json.MarshalIndent(globalBlockchain.Blocks[i], "", " ")
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


func handleWriteTransaction(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)

	var jsonMessage common.JSONSignableTransaction
	var decodedMessage common.SignableTransaction

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&jsonMessage); err != nil {
		fmt.Println("couldn't decode ")
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}

	decodedMessage.Origin = jsonMessage.Origin
	decodedMessage.TransactionType = jsonMessage.TransactionType
	decodedMessage.TxID = jsonMessage.TxID
	decodedMessage.R = jsonMessage.R
	decodedMessage.S = jsonMessage.S

	defer r.Body.Close()

	switch jsonMessage.TransactionType {
	case "PUBLISH_TORRENT":
		var mm common.PublishTorrentTrans
		if err := json.Unmarshal([]byte(jsonMessage.Transaction), &mm); err != nil {
			log.Fatal(err)
		}
		decodedMessage.Transaction = mm
		break
	case "TORRENT_REP":
		var mm common.TorrentRepTrans
		if err := json.Unmarshal([]byte(jsonMessage.Transaction), &mm); err != nil {
			log.Fatal(err)
		}
		decodedMessage.Transaction = mm
		break
	}

	if !decodedMessage.Verify() {
		respondWithJSON(w, r, http.StatusBadRequest, "Transaction provided is invalid")
		return
	}

	message, success := globalBlockchain.AddTransaction(decodedMessage, jsonMessage.Origin.Address)
	if !success {
		respondWithJSON(w, r, http.StatusBadRequest, message)
	} else {
		respondWithJSON(w, r, http.StatusCreated, message)
	}
	// BroadcastToAllPeers([]string{"http://localhost:8050/" + vars["channel"] + "/chain"}, newChain)
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
package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"io"
	"time"
	"log"
	"github.com/denverquane/goblockshare/blockchain"
	"github.com/denverquane/goblockshare/common"
	"encoding/json"
	"strconv"
)

var globalBlockchain *blockchain.BlockChain

var env common.EnvVars

func main() {
	env = common.LoadEnvFromFile("blockchain")

	temp := blockchain.MakeInitialChain()
	globalBlockchain = &temp

	globalBlockchain.AddMockTransactions()

	log.Fatal(run())
}

func run() error {
	muxx := makeMuxRouter()
	log.Println("Starting blockchain server on port " + env.BlockchainPort)
	s := &http.Server{
		Addr:           ":" + env.BlockchainPort,
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
	muxRouter.HandleFunc("/reputation/{address}", handleGetReputation).Methods("GET")
	muxRouter.HandleFunc("/alias/{address}", handleGetAlias).Methods("GET")

	return muxRouter
}

func handleIndexHelp(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "Please use the following endpoints:\n\nGET /blockchain to see the entire recorded blockchain\n" +
		"GET /block/<index> to see a specific block of the chain\nPOST /addTransaction to POST a signed transaction " +
		"onto the blockchain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
}

func handleGetAlias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	alias := globalBlockchain.GetAddressAlias(common.Base64Address(address))

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	fmt.Println("Returning: " + alias)
	io.WriteString(w, alias)
}

func handleGetReputation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	rep := globalBlockchain.GetAddressRep(common.Base64Address(address))

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	resp, err := json.Marshal(rep)
	if err != nil {
		fmt.Println(err)
	}
	io.WriteString(w, string(resp))
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
	var jsonMessage common.JSONSignableTransaction

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&jsonMessage); err != nil {
		fmt.Println("couldn't decode ")
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}

	decoded := jsonMessage.ConvertToSignable()

	if decoded.TransactionType == "TORRENT_REP"{
		rep := decoded.Transaction.(common.TorrentRepTrans)
		if globalBlockchain.ProcessingReferencedTX(rep.TxID) {
			respondWithJSON(w, r, http.StatusBadRequest, "TX referenced in REP transaction is still being processed")
			return
		}
	}

	//TODO reject layer trans if being processed?

	if !decoded.Verify() {
		respondWithJSON(w, r, http.StatusBadRequest, "Transaction provided is invalid")
		return
	}

	success, err := globalBlockchain.AddTransaction(decoded, jsonMessage.Origin.Address)
	if !success {
		respondWithJSON(w, r, http.StatusBadRequest, err.Error())
	} else {
		respondWithJSON(w, r, http.StatusCreated, "Added transaction!")
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
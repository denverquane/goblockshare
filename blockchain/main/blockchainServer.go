package main

import (
	"encoding/json"
	"fmt"
	"github.com/denverquane/goblockshare/blockchain"
	"github.com/denverquane/goblockshare/common"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
	"github.com/gorilla/websocket"
)

var globalBlockchain *blockchain.BlockChain

var env common.EnvVars

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan blockchain.Block) // broadcast channel

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

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
	go handleMessages()

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

	muxRouter.HandleFunc("/ws", handleConnections)

	return muxRouter
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	clients[ws] = true

	for {
		var msg blockchain.Block
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
	}
}

func handleIndexHelp(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "Please use the following endpoints:\n\nGET /blockchain to see the entire recorded blockchain\n"+
		"GET /block/<index> to see a specific block of the chain\nPOST /addTransaction to POST a signed transaction "+
		"onto the blockchain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
}

type Alias struct {
	Data string
}

func handleGetAlias(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	alias := globalBlockchain.GetAddressAlias(common.Base64Address(address))
	aliass := Alias{alias}
	aliases := make([]Alias, 1)
	aliases[0] = aliass

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	resp, err := json.Marshal(aliases)
	if err != nil {
		fmt.Println(err)
	}
	io.WriteString(w, string(resp))
}

func handleGetReputation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]
	rep := globalBlockchain.GetAddressRep(common.Base64Address(address))
	reps := make([]common.JSONRepSummary, 1)
	reps[0] = rep

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	resp, err := json.Marshal(reps)
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

	if decoded.TransactionType == "TORRENT_REP" {
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

	success, err := globalBlockchain.AddTransaction(decoded, jsonMessage.Origin.Address, &broadcast)
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

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		fmt.Println("Have msg to transmit")
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

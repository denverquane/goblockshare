package network

import (
	"encoding/json"
	"fmt"
	"github.com/denverquane/GoBlockShare/Go/blockchain"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func MakeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/{channel}/chain", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/{channel}/users", handleGetUsers).Methods("GET")
	muxRouter.HandleFunc("/{channel}/postTransaction", handleWriteTransaction).Methods("POST")
	muxRouter.HandleFunc("/{channel}/chain", handleChainUpdate).Methods("POST")
	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	data, err := json.MarshalIndent(blockchain.GetChainByValue(vars["channel"]), "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("GET")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	io.WriteString(w, string(data))
}

func handleChainUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var m blockchain.BlockChain

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	chain, err := blockchain.CheckReplacementChain(vars["channel"], m)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.MarshalIndent(chain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(data))
}

func handleWriteTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var m blockchain.AuthTransaction

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	newChain, err := blockchain.WriteTransaction(vars["channel"], m)

	if err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, r, http.StatusCreated, newChain.GetNewestBlock())
}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var chain = blockchain.GetChainByValue(vars["channel"])
	authors := chain.GetNewestBlock().Users
	data, err := json.MarshalIndent(authors, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("GET")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	io.WriteString(w, string(data))
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

//func getPeersFromInput(){
//	var done = false
//	for !done {
//		reader := bufio.NewReader(os.Stdin)
//		fmt.Println("Enter URL to broadcast to, WITH port (ex: 127.0.0.1:8090), or \"quit\" if you're done: ")
//		text, _ := reader.ReadString('\n')
//		text = strings.Replace(text, "\n", "", 1)
//		text = strings.Replace(text, " ", "", 1)
//		fmt.Println("Entered: \"" + text + "\"")
//		if text == "quit" {
//			done = true
//		} else if text == "" {
//			fmt.Println("Empty string supplied")
//		} else {
//			BroadcastChain("http://" + text + "/chainUpdate", *globalChain)
//			Peers = append(Peers, "http://" + text + "/chainUpdate")
//		}
//	}
//}

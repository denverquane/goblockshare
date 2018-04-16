package network

import (
	"encoding/json"
	"fmt"
	"github.com/denverquane/GoBlockShare/blockchain"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func MakeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/createChannel", handleCreateChannel).Methods("POST")
	muxRouter.HandleFunc("/channels", handleGetChannels).Methods("GET")

	muxRouter.HandleFunc("/{channel}", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/{channel}", handleWriteTransaction).Methods("POST")

	muxRouter.HandleFunc("/{channel}/users", handleGetUsers).Methods("GET")
	muxRouter.HandleFunc("/{channel}/chain", handleChainUpdate).Methods("POST")
	return muxRouter
}

func handleCreateChannel(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)

	var m blockchain.AuthTransaction

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	chain, err := blockchain.CreateNewGlobalChannel(m)

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

func handleGetChannels(w http.ResponseWriter, r *http.Request) {
	channels := blockchain.GetChannelNames()
	data, err := json.MarshalIndent(channels, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("GET Channels")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "PUT")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	io.WriteString(w, string(data))
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	chain, err := blockchain.GetChainByValue(vars["channel"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.MarshalIndent(chain, "", "  ")
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

func handleChainUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var m blockchain.BlockChain

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	chain, err := blockchain.AttemptReplaceChain(vars["channel"], m)

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
	// BroadcastToAllPeers([]string{"http://localhost:8050/" + vars["channel"] + "/chain"}, newChain)
}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var chain, err = blockchain.GetChainByValue(vars["channel"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	authors := chain.GetNewestBlock().Users
	data, err := json.MarshalIndent(authors, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("GET Users")
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

package network

import (
	"encoding/json"
	"fmt"
	"github.com/denverquane/GoBlockShare/Go/blockchain"
	"github.com/gorilla/mux"
	"errors"
	"io"
	"net/http"
)

var ADMIN_CHANNEL_NAME string

func MakeMuxRouter(adminChannelName string) http.Handler {
	ADMIN_CHANNEL_NAME = adminChannelName
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/{channel}/chain", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/{channel}/users", handleGetUsers).Methods("GET")
	muxRouter.HandleFunc("/{channel}/postTransaction", handleWriteTransaction).Methods("POST")
	muxRouter.HandleFunc("/{channel}/chain", handleChainUpdate).Methods("POST")
	muxRouter.HandleFunc("/{channel}/create", handleCreateChannel).Methods("POST")
	return muxRouter
}

func handleCreateChannel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["channel"] != ADMIN_CHANNEL_NAME {
		respondWithJSON(w, r, http.StatusBadRequest, errors.New("Please use the endpoint for the admin channel"+
			"when attempting to create a new channel. For example, .../ADMIN/create"))
		return
	}

	var m blockchain.AuthTransaction

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	chain, err := blockchain.CreateNewChannel(m, ADMIN_CHANNEL_NAME)

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

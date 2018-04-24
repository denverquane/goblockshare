package network

import (
	"encoding/json"
	"fmt"
	"github.com/denverquane/GoBlockShare/blockchain"
	"github.com/denverquane/GoBlockShare/blockchain/transaction"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

var globalBlockchain *blockchain.BlockChain

func MakeMuxRouter(chain *blockchain.BlockChain) http.Handler {
	muxRouter := mux.NewRouter()
	globalBlockchain = chain

	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/addTransaction", handleWriteTransaction).Methods("POST")

	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)

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

/* Below is an example of the input format for writing a transaction via the REST API:

{
"originPubKeyX":"41465825910018896748506442457299597466934834109962972962658476739222369973795",
"originPubKeyY":"59682448553058160470866529273575025018646442288865005510432207524813798486893",
"originAddress":"05+tccYNwv6kwMDHFHLhpT2+syGQYhcvZrIUGMkj9vE=",
"signedMsg":"dsfgsd",
"txref":["tx1", "tx2"],
"r":"83272896655727237885461857009977546962509371591045400188157617593583499140053",
"s":"77837220821200439760189315101894538440367033391263344979880555787602867385798",
"destAddr":"UVoty8GhdxPK4ZfxNUGIGSmDcumFGk4+3Sc8R1e7D08="
}

*/

func handleWriteTransaction(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)

	var m transaction.RESTWrappedFullTransaction

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	fmt.Println("received transaction")
	trans := m.ConvertToFull()
	if !trans.Verify() {
		respondWithJSON(w, r, http.StatusBadRequest, "Transaction provided is invalid")
		return
	}
	globalBlockchain.AddTransaction(trans)
	// Return some JSON response, even if the block isn't mined yet (but first validate the transaction's validity)

	respondWithJSON(w, r, http.StatusCreated, globalBlockchain.GetNewestBlock())
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

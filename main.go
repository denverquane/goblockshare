package main

import (
	"chatProgram/blockchain"
	"github.com/gorilla/mux"
	//"fmt"
	"fmt"
	"os"
	"log"
	"net/http"
	"time"
	"encoding/json"
	"io"
	"github.com/joho/godotenv"
	"github.com/davecgh/go-spew/spew"
)

var Chain blockchain.BlockChain

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		Chain = blockchain.BlockChain{1, make([]blockchain.Block, 1)}
		block := blockchain.InitialBlock()
		fmt.Println(block.ToString())
		Chain.Blocks[0] = block
	}()
	log.Fatal(run())

}

func run() error {
	mux := makeMuxRouter()
	httpAddr := os.Getenv("ADDR")
	log.Println("Listening on ", os.Getenv("ADDR"))
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
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
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/postTransaction", handleWriteBlock).Methods("POST")
	muxRouter.HandleFunc("/chainUpdate", handleChainUpdate).Methods("POST")
	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.MarshalIndent(Chain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleChainUpdate(w http.ResponseWriter, r *http.Request) {
	var m blockchain.BlockChain

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	if m.IsValid() {
		if m.Length > Chain.Length {
			if blockchain.AreChainsSameBranch(m, Chain) {
				Chain = m
				fmt.Println("Valid blockchain supplied! Replaced with: ")
				spew.Dump(Chain)
			} else {
				fmt.Println("Chains are of different branches! Not replacing mine")
			}
		} else {
			fmt.Println("Chains are the same length; not replacing anything")
		}
	} else {
		fmt.Println("Invalid blockchain supplied; not replacing anything")
	}

	bytes, err := json.MarshalIndent(Chain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m blockchain.Transaction

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	oldBlock := Chain.GetNewestBlock()
	newBlock, err := blockchain.GenerateBlock(oldBlock, m)
	fmt.Println("New block:\n" + newBlock.ToString())
	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, m)
		return
	}
	if blockchain.IsBlockSequenceValid(newBlock, oldBlock) {
		Chain.Blocks = append(Chain.Blocks, newBlock)
		Chain.Length++
		//Block = blockchain.CheckLongerChain(newBlock, Block)
		fmt.Println("Successfully added: {" + m.ToString() + "} to the chain")
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)

}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

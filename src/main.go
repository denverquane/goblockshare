package main

import (
	"chatProgram/src/blockchain"
	"github.com/gorilla/mux"
	//"fmt"
	"fmt"
	"os"
	"log"
	"net/http"
	"time"
	"encoding/json"
	"io"
	"github.com/davecgh/go-spew/spew"
	"bytes"
	"io/ioutil"
	"github.com/joho/godotenv"
	"bufio"
	"strings"
)

var Chain blockchain.BlockChain
var Peers []string

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		Peers = make([]string, 0)
		Chain = blockchain.BlockChain{make([]blockchain.Block, 1)}
		block := blockchain.InitialBlock()
		fmt.Println(block.ToString())
		Chain.Blocks[0] = block
	}()
	log.Fatal(run())

}

func getPeersFromInput(){
	var done = false
	for !done {
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Enter URL to broadcast to, WITH port (ex: 127.0.0.1:8090), or \"quit\" if you're done: ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", 1)
		text = strings.Replace(text, " ", "", 1)
		fmt.Println("Entered: \"" + text + "\"")
		if text == "quit" {
			done = true
		} else if text == "" {
			fmt.Println("Empty string supplied")
		} else {
			broadcastChain("http://" + text + "/chainUpdate", Chain)
			Peers = append(Peers, "http://" + text + "/chainUpdate")
		}
	}
}

func broadcastToAllPeers() {
	for _, v := range Peers {
		broadcastChain(v, Chain)
	}
}

func run() error {
	muxx := makeMuxRouter()
	httpAddr := os.Getenv("ADDR")
	log.Println("Listening on ", os.Getenv("ADDR"))
	go getPeersFromInput()

	s := &http.Server{
		Addr:           ":" + httpAddr,
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
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/postTransaction", handleWriteBlock).Methods("POST")
	muxRouter.HandleFunc("/chainUpdate", handleChainUpdate).Methods("POST")
	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, _ *http.Request) {
	data, err := json.MarshalIndent(Chain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("GET")
	io.WriteString(w, string(data))
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
		if m.Len() > Chain.Len() {
			if blockchain.AreChainsSameBranch(m, Chain) {
				Chain = m
				fmt.Println("Valid blockchain supplied! Replaced with: ")
				spew.Dump(Chain)
				broadcastToAllPeers()
			} else {
				fmt.Println("Chains are of different branches! Not replacing mine")
			}
		} else if Chain.Len() == 1 && m.Len() == 1 {
			Chain = m
			fmt.Println("Both chains are 1 length; replacing mine!")
			//spew.Dump(Chain)
		}else {
			fmt.Println("Chains are the same or lesser length; not replacing anything")
		}
	} else {
		fmt.Println("Invalid blockchain supplied; not replacing anything")
	}

	data, err := json.MarshalIndent(Chain, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(data))
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
	newBlock := blockchain.GenerateBlock(oldBlock, m)
	fmt.Println("New block:\n" + newBlock.ToString())

	if blockchain.IsBlockSequenceValid(newBlock, oldBlock) {
		Chain.Blocks = append(Chain.Blocks, newBlock)
		//Block = blockchain.CheckLongerChain(newBlock, Block)
		fmt.Println("Successfully added: {" + m.ToString() + "} to the chain")
		broadcastToAllPeers()
	}

	respondWithJSON(w, r, http.StatusCreated, newBlock)
}

func respondWithJSON(w http.ResponseWriter, _ *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func broadcastChain(url string, chain blockchain.BlockChain) {
	data, err := json.MarshalIndent(chain, "", "  ")
	//fmt.Println(string(data))
	var bytee = []byte(string(data))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bytee))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
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

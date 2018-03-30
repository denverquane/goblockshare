package main

import (
	"os"
	"log"
	"net/http"
	"time"
	"github.com/joho/godotenv"
	"github.com/denverquane/GoBlockChat/Go/network"
	"github.com/denverquane/GoBlockChat/Go/blockchain"
)

func main() {
	err := godotenv.Load("Go/.env")
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(run())
}

func run() error {
	muxx := network.MakeMuxRouter()

	httpAddr := os.Getenv("PORT")
	log.Println("Listening on ", os.Getenv("PORT"))

	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        muxx,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	makeGlobalChain()

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func makeGlobalChain() {
	users := make([]blockchain.UserPassPair, 2)
	users[0] = blockchain.UserPassPair{"admin", "pass"}
	users[1] = blockchain.UserPassPair{"user1", "pass"}
	chain := blockchain.MakeInitialChain(users)
	blockchain.SetGlobalChain(chain)
}


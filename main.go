package main

import (
	"github.com/denverquane/GoBlockShare/blockchain"
	"github.com/denverquane/GoBlockShare/network"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(run())
}

func run() error {
	httpAddr := os.Getenv("PORT")

	muxx := network.MakeMuxRouter()

	log.Println("Listening on ", os.Getenv("PORT"))
	if (httpAddr == "8080") {
		log.Println("(This is the same port used internally for running Docker builds - are you running within a container?)")
	}

	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        muxx,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	makeAGlobalChain("TEST_CHANNEL", blockchain.UserPassPair{"user", "password"})

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

//makeGlobalChain constructs an initial blockchain for the program, using a specified global admin username/password
func makeAGlobalChain(channelName string, user blockchain.UserPassPair) {
	chain := blockchain.MakeInitialChain([]blockchain.UserPassPair{user})
	blockchain.SetChannelChain(channelName, chain)
}

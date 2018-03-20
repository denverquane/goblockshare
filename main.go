package main

import (
	"os"
	"log"
	"net/http"
	"time"
	"chatProgram/network"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(run())
}

func run() error {
	muxx := network.MakeMuxRouter()

	httpAddr := os.Getenv("ADDR")
	log.Println("Listening on ", os.Getenv("ADDR"))

	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        muxx,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	network.FetchOrMakeChain(muxx)

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}


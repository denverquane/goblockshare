package main

import (
	"os"
	"log"
	"net/http"
	"time"
	"github.com/joho/godotenv"
	"GoBlockChat/Go/network"
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

	network.FetchOrMakeChain(muxx)

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}


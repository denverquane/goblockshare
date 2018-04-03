package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/denverquane/GoBlockChat/Go/blockchain"
	"github.com/denverquane/GoBlockChat/Go/network"
	"github.com/joho/godotenv"
	"hash"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	err := godotenv.Load("Go/.env")
	if err != nil {
		log.Fatal(err)
	}

	//file, err := ioutil.ReadFile("Go/main.go")
	//if err != nil {
	//	fmt.Print(err)
	//}
	//fmt.Println(string(file))
	log.Fatal(run())
}

func run() error {
	httpAddr := os.Getenv("PORT")
	version := os.Getenv("VERSION")
	h := hashDirectory("./Go")
	fmt.Printf("GoBlockShare Version: "+version+", Checksum: %x\n", h)

	muxx := network.MakeMuxRouter()

	log.Println("Listening on ", os.Getenv("PORT"))

	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        muxx,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	makeGlobalChain(version)

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func makeGlobalChain(version string) {
	users := make([]blockchain.UserPassPair, 1)
	users[0] = blockchain.UserPassPair{"admin", "pass"}
	chain := blockchain.MakeInitialChain(users, version)
	blockchain.SetChannelChain("Admin Channel", chain)
}

func hashDirectory(dir string) string {
	b, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Print(err)
	}

	hash := sha256.New()
	for _, v := range b {
		hash = recursivelyHashFiles(hash, v, dir+"/")
	}
	return (string)(hash.Sum(nil))
}

func recursivelyHashFiles(hasher hash.Hash, info os.FileInfo, path string) hash.Hash {
	if info.IsDir() && !strings.Contains(info.Name(), ".git") && !strings.Contains(info.Name(), ".idea") {
		// fmt.Println("Opening dir: " + path + info.Name())
		b, err := ioutil.ReadDir(path + info.Name())
		if err != nil {
			fmt.Print(err)
		}
		for _, v := range b {
			hasher = recursivelyHashFiles(hasher, v, path+info.Name()+"/")
		}
	} else if !strings.Contains(info.Name(), ".git") && !strings.Contains(info.Name(), ".idea") {
		// fmt.Println("Hashing File: " + (path + info.Name()))
		file, err := os.Open(path + info.Name())

		if err != nil {
			log.Fatal(err)
		}

		if _, err := io.Copy(hasher, file); err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("%x\n\n", hasher.Sum(nil))
		file.Close()
	}
	return hasher
}

package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/denverquane/GoBlockShare/Go/blockchain"
	"github.com/denverquane/GoBlockShare/Go/network"
	"github.com/joho/godotenv"
	"hash"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"encoding/hex"
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
	adminChannelName := os.Getenv("ADMIN_CHANNEL_NAME")
	h := hashDirectory("./Go")
	fmt.Printf("GoBlockShare Version: "+version+", Checksum: %s\n", h)

	muxx := network.MakeMuxRouter(adminChannelName, version, h)

	log.Println("Listening on ", os.Getenv("PORT"))

	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        muxx,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	makeGlobalChain(version, adminChannelName)

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func makeGlobalChain(version string, adminChannel string) {
	users := make([]blockchain.UserPassPair, 1)
	users[0] = blockchain.UserPassPair{"admin", "pass"}
	chain := blockchain.MakeInitialChain(users, version)
	blockchain.SetChannelChain(adminChannel, chain)
}

//hashDirectory receives a relative or absolute path, and hashes together all the files contained within the directory.
//This is for verification that separate program instances are all running the same version, with no modifications to
//the original source code
func hashDirectory(dir string) string {
	b, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Print(err)
	}

	h := sha256.New()
	for _, v := range b {
		h = recursivelyHashFiles(h, v, dir+"/")
	}
	return hex.EncodeToString(h.Sum(nil))
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

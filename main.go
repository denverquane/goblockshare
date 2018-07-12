package main

import (
	"bufio"
	"fmt"
	"github.com/denverquane/GoBlockShare/files"
	"github.com/denverquane/GoBlockShare/network"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

//var Signed transaction.SignableTransaction
//var Wallet1 wallet.Wallet
//var Wallet2 wallet.Wallet

func main() {
	//Wallet1 = wallet.MakeNewWallet()
	//Wallet2 = wallet.MakeNewWallet()
	//Signed = Wallet1.MakeTransaction(5.99, Wallet2.GetAddress().Address)
	//fmt.Println(Signed.ToString())
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(run())
}

func run() error {
	httpAddr := os.Getenv("PORT")

	globalHashes := make([]string, 0)

	//globalChain := blockchain.MakeInitialChain(Wallet1.GetAddress().Address)

	/************ Testing wallet block ***************/

	//message, _ := globalChain.AddTransaction(transaction.MakeFull(Signed, []string{}), Wallet1.GetAddress().Address) //empty TXREF for now
	//fmt.Println(message)
	//Wallet1.UpdateBalances(globalChain)
	//Wallet2.UpdateBalances(globalChain)
	/*************************************************/
	torr, err := files.MakeTorrentFileFromFile(1000, "README.md", "readme.md")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Total checksum: " + torr.TotalHash)

	//TODO don't store the actual data? Just store a reference to the file's location, and the offset bytes
	for _, v := range torr.SegmentHashKeys {
		fmt.Println("I know of layer " + v)
		globalHashes = append(globalHashes, v)
	}

	fmt.Println("Please enter the port to be used for communicating with this node. Type \"done\" when you are finished")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Text() != "" {
			httpAddr = scanner.Text()
			break
		}
	}

	if scanner.Err() != nil {
		// handle error.
	}

	//TODO Remove nil address reference
	muxx := network.MakeMuxRouter()
	network.Torrents = make([]files.TorrentFile, 0)
	network.Torrents = append(network.Torrents, torr)

	log.Println("Listening on ", httpAddr)
	if httpAddr == "8080" {
		log.Println("(This is the same port used internally for running Docker builds - are you running within a container?)")
	}

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

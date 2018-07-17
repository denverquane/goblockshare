package main

import (
	"fmt"
	"github.com/denverquane/GoBlockShare/files"
	"github.com/denverquane/GoBlockShare/network"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/denverquane/GoBlockShare/blockchain"
	"github.com/denverquane/GoBlockShare/blockchain/transaction"
	"bufio"
	"strconv"
)

//var Signed transaction.SignableTransaction
//var Wallet1 wallet.Wallet
//var Wallet2 wallet.Wallet

type torrFileSpecs struct {
	url string
	layerByteSize int64
	name string
}

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

	//globalHashes := make([]string, 0)

	globalChain := blockchain.MakeInitialChain()

	myAddress := transaction.GenerateNewPersonalAddress()

	jobs := make(chan torrFileSpecs, 10)

	results := make(chan files.TorrentFile, 10)

	for w:=1; w < 2; w++ {
		go torrentWorker(w, jobs, results)
	}

	/************ Testing wallet block ***************/

	//message, _ := globalChain.AddTransaction(transaction.MakeFull(Signed, []string{}), Wallet1.GetAddress().Address) //empty TXREF for now
	//fmt.Println(message)
	//Wallet1.UpdateBalances(globalChain)
	//Wallet2.UpdateBalances(globalChain)
	/*************************************************/


	scanner := bufio.NewScanner(os.Stdin)
	totalJobs := 0
	fmt.Println("Please enter the filename you wish to broadcast as a torrent on the blockchain.")
	fmt.Println("(The path should be relative to main.go). Type \"done\" when complete")
	for scanner.Scan() {
		if scanner.Text() != "" && scanner.Text() != "done" {
			url := scanner.Text()
			fmt.Println("What would you like this torrent to be called?")
			for scanner.Scan() {
				name := scanner.Text()
				if name == "" {
					name = url
				}
				jobs <- torrFileSpecs{url, 1000 * 1000, name}
				totalJobs++
				fmt.Println("Added " + name + " to job queue")
				break
			}
		} else {
			close(jobs)
			break
		}
		fmt.Println("Please enter the filename you wish to broadcast as a torrent on the blockchain.")
		fmt.Println("(The path should be relative to main.go). Type \"done\" when complete")
	}

	if scanner.Err() != nil {
		// handle error.
	}
	muxx := network.MakeMuxRouter()
	network.Torrents = make([]files.TorrentFile, 0)
	network.GlobalBlockchain = &globalChain

	origin := transaction.AddressToOriginInfo(myAddress)
	for i := 0; i<totalJobs; i++ {
		file := <-results
		trans := transaction.PublishTorrentTrans{file}
		btt := transaction.TorrentTransaction{origin, trans, nil, nil}
		signed := transaction.Sign(&myAddress.PrivateKey, btt)
		full := transaction.MakeFull(signed, nil)
		globalChain.AddTransaction(full, myAddress.Address)
		network.Torrents = append(network.Torrents, file)

		trans2 := transaction.TorrentRepTrans{full.TxID, transaction.RepMessage{true, true, true}}
		btt2 := transaction.TorrentTransaction{origin, trans2, nil, nil}
		signed2 := transaction.Sign(&myAddress.PrivateKey, btt2)
		full2 := transaction.MakeFull(signed2, nil)
		globalChain.AddTransaction(full2, myAddress.Address)
	}

	//torr, err := files.MakeTorrentFileFromFile(1000, "README.md", "readme.md")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println("Total checksum: " + torr.TotalHash)
	//
	////TODO don't store the actual data? Just store a reference to the file's location, and the offset bytes
	//for _, v := range torr.LayerHashKeys {
	//	fmt.Println("I know of layer " + v)
	//	globalHashes = append(globalHashes, v)
	//}

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

func torrentWorker(id int, jobs <-chan torrFileSpecs, results chan<- files.TorrentFile) {
	for job := range jobs {
		torr, err := files.MakeTorrentFileFromFile(job.layerByteSize, job.url, job.name)
		if err != nil {
			fmt.Println("Worker " + strconv.Itoa(id) + " has error: " + err.Error())
		} else {
			fmt.Println("Worker " + strconv.Itoa(id) + " finished " + job.name)
		}
		results <- torr
	}
}

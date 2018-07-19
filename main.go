package main

import (
	"bufio"
	"fmt"
	"github.com/denverquane/GoBlockShare/blockchain"
	"github.com/denverquane/GoBlockShare/blockchain/transaction"
	"github.com/denverquane/GoBlockShare/files"
	"github.com/denverquane/GoBlockShare/network"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type torrFileSpecs struct {
	url           string
	layerByteSize int64
	name          string
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(run())
}

var MyAddress transaction.PersonalAddress

func run() error {
	var httpAddr string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter the port this app should use [Ex: 7070]:")
	httpAddr = getStdin(scanner)

	jobs := make(chan torrFileSpecs, 10)
	results := make(chan files.TorrentFile, 10)

	for w := 1; w < 2; w++ {
		go torrentWorker(w, jobs, results)
	}
	totalJobs := 0

	fmt.Println("Would you like to generate a default torrent from README.md? [y/n]")
	text := getStdin(scanner)
	if text == "y" || text == "Y" || text == "yes" || text == "Yes" || text == "YES" {
		jobs <- torrFileSpecs{"README.md", 1000, "readme.md"}
		totalJobs++
	}

	for {
		fmt.Println("Any other torrents to provide? [y/n]")
		text = getStdin(scanner)
		if text == "n" || text == "no" || text == "No" || text == "done" || text == "quit" {
			break
		}
		fmt.Println("What is the location of your file? [Ex: test.txt, C:/Users/<...>/file.txt]:")
		url := getStdin(scanner)
		if url == "done" || url == "quit" {
			break
		}
		fmt.Println("What to call this torrent? [Ex: test.txt]")
		name := getStdin(scanner)
		if name == "done" || name == "quit" {
			break
		}
		jobs <- torrFileSpecs{url, 1000, name}
		totalJobs++
		fmt.Println("Added " + name + " to job queue")
	}

	globalChain := blockchain.MakeInitialChain()

	MyAddress = transaction.GenerateNewPersonalAddress()
	//trans := MyAddress.GenerateNullTransaction()
	//fmt.Println(trans.ToString())

	if scanner.Err() != nil {
		// handle error.
	}
	muxx := network.MakeMuxRouter()
	network.RegisterBlockchain(&globalChain)

	for i := 0; i < totalJobs; i++ {
		file := <-results
		if file.Name != "" {
			trans := transaction.PublishTorrentTrans{file, "PUBLISH_TORRENT"}
			full := globalChain.CreateAndAddTransaction(MyAddress, trans)
			network.RegisterTorrent(file)

			trans2 := transaction.TorrentRepTrans{full.TxID, transaction.RepMessage{true, true, true}, "TORRENT_REP"}
			globalChain.CreateAndAddTransaction(MyAddress, trans2)
		}
	}

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

	go listenForInput()

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func getStdin(scanner *bufio.Scanner) string {
	for scanner.Scan() {
		if scanner.Text() != "" {
			return scanner.Text()
		}
	}
	return ""
}

//torrentWorker represents a worker that should process a file on the local filesystem and process it into an
//internal TorrentFile. This could potentially still be a bottleneck from the underlying I/O filesystem, but processing
//entirely separate files is inherently parallelizable
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

func listenForInput() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter the IP and port of a node you wish to query for torrents and layers")
	fmt.Println("For example (don't include http://): localhost:7070")
	fmt.Println("Enter \"done\" or \"quit\" at any time to finish querying peers")

	ip := getStdin(scanner)

	for ip != "quit" && ip != "done" {
		fmt.Println("Query \"" + ip + "\" for torrents, or layers? [T/L]")
		choice := getStdin(scanner)
		var end string
		if choice == "done" || choice == "quit" {
			break
		} else if choice == "T" || choice == "t" || choice == "torrent" || choice == "tor" || choice == "torr" {
			end = "torrents"
		} else {
			end = "layers"
		}
		resp, err := http.Get("http://" + ip + "/" + end)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Please enter new IP:")
			ip = getStdin(scanner)
			continue
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println(end + ": " + string(body))
			resp.Body.Close()
		}
		fmt.Println("Query same IP for more info? [y/n]")
		t := getStdin(scanner)
		if t == "done" || t == "quit" {
			break
		} else if t == "y" || t == "yes" || t == "T" || t == "Yes" {
			//nothing, keep same ip
		} else {
			fmt.Println("Please enter new IP:")
			ip = getStdin(scanner)
		}
	}

	//for scanner.Scan() {
	//	if scanner.Text() != "" {
	//		layer := scanner.Text()
	//
	//		resp, err := http.Get("http://" + ip + "/layers/" + layer)
	//		if err != nil {
	//			fmt.Println(err)
	//			break
	//		} else {
	//			body, err := ioutil.ReadAll(resp.Body)
	//			if err == nil {
	//				fmt.Println("Layer: " + string(body))
	//				meta := files.AppendLayerDataToFile(layer, body)
	//				network.AddLayer(layer, meta)
	//			} else {
	//				fmt.Println(err)
	//			}
	//			resp.Body.Close()
	//			break
	//		}
	//	}
	//	break
	//}
	log.Println("Done receiving input, will only respond to http endpoints now")
}

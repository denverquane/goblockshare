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
	"io/ioutil"
)

type torrFileSpecs struct {
	url string
	layerByteSize int64
	name string
}

func main() {
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
				jobs <- torrFileSpecs{url, 1000, name}
				totalJobs++
				fmt.Println("Added " + name + " to job queue")
				break
			}
		} else {
			fmt.Println("Done receiving inputs")
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
	network.RegisterBlockchain(&globalChain)
	
	for i := 0; i<totalJobs; i++ {
		file := <-results
		if file.Name != "" {
			trans := transaction.PublishTorrentTrans{file, "PUBLISH_TORRENT"}
			full := globalChain.CreateAndAddTransaction(myAddress, trans)
			network.RegisterTorrent(file)

			trans2 := transaction.TorrentRepTrans{full.TxID, transaction.RepMessage{true, true, true}, "TORRENT_REP"}
			globalChain.CreateAndAddTransaction(myAddress, trans2)
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

	for scanner.Scan() {

		if scanner.Text() != "quit" && scanner.Text() != "done" {
			ip := scanner.Text()
			resp, err := http.Get("http://" + ip + "/torrents")
			if err != nil {
				fmt.Println(err)
				break
			} else {
				body, _ := ioutil.ReadAll(resp.Body)
				fmt.Println("Layers: " + string(body))
				resp.Body.Close()
				fmt.Println("Would you like to fetch one of the layers?")
				for scanner.Scan() {
					if scanner.Text() == "y" || scanner.Text() == "Y" || scanner.Text() == "yes" || scanner.Text() == "Yes" {
						fmt.Println("Which layer?")
						for scanner.Scan() {
							if scanner.Text() != "" {
								layer := scanner.Text()
								resp, err := http.Get("http://" + ip + "/layers/" + layer)
								if err != nil {
									fmt.Println(err)
									break
								} else {
									body, err := ioutil.ReadAll(resp.Body)
									if err == nil {
										fmt.Println("Layer: " + string(body))
										meta := files.AppendLayerDataToFile(layer, body)
										network.AddLayer(layer, meta)
									} else {
										fmt.Println(err)
									}
									resp.Body.Close()
									break
								}
							}
							break
						}
					}
					break
				}
			}
		} else {
			break
		}
		fmt.Println("Please enter the IP and port of a node you wish to query for torrents and layers")
		fmt.Println("For example (don't include http://): localhost:7070")
	}
}

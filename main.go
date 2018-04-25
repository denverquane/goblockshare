package main

import (
	"bufio"
	"fmt"
	"github.com/denverquane/GoBlockShare/blockchain"
	"github.com/denverquane/GoBlockShare/blockchain/transaction"
	"github.com/denverquane/GoBlockShare/network"
	"github.com/denverquane/GoBlockShare/wallet"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

var Signed transaction.SignedTransaction
var Wallet1 wallet.Wallet
var Wallet2 wallet.Wallet

func main() {
	Wallet1 = wallet.MakeNewWallet()
	Wallet2 = wallet.MakeNewWallet()
	Signed = Wallet1.MakeTransaction(5.99, Wallet2.GetAddress().Address)
	fmt.Println(Signed.ToString())
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(run())
}

func run() error {
	httpAddr := os.Getenv("PORT")

	globalChain := blockchain.MakeInitialChain(Wallet1.GetAddress().Address)

	/************ Testing wallet block ***************/

	globalChain.AddTransaction(Signed.MakeFull([]string{})) //empty TXREF for now
	Wallet1.InitializeBalances(globalChain)
	fmt.Println(Wallet1.GetBalances())
	Wallet2.InitializeBalances(globalChain)
	fmt.Println(Wallet2.GetBalances())

	/*************************************************/
	fmt.Println("Please enter the nodes you would like to communicate with. Type \"done\" when you are finished")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		if scanner.Text() == "done" || scanner.Text() == "quit" {
			break
		}
		fmt.Println(scanner.Text())
	}

	if scanner.Err() != nil {
		// handle error.
	}

	muxx := network.MakeMuxRouter(&globalChain)

	log.Println("Listening on ", os.Getenv("PORT"))
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

package wallet

import (
	"fmt"
	"github.com/denverquane/GoBlockShare/blockchain"
	"github.com/denverquane/GoBlockShare/blockchain/transaction"
	"strconv"
)

type Wallet struct {
	addresses   []transaction.PersonalAddress
	inputTxIDs  []string
	outputTxIDs []string
	balanceMap  map[string]float64
}

func MakeNewWallet() Wallet {
	address := transaction.GenerateNewPersonalAddress()
	fmt.Println(address.Address)
	return Wallet{[]transaction.PersonalAddress{address}, []string{},
		[]string{}, make(map[string]float64, 0)}
}

func (wallet Wallet) GetAddress() transaction.PersonalAddress {
	return wallet.addresses[0]
}

func (wallet Wallet) GetBalances() string {
	str := "Wallet Balances: \n"
	for i, v := range wallet.balanceMap {
		str += ("  " + i + " has balance of " + strconv.FormatFloat(v, 'f', -1, 64) + "\n")
	}
	return str
}

//TODO This is for testing!!! Don't rely on this!
func (wallet Wallet) getOriginInfo() transaction.OriginInfo {
	return transaction.AddressToOriginInfo(wallet.addresses[0])
}

func (wallet Wallet) MakeTransaction(quantity float64, dest transaction.Base64Address) transaction.SignedTransaction {
	unsigned := transaction.SignedTransaction{wallet.getOriginInfo(), dest, quantity,
		"Sending!", nil, nil}
	return unsigned.SignMessage(&wallet.addresses[0].PrivateKey)
}

func (wallet *Wallet) InitializeBalances(blockchain blockchain.BlockChain) {
	for _, block := range blockchain.Blocks { // look at all blocks
		wallet.UpdateBalances(block)
	}
}

func (wallet *Wallet) UpdateBalances(block blockchain.Block) {
	for _, tx := range block.Transactions { // all transactions in the block
		for _, addr := range wallet.addresses { // look through our personal addresses
			if addr.Address == tx.SignedTrans.DestAddr { // if the output address matches one of ours
				alreadyProcessed := false
				for _, input := range wallet.inputTxIDs { // make sure we haven't recorded the transaction already
					if input == tx.TxID {
						alreadyProcessed = true
					}
				}
				if !alreadyProcessed {
					wallet.inputTxIDs = append(wallet.inputTxIDs, tx.TxID)
					//TODO get specific currency
					currency := "REP"
					if _, ok := wallet.balanceMap[currency]; ok {
						wallet.balanceMap[currency] += tx.SignedTrans.Quantity
					} else {
						wallet.balanceMap[currency] = tx.SignedTrans.Quantity
					}

					fmt.Println("Recorded +" + strconv.FormatFloat(tx.SignedTrans.Quantity, 'f', -1, 64) +
						" " + currency + " to my wallet!")
				}

			} else if addr.Address == tx.SignedTrans.Origin.Address {
				alreadyProcessed := false
				for _, output := range wallet.outputTxIDs { // make sure we haven't recorded the transaction already
					if output == tx.TxID {
						alreadyProcessed = true
					}
				}
				if !alreadyProcessed {
					wallet.outputTxIDs = append(wallet.outputTxIDs, tx.TxID)
					//TODO get specific currency
					currency := "REP"
					if _, ok := wallet.balanceMap[currency]; ok {
						wallet.balanceMap[currency] -= tx.SignedTrans.Quantity
					} else {
						fmt.Println("Recorded a negative balance in my account...")
						wallet.balanceMap[currency] = -tx.SignedTrans.Quantity
					}

					fmt.Println("Recorded -" + strconv.FormatFloat(tx.SignedTrans.Quantity, 'f', -1, 64) +
						" " + currency + " from my wallet!")
				}
			}
		}
	}
}

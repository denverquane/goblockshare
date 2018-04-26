package wallet

import (
	"fmt"
	"github.com/denverquane/GoBlockShare/blockchain"
	"github.com/denverquane/GoBlockShare/blockchain/transaction"
	"strconv"
)

type Wallet struct {
	lastProcessedBlock int
	addresses          []transaction.PersonalAddress
	balance            float64
	ChannelRecords     map[string]ChannelRecord //map of all channels we (might) have access to
}

func MakeNewWallet() Wallet {
	address := transaction.GenerateNewPersonalAddress()
	fmt.Println(address.Address)
	return Wallet{-1, []transaction.PersonalAddress{address}, 0.0, make(map[string]ChannelRecord, 0)}
}

func (wallet Wallet) GetAddress() transaction.PersonalAddress {
	return wallet.addresses[0]
}

func (wallet Wallet) GetBalances() string {
	str := "Wallet Balance: " + strconv.FormatFloat(wallet.balance, 'f', -1, 64) + " REP"
	for name, record := range wallet.ChannelRecords {
		if record.haveToken {
			str += "\n                " + "1 " + name
		} else {
			str += "\n                " + "0 " + name
		}
	}

	return str
}

//TODO This is for testing!!! Don't rely on this!
func (wallet Wallet) getOriginInfo() transaction.OriginInfo {
	return transaction.AddressToOriginInfo(wallet.addresses[0])
}

func (wallet Wallet) MakeTransaction(quantity float64, currency string, dest transaction.Base64Address) transaction.SignedTransaction {
	unsigned := transaction.SignedTransaction{wallet.getOriginInfo(), dest, quantity, currency,
		"Sending!", nil, nil}
	return unsigned.SignMessage(&wallet.addresses[0].PrivateKey)
}

func (wallet *Wallet) UpdateBalances(blockchain blockchain.BlockChain) {
	for _, addr := range wallet.addresses {
		wallet.balance += blockchain.GetAddrBalanceFromInclusiveIndex(wallet.lastProcessedBlock+1, addr.Address, "REP")
	}
	newCurrencies := wallet.getNewCurrencyNamesAndAmts(blockchain)
	for name, amt := range newCurrencies {
		wallet.ChannelRecords[name] = GenerateNewChannelRecord(name, wallet.addresses[0].Address, blockchain, amt)
	} //TODO Don't use an arbitrary address here...

	wallet.lastProcessedBlock = int(blockchain.GetNewestBlock().Index)
	fmt.Println(wallet.GetBalances())
}

func (wallet Wallet) getNewCurrencyNamesAndAmts(chain blockchain.BlockChain) map[string]float64 {
	currencies := make(map[string]float64, 0)
	for i, block := range chain.Blocks {
		if i > wallet.lastProcessedBlock {
			for _, tx := range block.Transactions {
				if tx.SignedTrans.Currency != "REP" { //don't even bother with the rep ones we should've already processed
					if _, recordExists := wallet.ChannelRecords[tx.SignedTrans.Currency]; !recordExists { //we already know about this currency
						for _, addr := range wallet.addresses {
							if tx.SignedTrans.DestAddr == addr.Address { //we received a transaction
								fmt.Println("Received " + tx.SignedTrans.Currency)
								if _, ok := currencies[tx.SignedTrans.Currency]; ok {
									currencies[tx.SignedTrans.Currency] += tx.SignedTrans.Quantity
								} else {
									currencies[tx.SignedTrans.Currency] = tx.SignedTrans.Quantity
								}
							}
						}
					}
				}
			}
		}
	}
	return currencies
}

//func (wallet *Wallet) updateBalances(block blockchain.Block) {
//
//	for _, tx := range block.Transactions { // all transactions in the block
//		for _, addr := range wallet.addresses { // look through our personal addresses
//			if addr.Address == tx.SignedTrans.DestAddr { // if the output address matches one of ours
//				alreadyProcessed := false
//				for _, input := range wallet.inputTxIDs { // make sure we haven't recorded the transaction already
//					if input == tx.TxID {
//						alreadyProcessed = true
//					}
//				}
//				if !alreadyProcessed {
//					wallet.inputTxIDs = append(wallet.inputTxIDs, tx.TxID)
//
//					if tx.SignedTrans.Currency == "REP" {
//						wallet.balance += tx.SignedTrans.Quantity
//					}
//
//					fmt.Println("Recorded +" + strconv.FormatFloat(tx.SignedTrans.Quantity, 'f', -1, 64) +
//						" " + tx.SignedTrans.Currency + " to my wallet!")
//				}
//
//			} else if addr.Address == tx.SignedTrans.Origin.Address {
//				alreadyProcessed := false
//				for _, output := range wallet.outputTxIDs { // make sure we haven't recorded the transaction already
//					if output == tx.TxID {
//						alreadyProcessed = true
//					}
//				}
//				if !alreadyProcessed {
//					wallet.outputTxIDs = append(wallet.outputTxIDs, tx.TxID)
//
//					if tx.SignedTrans.Currency == "REP" {
//						wallet.balance -= tx.SignedTrans.Quantity
//					}
//
//					fmt.Println("Recorded -" + strconv.FormatFloat(tx.SignedTrans.Quantity, 'f', -1, 64) +
//						" " + tx.SignedTrans.Currency + " from my wallet!")
//				}
//			}
//		}
//	}
//}

package wallet

import (
	"fmt"
	"github.com/denverquane/GoBlockShare/blockchain"
	"github.com/denverquane/GoBlockShare/blockchain/transaction"
	"strconv"
)

//This is the size used for generating the personal decryption keys, NOT the channel decryption keys
const RSA_BIT_SIZE = 2048

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
	for name := range wallet.ChannelRecords {
		str += "\n                " + "1 " + name
	}

	return str
}

//TODO This is for testing!!! Don't rely on this!
func (wallet Wallet) getOriginInfo() transaction.OriginInfo {
	return transaction.AddressToOriginInfo(wallet.addresses[0])
}

func (wallet Wallet) MakeTransaction(quantity float64, dest transaction.Base64Address) transaction.SignableTransaction {
	unsigned := transaction.SignedTransaction{wallet.getOriginInfo(), dest, quantity,
		"Sending!", nil, nil}
	return transaction.Sign(&wallet.addresses[0].PrivateKey, unsigned)
}

func (wallet *Wallet) UpdateBalances(blockchain blockchain.BlockChain) {
	for _, addr := range wallet.addresses {
		wallet.balance += blockchain.GetAddrBalanceFromInclusiveIndex(wallet.lastProcessedBlock+1, addr.Address, "REP")
	}

	//newCurrencies := wallet.getNewTokenRecords(blockchain)
	//for _, v := range newCurrencies {
	//	fmt.Println(v.channelPublic)
	//	signed := v.makeTransactionForMyKey()
	//	fmt.Println("Sending back trans")
	//	//fmt.Println(signed.Payload)
	//	full := transaction.FullTransaction{signed, []string{}, ""}
	//	full.TxID = hex.EncodeToString(full.GetHash())
	//	message, added := blockchain.AddTransaction(full)
	//	if added {
	//		v.status = SentMyPubKey
	//	}
	//	fmt.Println(message)
	//	//bytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, &key, []byte(signed.Payload), []byte("key"))
	//}
	//wallet.ChannelRecords = mergeChannelMaps(wallet.ChannelRecords, newCurrencies)

	wallet.lastProcessedBlock = int(blockchain.GetNewestBlock().Index)
	fmt.Println(wallet.GetBalances())
}

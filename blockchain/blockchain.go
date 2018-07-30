package blockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/denverquane/GoBlockShare/common"
	"io/ioutil"
	"net/http"
)

type BlockChain struct {
	Blocks          []Block
	processingBlock *Block
	txMap			map[string]common.SignableTransaction
}

//IsProcessing checks the field for the block being processed, and if it is nil, indicates that the blocks for the
//chain have already been processed, and there isn't a lingering block being mined/hashed
func (chain BlockChain) IsProcessing() bool {
	return chain.processingBlock != nil
}

func (chain BlockChain) Len() int {
	return len(chain.Blocks)
}

func (chain BlockChain) ToString() string {
	str := "Chain: \n{\n"
	for _, v := range chain.Blocks {
		str += v.ToString()
	}
	str += "}"
	return str
}

func (chain BlockChain) GetTxById(txid string) common.SignableTransaction {
	if value, ok := chain.txMap[txid]; ok {
		return value
	}

	fmt.Println("Searching entire blockchain manually for: " + txid)
	for _, block := range chain.Blocks {
		for _, tx := range block.Transactions {
			if tx.TxID == txid {
				return tx
			}
		}
	}
	return common.SignableTransaction{TxID: "ERROR"}
}

//IsValid ensures that a blockchain's listed length is the same as the length of the array containing its blocks,
//and that the hashes linking blocks are valid linkages (make sure previous hash actually matches the previous block's
//hash, for example)
func (chain BlockChain) IsValid() bool {
	if chain.Len() != len(chain.Blocks) {
		return false
	}

	if chain.Len() < 2 {
		return true
	}

	for i := 0; i < chain.Len()-1; i++ {
		oldB := chain.Blocks[i]
		newB := chain.Blocks[i+1]

		if !IsBlockSequenceValid(newB, oldB) {
			return false
		}
	}
	return true
}

func (chain *BlockChain) AddTransaction(trans common.SignableTransaction, payableAddress common.Base64Address) (string, bool) {
	if chain.processingBlock != nil { //currently processing a block
		chain.processingBlock.AddTransaction(trans)
		fmt.Println("Added transaction to mining block")
		return "Added transaction to currently mining block", true
	} else {
		invalidBlock, err := GenerateInvalidBlock(chain.GetNewestBlock(), []common.SignableTransaction{trans}, payableAddress)
		if err != nil {
			return err.Error(), false
		}
		var c = make(chan bool)
		chain.processingBlock = &invalidBlock
		fmt.Println("Mining a new block")
		go chain.processingBlock.hashUntilValid(chain.GetNewestBlock().Difficulty, c)
		go chain.waitForProcessingSwap(c)
		return "Added transaction!", true
	}
}

//ProcessingReferencedTX checks to ensure that the blockchain isn't currently processing a transaction that is referred
//to by another future transaction
func (chain *BlockChain) ProcessingReferencedTX(txid string) bool {
	if chain.IsProcessing() {
		var b Block
		b = *chain.processingBlock
		for _, t := range b.Transactions {
			if t.TxID == txid {
				return true
			}
		}
	}
	return false
}

//waitForProcessingSwap waits until a block has finished mining (asynchronously) before adding it to the sequence of
//recorded/valid blocks
func (chain *BlockChain) waitForProcessingSwap(c chan bool) {
	for i := 0; !(<-c); i++ {
		if i%100000 == 0 {
			fmt.Println("Mining...")
		}
		// Wait until block is mined successfully
	}
	fmt.Println("Successfully mined block!")
	for _, v := range chain.processingBlock.Transactions {
		chain.txMap[v.TxID] = v
	}
	chain.Blocks = append(chain.Blocks, *chain.processingBlock)
	// fmt.Println(len(chain.Blocks[1].Transactions))
	chain.processingBlock = nil
}

//AreChainsSameBranch ensures that two chains are of the same structure and history, and therefore one might be a
//possible replacing chain of longer length than the other
func AreChainsSameBranch(chain1, chain2 BlockChain) bool {
	var min = 0
	if chain1.Len() > chain2.Len() {
		min = chain2.Len()
	} else {
		min = chain1.Len()
	}

	for i := 0; i < min; i++ {
		a := chain1.Blocks[i]
		b := chain2.Blocks[i]
		ah, _ := a.GetHash(true)
		bh, _ := b.GetHash(true)
		if ah != bh {
			return false
		}
	}
	return true
}

func (chain BlockChain) GetNewestBlock() Block {
	return chain.Blocks[chain.Len()-1]
}

//MakeInitialChain constructs a simple new blockchain, with an initial block paying out to the provided address
//This is a basic test to stimulate the network with an initial balance/transaction
func MakeInitialChain() BlockChain {
	chain := BlockChain{make([]Block, 1), nil, make(map[string]common.SignableTransaction)}
	chain.Blocks[0] = InitialBlock()
	return chain
}

//AppendMissingBlocks takes a chain, and appends all the transactions that are found on a longer chain to it
//This is handy when using a single Global chain that should never be entirely replaced; only appended to
func (chain BlockChain) AppendMissingBlocks(longerChain BlockChain) BlockChain {
	if AreChainsSameBranch(chain, longerChain) && longerChain.IsValid() {
		for i := len(chain.Blocks); i < len(longerChain.Blocks); i++ {
			chain.Blocks = append(chain.Blocks, longerChain.Blocks[i])
		}
	}
	return chain
}

func (chain BlockChain) GetAddressRep(addr common.Base64Address) common.ReputationSummary {
	totalSummary := common.ReputationSummary{make(map[string]common.TorrentRep, 0),
		make(map[string]common.LayerRep, 0)}

	for _, block := range chain.Blocks {
		for _, tx := range block.Transactions {
			tType := tx.TransactionType
			if tType == "TORRENT_REP" {
				torrentRep := tx.Transaction.(common.TorrentRepTrans)
				var summary common.TorrentRep
				if val, ok := totalSummary.AllTorrentRep[string(torrentRep.TorrentHash)]; ok {
					summary = val
				} else {
					summary = common.TorrentRep{0,0,0}
				}
				referencedTXID := torrentRep.TxID
				if chain.GetTxById(referencedTXID).Origin.Address == addr {
					if torrentRep.RepMessage.AccurateName {
						summary.AccurateReports++
					}
					if torrentRep.RepMessage.HighQuality {
						summary.QualityReports++
					}
					if torrentRep.RepMessage.WasValid {
						summary.ValidReports++
					}
				}
				totalSummary.AllTorrentRep[string(torrentRep.TorrentHash)] = summary
			} else if tType == "LAYER_REP" {
				layerRep := tx.Transaction.(common.LayerRepTrans)
				var summary common.LayerRep
				if val, ok := totalSummary.AllLayerRep[string(layerRep.LayerHash)]; ok {
					summary = val
				} else {
					summary = common.LayerRep{0, 0, 0}
				}
				referencedTXID := layerRep.TxID
				tx := chain.GetTxById(referencedTXID)
				if tx.Origin.Address == addr {
					if layerRep.WasLayerValid {
						summary.ValidReports++
					}
					if !layerRep.WasLayerReceived {
						summary.NotReceived++
					}
				}
				totalSummary.AllLayerRep[string(layerRep.LayerHash)] = summary
			} else if tType == "SHARED_LAYER"{
				sharedLayer := tx.Transaction.(common.SharedLayerTrans)
				var summary common.LayerRep
				if val, ok := totalSummary.AllLayerRep[string(sharedLayer.SharedLayerHash)]; ok {
					summary = val
				} else {
					summary = common.LayerRep{0, 0, 0}
				}
				if tx.Origin.Address == addr {
					summary.SharedQuantity++
				}
				totalSummary.AllLayerRep[string(sharedLayer.SharedLayerHash)] = summary
			}
		}
	}

	return totalSummary
}

func BroadcastChain(url string, chain BlockChain) {
	data, err := json.MarshalIndent(chain, "", "  ")
	//fmt.Println(string(data))
	var bytee = []byte(string(data))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bytee))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Access-Control-Allow-Origin", "*")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

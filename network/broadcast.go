package network

import (
	"chatProgram/blockchain"
	"encoding/json"
	"net/http"
	"bytes"
	"fmt"
	"io/ioutil"
)

func BroadcastChain(url string, chain blockchain.BlockChain) {
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

func BroadcastToAllPeers(peers []string, chain blockchain.BlockChain) {
	for _, v := range peers {
		BroadcastChain(v, chain)
	}
}
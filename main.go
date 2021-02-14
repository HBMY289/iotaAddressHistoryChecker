package main

import (
	"fmt"
	"net/http"

	//"net/url"
	"bytes"
	"encoding/json"
	"io/ioutil"

	//"github.com/iotaledger/iota.go/encoding/ascii"
	"github.com/iotaledger/iota.go/transaction"
	//"github.com/iotaledger/iota.go/trinary"
)

const trytesSearchURL = "https://explorer-api.iota.org/trytes/mainnet/"
const transactionSearchURL = "https://explorer-api.iota.org/transactions/mainnet/"
const addrFileName = "addressExport.txt"

func main() {

	state := ImportAddresses(addrFileName)
	fmt.Println(state[0].address)
	populateAddressTxs(state)
	//fmt.Println(state[0].txHashes)
	populateTxInfo(&state)
	//res := GetTXHashesOfAdresses(addrs)
	//fmt.Println("res:\n", res)
	//fmt.Println(addrs)
	//SearchAddresses(addrs)
}

func ImportAddresses(fileName string) []addrInfo {
	var addrs []string
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Could not find '%s'. This program requires a file with addresses to work.", fileName)
		panic(err)
	}

	err = json.Unmarshal(data, &addrs)
	if err != nil {
		fmt.Println("The file does not have the expected format.")
		panic(err)
	}
	//fmt.Println(string(data))

	return initializeStateWithAddresses(addrs[0:1]) // xxx only first one for testing
}

func initializeStateWithAddresses(addrs []string) []addrInfo {
	var state []addrInfo
	for _, addr := range addrs {
		info := addrInfo{}
		info.address = addr
		state = append(state, info)
	}
	return state
}

func SearchAddresses(addrs []string) {
	//addrs = make([]string,1)
	//addrs="HTTISYARVKRWCTUAOKH9I9VTQVTLSQXCZGPKWJN9RFISYXSXNDJIWPAAJXTYILIKHSVCEJWISMGUA9999"
	//	message := req{}
	//	message.network = "mainnet"
	//	message.hashes = []string{"HTTISYARVKRWCTUAOKH9I9VTQVTLSQXCZGPKWJN9RFISYXSXNDJIWPAAJXTYILIKHSVCEJWISMGUA9999"}
	//hashes := []string{"HTTISYARVKRWCTUAOKH9I9VTQVTLSQXCZGPKWJN9RFISYXSXNDJIWPAAJXTYILIKHSVCEJWISMGUA9999"}

	message := map[string]interface{}{
		"network": "mainnet",
		"hashes":  addrs,
	}
	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	fmt.Println("message: ", message)
	fmt.Println("json: ", string(bytesRepresentation))
	response, err := http.Post(trytesSearchURL, "application/json", bytes.NewBuffer(bytesRepresentation))

	//okay, moving on...
	if err != nil {
		//handle postform error
		panic(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		panic(err)
	}

	fmt.Printf("body: %s\n", string(body))
}

func GetTxsTryteResponse(hashes []string) txTrytesResponse {
	var result txTrytesResponse
	message := map[string]interface{}{
		"network": "mainnet",
		"hashes":  hashes,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}

	response, err := http.Post(trytesSearchURL, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		panic(err)
	}

	//fmt.Printf("body: %s\n", string(body))

	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}

	return result
}

func populateTxInfo(state *[]addrInfo) {
	for j, _ := range *state {
		resp := GetTxsTryteResponse(*(state)[j].txHashes)
		for i, trytes := range resp.Trytes {
			tx, err := transaction.AsTransactionObject(trytes)
			if err != nil {
				panic(err)
			}
			if tx.Value != 0 {
				bundle := &(getBundle(&state[j].valueBundles, tx.Bundle))
				valueTx := txInfo{}
				valueTx.attachmentTimestamp = tx.AttachmentTimestamp
				if resp.MilestoneIndexes[i] > 0 {
					valueTx.confMilestone = resp.MilestoneIndexes[i]
					valueTx.confirmed = true
				}
				valueTx.hash = tx.Hash
				valueTx.messageASCII = tx.SignatureMessageFragment //TODO decode
				valueTx.obsoleteTag = tx.ObsoleteTag
				valueTx.rawTrytes = trytes
				valueTx.tag = tx.Tag
				valueTx.timestamp = tx.Timestamp
				valueTx.value = tx.Value
				//bundle.txs = append(bundle.txs, valueTx)
				*bundle.txs = append(*bundle.txs, valueTx)
			}
			//fmt.Println("t0:\n", ts[0].Hash, ts[0].Bundle, ts[0].Confirmed)
		}
		//fmt.Println(i, resp)
		//state[i].txInfos = resp
		//state[i].  //xxx
	}
}

func getBundle(bundles *[]valueBundle, hash string) valueBundle {
	for _, bundle := range *bundles {
		if bundle.bundleHash == hash {
			return bundle
		}
	}
	newBundle := valueBundle{}
	newBundle.bundleHash = hash
	*bundles = append(*bundles, newBundle)
	return (*bundles)[len(*bundles)-1]
}

func populateAddressTxs(state []addrInfo) {

	for i, addr := range state {
		state[i].txHashes = GetTxHashesOfAddress(addr.address)
	}

}

func GetTxHashesOfAddress(addr string) []string {
	addrHashes := addrTxResponse{}
	resp, err := http.Get(transactionSearchURL + addr)
	if err != nil {
		panic(err)
	}

	if resp != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		//fmt.Println(string(body))
		addrHashes = addrTxResponse{}
		err = json.Unmarshal(body, &addrHashes)
		if err != nil {
			panic(err)
		}

		resp.Body.Close()
	} else {
		fmt.Println("BAD")
	}

	return addrHashes.Hashes
}

type req struct {
	network string
	hashes  []string
}

type addrTxResponse struct {
	Mode   string   `json:"mode"`
	Hashes []string `json:"hashes"`
	Cursor struct {
		Node    int  `json:"node"`
		HasMore bool `json:"hasMore"`
		Perma   int  `json:"perma"`
		NextInt int  `json:"nextInt"`
	} `json:"cursor"`
}

type txTrytesResponse struct {
	MilestoneIndexes []int64  `json:"milestoneIndexes"`
	Trytes           []string `json:"trytes"`
}

type addrInfo struct {
	address      string
	txHashes     []string
	txInfos      []txInfo
	bundelHashes []string
	valueBundles []valueBundle
	balance      uint64
}

type valueBundle struct {
	bundleHash string
	txs        []txInfo
}

type txInfo struct {
	hash                string
	value               int64
	attachmentTimestamp int64
	timestamp           uint64
	confMilestone       int64
	confirmed           bool
	tag                 string
	obsoleteTag         string
	messageASCII        string
	bundlehash          string
	rawTrytes           string
}

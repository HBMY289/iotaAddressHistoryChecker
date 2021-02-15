package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	//"github.com/iotaledger/iota.go/encoding/ascii"
	"github.com/HBMY289/iotaAddressHistoryChecker/analysis"
	. "github.com/HBMY289/iotaAddressHistoryChecker/types"
	"github.com/iotaledger/iota.go/transaction"
)

const trytesSearchURL = "https://explorer-api.iota.org/trytes/mainnet/"
const transactionSearchURL = "https://explorer-api.iota.org/transactions/mainnet/"

const addrFileName = "addressExport.txt"
const stateFileName = "stateExport.txt"

func main() {
	state := StateInfo{}
	importStateFromFile(&state, stateFileName)
	//state.Addresses = state.Addresses[0:1] // TODO only use first address for testing
	//fmt.Println(state.Addresses[0].TxInfos[0])
	bundles := analysis.GetConfirmedBundles(state)
	fmt.Println("found bundles: ", len(bundles))
	fmt.Println("bundles[0]:", bundles[0].TxInfos[0].AttachmentTimestamp)
	for _, bundle := range bundles {
		fmt.Println(bundle.TxInfos[0].AttachmentTimestamp, bundle.TxInfos[0].Bundlehash)
	}
	fmt.Println(analysis.GetAnalyzedBundlesReport(bundles))

}

func main2() {
	state := StateInfo{}
	importAddressesFromFile(&state, addrFileName)
	//state.Addresses = state.Addresses[0:1] // TODO only use first address for testing
	debug(state.Addresses[0].Address)
	populateAddressInfo(&state)
	debug(state.Addresses[0].BundleHashes)
	fmt.Println(state.Addresses[0].TxInfos[0])

	exportState(state, stateFileName)
}

func exportState(state StateInfo, fileName string) {
	fmt.Println(len(state.Addresses))
	fmt.Println(state.Addresses[0].TxInfos[0])
	j, err := json.Marshal(state)
	if err != nil {
		fmt.Println(err)
		return
	}
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("j: ", string(j))
	l, err := f.WriteString(string(j))
	if err != nil {
		fmt.Println(err)
		fmt.Println(l)
		f.Close()
		return
	}
}

func importStateFromFile(state *StateInfo, fileName string) {

	data, err := ioutil.ReadFile(fileName)
	var tempState StateInfo
	if err != nil {
		fmt.Printf("Could not find '%s'. This program requires a file with addresses to work.", fileName)
		panic(err)
	}

	err = json.Unmarshal(data, &tempState)
	if err != nil {
		fmt.Println("The file does not have the expected format.")
		panic(err)
	}

	*state = tempState
	fmt.Printf("successfully imported state with %d addresses from file: %s\n", len(state.Addresses), fileName)
}

func importAddressesFromFile(state *StateInfo, fileName string) {
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

	for _, addr := range addrs {
		info := AddrInfo{}
		info.Address = addr
		state.Addresses = append(state.Addresses, info)
	}

	fmt.Printf("successfully imported %d addresses from file: %s\n", len(state.Addresses), fileName)
}

func GetTxsTryteResponse(hashes []string) TxTrytesResponse {
	var result TxTrytesResponse
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

	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}

	return result
}

func populateAddressInfo(state *StateInfo) {
	for i, _ := range state.Addresses {

		fmt.Printf("\rgetting info for address (%d/%d)", i+1, len(state.Addresses))
		txHashes := findTxHashes(state.Addresses[i].Address)
		debug(txHashes)
		bundleHashes := getBundleHashes(txHashes)
		state.Addresses[i].BundleHashes = bundleHashes
		debug("bundle hashes:")
		debug(state.Addresses[i].BundleHashes)
		debug("getting all tx hashes of bundles")
		txHashes = getAllTxHashesOfBundles(state.Addresses[i].BundleHashes)
		state.Addresses[i].TxHashes = txHashes
		debug("all tx hashes:")
		debug(state.Addresses[i].TxHashes)
		debug("getting tx info for txs")
		state.Addresses[i].TxInfos = getValueTxs(txHashes)
		debug("txInfos:")
		debug(state.Addresses[i].TxInfos)
		//TODO get balances
	}
}

func getValueTxs(txHashes []string) []TxInfo {
	var txInfos []TxInfo
	resp := GetTxsTryteResponse(txHashes)
	for i, trytes := range resp.Trytes {
		tx, err := transaction.AsTransactionObject(trytes)
		if err != nil {
			panic(err)
		}
		if tx.Value != 0 {
			valueTx := TxInfo{}
			valueTx.AttachmentTimestamp = tx.AttachmentTimestamp
			if resp.MilestoneIndexes[i] > 0 {
				valueTx.ConfMilestone = resp.MilestoneIndexes[i]
				valueTx.Confirmed = true
			}
			valueTx.Hash = tx.Hash
			valueTx.MessageASCII = tx.SignatureMessageFragment //TODO decode
			valueTx.ObsoleteTag = tx.ObsoleteTag
			//valueTx.rawTrytes = trytes
			valueTx.Tag = tx.Tag
			valueTx.Timestamp = tx.Timestamp
			valueTx.Value = tx.Value
			valueTx.Bundlehash = tx.Bundle
			valueTx.Address = tx.Address
			txInfos = append(txInfos, valueTx)
		}
	}
	return txInfos
}

func getAllTxHashesOfBundles(bHashes []string) []string {
	var hashes []string
	for _, bHash := range bHashes {
		debug("finding hashes for bundle: " + bHash)
		txHashes := findTxHashes(bHash)
		hashes = append(hashes, txHashes...)
		debug(txHashes)
	}
	return hashes
}

func findTxHashes(hash string) []string {
	hashes := FindTxResponse{}
	resp, err := http.Get(transactionSearchURL + hash)
	if err != nil {
		panic(err)
	}
	if resp != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(body, &hashes)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()
	} else {
		panic(err)
	}
	return hashes.Hashes
}

func getBundleHashes(txHashes []string) []string {
	var hashes []string
	resp := GetTxsTryteResponse(txHashes)
	for _, trytes := range resp.Trytes {
		tx, err := transaction.AsTransactionObject(trytes)
		if err != nil {
			panic(err)
		}
		if tx.Value != 0 {
			if !known(hashes, tx.Bundle) {
				hashes = append(hashes, tx.Bundle)
			}
		}
	}
	return hashes
}

func known(names []string, newName string) bool {
	for _, name := range names {
		if name == newName {
			return true
		}
	}
	return false
}

func GetTxHashesOfAddress2(addr string) []string {
	addrHashes := FindTxResponse{}
	resp, err := http.Get(transactionSearchURL + addr)
	if err != nil {
		panic(err)
	}

	if resp != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		addrHashes = FindTxResponse{}
		err = json.Unmarshal(body, &addrHashes)
		if err != nil {
			panic(err)
		}

		resp.Body.Close()
	} else {
		panic(err)
	}

	return addrHashes.Hashes
}

func debug(item interface{}) {
	//fmt.Println("debug:", item)
}

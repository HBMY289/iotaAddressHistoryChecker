package explorer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/HBMY289/iotaAddressHistoryChecker/types"
	"github.com/iotaledger/iota.go/converter"
	"github.com/iotaledger/iota.go/transaction"
)

const trytesSearchURL = "https://explorer-api.iota.org/trytes/legacy-mainnet/"
const transactionSearchURL = "https://explorer-api.iota.org/transactions/legacy-mainnet/"
const balanceSearchURL = "https://explorer-api.iota.org/address/legacy-mainnet/"
const network = "legacy-mainnet"

func PopulateAddressInfo(state *StateInfo) error {
	for i, _ := range state.Addresses {

		fmt.Printf("\rgetting info for address (%d/%d)", i+1, len(state.Addresses))
		txHashes, err := findTxHashes(state.Addresses[i].Address)
		if err != nil {
			return err
		}
		bundleHashes, err := getBundleHashes(txHashes)
		if err != nil {
			return err
		}
		state.Addresses[i].BundleHashes = bundleHashes

		txHashes, err = getAllTxHashesOfBundles(state.Addresses[i].BundleHashes)
		if err != nil {
			return err
		}
		state.Addresses[i].TxHashes = txHashes

		state.Addresses[i].TxInfos, err = getValueTxs(txHashes)
		if err != nil {
			return err
		}

		state.Addresses[i].Balance, err = getBalance(state.Addresses[i].Address)
		if err != nil {
			return err
		}

	}
	return nil
}

func findTxHashes(hash string) ([]string, error) {
	hashes := FindTxResponse{}
	resp, err := http.Get(transactionSearchURL + hash)
	if err != nil {
		return []string{}, err
	}
	if resp != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return []string{}, err
		}
		err = json.Unmarshal(body, &hashes)
		if err != nil {
			return []string{}, err
		}
		resp.Body.Close()
	} else {
		return []string{}, err
	}
	return hashes.Hashes, nil
}

func getBalance(addr string) (uint64, error) {
	var bResp BalanceResponse
	resp, err := http.Get(balanceSearchURL + addr)
	if err != nil {
		return 0, err
	}
	if resp != nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return 0, err
		}
		err = json.Unmarshal(body, &bResp)
		if err != nil {
			return 0, err
		}
		resp.Body.Close()
	} else {
		return 0, err
	}
	return bResp.Balance, nil
}

func getBundleHashes(txHashes []string) ([]string, error) {
	var hashes []string
	resp, err := GetTxsTryteResponse(txHashes)
	if err != nil {
		return []string{}, err
	}
	for _, trytes := range resp.Trytes {
		tx, err := transaction.AsTransactionObject(trytes)
		if err != nil {
			return []string{}, err
		}
		if tx.Value != 0 {
			if !known(hashes, tx.Bundle) {
				hashes = append(hashes, tx.Bundle)
			}
		}
	}
	return hashes, nil
}

func known(names []string, newName string) bool {
	for _, name := range names {
		if name == newName {
			return true
		}
	}
	return false
}

func getAllTxHashesOfBundles(bHashes []string) ([]string, error) {
	var hashes []string
	for _, bHash := range bHashes {
		txHashes, err := findTxHashes(bHash)
		if err != nil {
			return []string{}, err
		}
		hashes = append(hashes, txHashes...)
	}
	return hashes, nil
}

func getValueTxs(txHashes []string) ([]TxInfo, error) {
	var txInfos []TxInfo
	resp, err := GetTxsTryteResponse(txHashes)
	if err != nil {
		return []TxInfo{}, err
	}
	for i, trytes := range resp.Trytes {
		tx, err := transaction.AsTransactionObject(trytes)
		if err != nil {
			return txInfos, err
		}
		if tx.Value != 0 {
			valueTx := TxInfo{}
			valueTx.AttachmentTimestamp = tx.AttachmentTimestamp
			if resp.MilestoneIndexes[i] > 0 {
				valueTx.ConfMilestone = resp.MilestoneIndexes[i]
				valueTx.Confirmed = true
			}
			valueTx.Hash = tx.Hash
			valueTx.MessageASCII = trytesToAscii(tx.SignatureMessageFragment)

			valueTx.ObsoleteTag = tx.ObsoleteTag
			valueTx.Tag = tx.Tag
			valueTx.Timestamp = tx.Timestamp
			valueTx.Value = tx.Value
			valueTx.Bundlehash = tx.Bundle
			valueTx.Address = tx.Address
			txInfos = append(txInfos, valueTx)
		}
	}
	return txInfos, nil
}

func GetTxsTryteResponse(hashes []string) (TxTrytesResponse, error) {
	var result TxTrytesResponse
	message := map[string]interface{}{
		"network": network,
		"hashes":  hashes,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		return result, err
	}

	response, err := http.Post(trytesSearchURL, "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return result, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func trytesToAscii(trytes string) string {
	if len(trytes)%2 != 0 {
		trytes = trytes + "9"
	}
	message, err := converter.TrytesToASCII(trytes)
	if err != nil {
		panic(err)
	}
	return message
}

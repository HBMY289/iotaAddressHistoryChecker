package analysis

import (
	"fmt"
	"sort"

	. "github.com/HBMY289/iotaAddressHistoryChecker/types"
)

func GetConfirmedBundles(state StateInfo) []ValueBundle {
	bundles := collectBundles((state))
	sort.Slice(bundles, func(i1, i2 int) bool {
		return bundles[i1].TxInfos[0].AttachmentTimestamp < bundles[i2].TxInfos[0].AttachmentTimestamp
	})
	return bundles
}

func analyzeBundles(bundles []ValueBundle) {

}

func collectBundles(state StateInfo) []ValueBundle {

	var bundles []ValueBundle
	for _, addr := range state.Addresses {
		for _, tx := range addr.TxInfos {
			if tx.Confirmed {
				bundles = addTx(bundles, tx)
			}
		}
	}
	return bundles
}

func addTx(bundles []ValueBundle, tx TxInfo) []ValueBundle {
	bIndex := getIndexOfBundle(bundles, tx.Bundlehash)
	if bIndex == -1 {
		bundles = append(bundles, ValueBundle{})
		bIndex = len(bundles) - 1
	}
	tIndex := getIndexOfTx(bundles[bIndex].TxInfos, tx.Hash)
	if tIndex == -1 {
		bundles[bIndex].Hash = tx.Bundlehash
		bundles[bIndex].TxInfos = append(bundles[bIndex].TxInfos, tx)
	}
	return bundles
}

func getIndexOfBundle(bundles []ValueBundle, hash string) int {
	for i, bundle := range bundles {
		if bundle.Hash == hash {
			fmt.Println("found match")
			return i
		}
	}
	return -1
}

func getIndexOfTx(txs []TxInfo, hash string) int {
	for i, tx := range txs {
		if tx.Hash == hash {
			return i
		}
	}
	return -1
}

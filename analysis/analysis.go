package analysis

import (
	"fmt"
	"sort"
	"strings"
	"time"

	. "github.com/HBMY289/iotaAddressHistoryChecker/types"
)

func getConfirmedBundles(state StateInfo) []ValueBundle {
	bundles := collectBundles((state))
	sort.Slice(bundles, func(i1, i2 int) bool {
		return bundles[i1].TxInfos[0].AttachmentTimestamp < bundles[i2].TxInfos[0].AttachmentTimestamp
	})
	bundles = analyzeBundles(bundles, seedAddresses(state))
	return bundles
}

func seedAddresses(state StateInfo) []string {
	var addrs []string
	for _, addrInfo := range state.Addresses {
		addrs = append(addrs, addrInfo.Address)
	}
	return addrs
}

func analyzeBundles(bundles []ValueBundle, addrs []string) []ValueBundle {
	for i, bundle := range bundles {
		bundles[i] = analyzeBundle(bundle, addrs)
	}
	return bundles

}

func GetBalanceReport(state StateInfo) string {
	report := "\nThe following balances could be found in the data of explorer.iota.org:\n"
	var total uint64
	for _, addr := range state.Addresses {
		report += fmt.Sprintf("%di on address %s\n", addr.Balance, addr.Address)
		total += addr.Balance
	}
	report += fmt.Sprintf("\ntotal balance: %di for %d addresses\n", total, len(state.Addresses))
	return report
}

func GetAnalyzedBundlesReport(state StateInfo) string {
	var report string
	var plural string

	bundles := getConfirmedBundles(state)
	report = "\nThe following value movements could be found in the data of explorer.iota.org:\n"
	for _, bundle := range bundles {
		if !bundle.Internal {
			if len(bundle.Addresses) > 1 {
				plural = "es"
			}
			report += fmt.Sprintf("%s\n%di %s address%s %s\nvia bundle %s\n\n", bundle.Date, bundle.Value, direction(bundle.Outgoing), plural, strings.Join(bundle.Addresses, "\nand "), bundle.Hash)
		}
	}
	return report
}

func direction(out bool) string {
	if out {
		return "sent to"
	}
	return "received from"
}

func analyzeBundle(bundle ValueBundle, addrs []string) ValueBundle {

	var knownInputs, unknownInputs, knownOutputs, unknownOutputs []string
	var knownBal, unknownBal int64
	for _, txInfo := range bundle.TxInfos {
		if contains(addrs, txInfo.Address) {
			if txInfo.Value < 0 {
				knownInputs = append(knownInputs, txInfo.Address)
			} else {
				knownOutputs = append(knownOutputs, txInfo.Address)
			}
			knownBal += txInfo.Value
		} else {
			if txInfo.Value < 0 {
				unknownInputs = append(unknownInputs, txInfo.Address)
			} else {
				unknownOutputs = append(unknownOutputs, txInfo.Address)
			}
			unknownBal += txInfo.Value
		}
	}

	switch {
	case knownBal > 0:
		bundle.Outgoing = false
		bundle.Addresses = unknownInputs
		bundle.Value = knownBal
	case knownBal < 0:
		bundle.Outgoing = true
		bundle.Addresses = unknownOutputs
		bundle.Value = knownBal * -1
	case knownBal == 0:
		bundle.Internal = true
		bundle.Value = 0
	}
	bundle.Date = time.Unix(bundle.TxInfos[0].AttachmentTimestamp/1000, 0).Format("Mon, 02 Jan 2006 15:04:05 MST")
	return bundle
}

func contains(items []string, matchItem string) bool {
	for _, item := range items {
		if item == matchItem {
			return true
		}
	}
	return false
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

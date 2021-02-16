package main

import (
	"fmt"
	"os"

	//"github.com/iotaledger/iota.go/encoding/ascii"
	"github.com/HBMY289/iotaAddressHistoryChecker/analysis"
	. "github.com/HBMY289/iotaAddressHistoryChecker/data"
	"github.com/HBMY289/iotaAddressHistoryChecker/explorer"
	. "github.com/HBMY289/iotaAddressHistoryChecker/types"

	"github.com/HBMY289/iotaZeroBalanceHelper/userIO"
)

const addrFileName = "addressExport.txt"
const stateFileName = "stateExport.txt"

func main() {
	state := StateInfo{}
	mainMenu(state)
}

func mainMenu(state StateInfo) {
	fmt.Println("This program aggregates and analyses information about Iota token movements on your account.")
	fmt.Println("It requires a list of known addresses of your seed and will then request the transaction history from the Iota tangle explorer 'explorer.iota.org'.")
	fmt.Println("The required address file can be generated using the iotaZeroBalanceHelper ('https://github.com/HBMY289/iotaZeroBalanceHelper').")
	for {
		opt := userIO.GetOption("You can either start by importing an address file from the iotaZeroBalanceHelper or a previously exported state of this program.",
			[]string{"Import address file", "Import saved state file", "Exit"})
		switch opt {
		case 1:
			err := ImportAddressesFromFile(&state, addrFileName)
			if err == nil {
				fmt.Println("Make sure the computer is connected to the internet. The transaction history will now be requested from explorer.iota.org.")
				userIO.WaitforEnter()
				err := explorer.PopulateAddressInfo(&state)
				if err == nil {
					export := userIO.GetConfirmation("\nDo you want to export the collected information as a state file for future use?")
					if export {
						err := ExportState(state, stateFileName)
						if err != nil {
							fmt.Println("Could not export state file.")
						}
					}
					analysisMenu(state)
				}
			}
		case 2:
			err := ImportStateFromFile(&state, stateFileName)
			if err == nil {
				analysisMenu(state)
			}
		case 3:
			os.Exit(0)
		}
	}
}

func analysisMenu(state StateInfo) {
menu:
	for {
		opt := userIO.GetOption("\nThe imported state file can now be analyzed. Choose an option:", []string{"Show balances of addresses", "Show value movements by date", "Back"})
		switch opt {
		case 1:
			fmt.Println(analysis.GetBalanceReport(state))
			userIO.WaitforEnter()
		case 2:
			fmt.Println(analysis.GetAnalyzedBundlesReport(state))
			userIO.WaitforEnter()
		case 3:
			break menu
		}
	}
}

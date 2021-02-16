package data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/HBMY289/iotaAddressHistoryChecker/types"
)

func ImportStateFromFile(state *StateInfo, fileName string) error {

	var tempState StateInfo
	text, err := getTextFromFile(fileName)
	if err != nil {
		return err
	}
	err = json.Unmarshal(text, &tempState)
	if err != nil {
		fmt.Println("The file does not have the expected format.")
		return err
	}

	*state = tempState
	fmt.Printf("\n\nSuccessfully imported state with %d addresses from file: %s\n", len(state.Addresses), fileName)
	return nil
}

func ImportAddressesFromFile(state *StateInfo, fileName string) error {
	var addrs []string
	text, err := getTextFromFile(fileName)
	if err != nil {
		return err
	}
	err = json.Unmarshal(text, &addrs)
	if err != nil {
		fmt.Println("The file does not have the expected format.")
		return err
	}

	for _, addr := range addrs {
		info := AddrInfo{}
		info.Address = addr
		state.Addresses = append(state.Addresses, info)
	}

	fmt.Printf("\n\nSuccessfully imported %d addresses from file: %s\n", len(state.Addresses), fileName)
	return nil
}

func getTextFromFile(fileName string) ([]byte, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Could not open '%s'.Make sure the file has the correct name and is placed in the same folder as this program.", fileName)
		return data, err
	}
	return data, nil
}

func ExportState(state StateInfo, fileName string) error {

	j, err := json.Marshal(state)
	if err != nil {
		panic(err)
	}
	err = writeToFile(string(j),fileName)
	if err != nil {
		return err
	}
	return nil
}

func writeToFile (text,fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(text)
	if err != nil {
		return err
	}
	return nil
}
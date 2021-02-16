# iotaAddressHistoryChecker

# What is it for?
This program was written to help Iota users who unexpectedly see a zero balance in their Trinity wallet. There are multiple reasons why this could be happening and often it is not easy for non-experts to find out what is going on.

## What does it do?
The program automatically downloads the complete tranasaction history for all known addresses of your account. It the analyzes them and creates a short report listing all relevant information about token movements for incoming and outgoing tranasctions.
While this program requires an internet connection to work it does not need the seed for your account. Instead it relies on a list of known addresses that can be securely generated on an air-gapped computer using another program I have written ([iotaZeroBalanceHelper](https://github.com/HBMY289/iotaZeroBalanceHelper).


## How it works


## Disclaimer
NEVER share your seed with anyone. No Iota community member or member of the Iota Foundation will ever ask for your seed. If someone does it is 100% a scam to steal your money. 


## How to start the tool
The simplest way is to download the appropriate binary executable for your operating system from [releases](https://github.com/HBMY289/iotaAddressHistoryChecker/releases). You can also build the tool from source, which is rather easy as well. Assuming you have [go](https://golang.org/doc/install) and [git](https://www.atlassian.com/git/tutorials/install-git) installed already you can just execute this command for example in your user folder to get a copy of the source code.
```
git clone https://github.com/HBMY289/iotaAddressHistoryChecker.git
```

Then you change into the new folder and build the excutable.
```
cd iotaAddressHistoryChecker
go build
```
After that you can just start the newly created binary file by typing
```
./iotaAddressHistoryChecker
```
or on Windows
```
iotaAddressHistoryChecker.exe
```
## How to use the tool
Once the program is running you will have to import the address file generated earlier by the ([iotaZeroBalanceHelper](https://github.com/HBMY289/iotaZeroBalanceHelper). Place the file `addressExport.txt` next to this tool's executable. After the successful import all available transaction information will be downloaded.

##### Export state
Depending on the number of supplied addresses the download of all related transactions can take quite a while. The program allows to export this aggregated information to a state file so it can later be used and analyzed without having to request the information from the explorer again.

## Need additonal help?
If you need any additonal help either with the tool itself or with checking the exported addresses you can contact me (HBMY289) or any other community member in the #help channel via the official via the official [Iota Discord server](https://discord.iota.org/).

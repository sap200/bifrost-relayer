package main

import (
	"fmt"
	"os/exec"
	"bytes"
	"github.com/sap200/bifrost-relayer/tzpengine"
	"github.com/sap200/bifrost-relayer/commands"
	"time"
	"strings"
	"sync"
	gop "github.com/sap200/bifrost-relayer/goperations"
	"github.com/sap200/bifrost-relayer/verifier"
)

var contract string

func init() {
	commands.ContractAddress, commands.FA12 = commands.SetFa12AndBifrostContractAddress("../tezos/contract_addr/bifrost.txt", "../tezos/contract_addr/fa12.txt")
	contract = commands.ContractAddress
}


var wg sync.WaitGroup

func main() {

	printHeaders()
	// Doesn't works in dev mode
	// InitializeKey()
	wg.Add(4)
	// Run the engine to constantly listen on storage and rly txs...
	go runBifrostToCosmosEngine()
	// Run the verifier engine to check for malicious txs...
	go verifier.LaunchVerificationChannel(commands.VerificationPort)
	// Run the operation engine
	go gop.GeneralOperationChannel(commands.GPPORT)

	go runFa12RelayEngine()

	wg.Wait()
	
}

func printHeaders() {
	fmt.Println()
	fmt.Println("				 /--------/")
	fmt.Println("				/--------/")
	fmt.Println("				  / /")
	fmt.Println("				 / /")
	fmt.Println("				/ /")
	fmt.Println("----------------------------------------------------")
	fmt.Println("starting bifrost relayer and verification engine...")
	fmt.Println("----------------------------------------------------")
	fmt.Println("[tezos_chain] ---> [bifrostzone] @via bifrostzone...")
	fmt.Println("[tezos_chain] <--- [bifrostzone] @via bifrostzone...")
	fmt.Println("----------------------------------------------------")
	fmt.Println()
}

// TODO: Add keys at both end for relayer.

func InitializeKey() {
	var outb, errb bytes.Buffer
	qkey := commands.QueryCosmosKeyPair()
	cmd := exec.Command(qkey[0], qkey[1:]...)
	cmd.Stderr = &errb
	cmd.Stdout = &outb
	cmd.Run()


	if &errb != nil && strings.Contains(errb.String(), commands.KeyDoesNotExists) {
		// make the key else key already exists....
		var errb1, outb1 bytes.Buffer
		addCmd := commands.CreateKeyPair()
		cmd = exec.Command(addCmd[0], addCmd[1:]...)
		cmd.Stdout = &outb1
		cmd.Stderr = &errb1
		cmd.Run()
		
		commands.WriteKey(commands.CosmosKeyPath, append(outb1.Bytes(), errb1.Bytes()...))
	}
}

func runBifrostToCosmosEngine() {
	defer wg.Done()
	for { 

		fetchStorage := commands.FetchTezosStorage(contract)
		out, err := exec.Command(fetchStorage[0], fetchStorage[1:]...).Output()
		if err != nil {
			fmt.Println(err)
		}

		x := tzpengine.GetTxInSlice(string(out))
		txs := tzpengine.GetCosmosTxs(x)
		
		for _, tx := range txs {
			// execute peg-zone command
			fmt.Println("\u2714 Relaying 1 transmit packet [tezos] ---> [bifrostzone] @via {{bifrost}}...")
			var outbt, errbt bytes.Buffer
			pegcmd := commands.RelayPacketCommand(tx, "alice")
			cmdxs := exec.Command(pegcmd[0], pegcmd[1:]...)
			cmdxs.Stdout = &outbt
			cmdxs.Stderr = &errbt
			cmdxs.Run()

			// fmt.Println(outbt.String())
			// fmt.Println(errbt.String())

			if &errbt != nil && (strings.Contains(errbt.String(), commands.VerificationFail) ||
				strings.Contains(errbt.String(), "Error:")) || 
			    &outbt != nil && strings.Contains(outbt.String(), `"code":1`) {

				fmt.Println("\u2718 \u269D Relaying 1 acknowledgement-Failed packet [tezos] <--- [bifrostzone] @via {{bifrost}}..")
				// revert back the txs
				var outud, errud bytes.Buffer
				updateCmd := commands.UpdateStatus(contract, "alice", tzpengine.Failed, tx.TxId)
				cmdud := exec.Command(updateCmd[0], updateCmd[1:]...)
				cmdud.Stdout = &outud
				cmdud.Stderr = &errud
				cmdud.Run()
				// fmt.Println(outud.String())
				// fmt.Println(errud.String())
			}
			

			if &outbt != nil && strings.Contains(outbt.String(), `"code":0`) {
				updateCmd := commands.UpdateStatus(contract, "alice", tzpengine.Success, tx.TxId)
				var outub, errub bytes.Buffer
				cmdxs = exec.Command(updateCmd[0], updateCmd[1:]...)
				cmdxs.Stdout = &outub
				cmdxs.Stderr = &errub
				cmdxs.Run()
				fmt.Println("\u2714 \u269D Relaying 1 acknowledgement-Success packet [tezos] <--- [bifrostzone] @via {{bifrost}}..")
			}

			fmt.Println()
			time.Sleep(time.Second * 2)
		}

		time.Sleep(time.Second*20)
	}
}

func runFa12RelayEngine() {
	for {
		// testing....
		fetchStorage := commands.FetchTezosStorage(commands.FA12)
		var outf, errf bytes.Buffer
		fcmd := exec.Command(fetchStorage[0], fetchStorage[1:]...)
		fcmd.Stdout = &outf
		fcmd.Stderr = &errf
		fcmd.Run()

		x := tzpengine.GetTxInSlice(outf.String())
		txs, err := tzpengine.GetFa12Txs(x)
		if err != nil {
			fmt.Println(err)
		}

		for _, tx := range txs {
			// relay this tx...
			var outbt, errbt bytes.Buffer
			fmt.Println("\u2714 Relaying 1 transmit packet [tezos] ---> [bifrostzone] @via {{bifrost}}...")
			pegcmd := commands.RelayFa12PacketCommand(tx, "alice")
			cmdxs := exec.Command(pegcmd[0], pegcmd[1:]...)
			cmdxs.Stdout = &outbt
			cmdxs.Stderr = &errbt
			cmdxs.Run()

			if (&errbt != nil && (strings.Contains(errbt.String(), commands.Fa12Error) || 
				strings.Contains(errbt.String(), commands.VerificationFail))) || 
				(&outbt != nil && strings.Contains(outbt.String(), `"code:1"`)) {

				// this is when transaction fails 
				// send a update status with failed message 
				fmt.Println("\u2718 \u269D Relaying 1 acknowledgement-Failed packet [tezos] <- [bifrostzone] @via {{bifrost}}...")
				var errut, outut bytes.Buffer
				updateCmd := commands.UpdateFa12Status(commands.FA12, "alice", tzpengine.Failed, tx.TxId)
				cmdud := exec.Command(updateCmd[0], updateCmd[1:]...)
				cmdud.Stdout = &outut
				cmdud.Stderr = &errut
				cmdud.Run()
	
			}

			if &outbt != nil && strings.Contains(outbt.String(), `"code":0`) {
				fmt.Println("\u2714 \u269D Relaying 1 acknowledgement-Success packet [tezos] <--- [bifrostzone] @via {{bifrost}}..")
				var errut, outut bytes.Buffer
				updateCmd := commands.UpdateFa12Status(commands.FA12, "alice", tzpengine.Success, tx.TxId)
				cmdud := exec.Command(updateCmd[0], updateCmd[1:]...)
				cmdud.Stdout = &outut
				cmdud.Stderr = &errut
				cmdud.Run()
			}
			fmt.Println()
			time.Sleep(time.Second * 2)
		}
		time.Sleep(time.Second * 20)
	}
}
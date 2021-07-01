package goperations

import (
	"net"
	"fmt"
	"os"
	"bufio"
	"strings"
	"io"
	"github.com/sap200/bifrost-relayer/commands"
	"bytes"
	"os/exec"
	"strconv"
	"encoding/json"
	"github.com/sap200/bifrost-relayer/tzpengine"
)

func GeneralOperationChannel(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		// handle error
		fmt.Println("\nunable to launch general-operation channel...")
		os.Exit(1)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			fmt.Println("\n\u2718 Relaying 1 Operation-Failed packet [tezos] <--- [bifrostzone] @via {{bifrost}}...")
			fmt.Println(err)
		}
	go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("\n\u2714 Relaying 1 Operation-Perform packet [tezos] <--- [bifrostzone] @via {{bifrost}}...")
	// 1st read what is said in the conn...
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("\nUnable to read opcode.")
		return
	}
	instruction := strings.TrimSpace(string(status))
	instArr := strings.Split(instruction, " ")
	opcode := instArr[0]
	// perform operation:
	// option a: burn
	switch opcode {
	case "burn":
		// fmt.Println(instArr)
		// do operations like 
		if len(instArr) != 3 {
			io.WriteString(conn, FAIL + " :length of instruction > 3\n")
			break
		}

		address := instArr[1]
		amt := instArr[2]
		amount, err := strconv.ParseUint(amt, 10, 64)
		if err != nil {
			io.WriteString(conn, FAIL + " :unable to convert amount to uint\n")
			break
		}
		isFine, _, _:= unlockTezos(address, amount)
		// fmt.Println(outer)
		// fmt.Println(errorss)
		fmt.Println("\u2714 Relaying 1 Operation-Result packet [tezos] ---> [bifrostzone] @via {{bifrost}}...")
		if isFine {
			io.WriteString(conn, SUCCESS + " :Seems like there is no error.\n")
		} else {
			io.WriteString(conn, FAIL + " :Seems like there is error in unlocking tezos\n")
		}
	case "mint":
		if len(instArr) != 3 {
			io.WriteString(conn, FAIL + ": length of instruction != 3\n")
			break
		}

		address := instArr[1]
		amt := instArr[2]
		amount, err := strconv.ParseUint(amt, 10, 64)
		if err != nil {
			io.WriteString(conn, FAIL + ": unable to convert amount to uint\n")
			break
		}

		isFine, _, _ := mintFA12(address, amount)
		fmt.Println("\u2714 Relaying 1 Operation-Result packet [tezos] ---> [bifrostzone] @via {{bifrost}}...")
		if isFine {
			io.WriteString(conn, SUCCESS + " :Seems like there is no error.\n")
		} else {
			io.WriteString(conn, FAIL + " :Seems like there is error in minting new tokens\n")
		}

	case "verify":
		if len(instArr) != 1 {
			io.WriteString(conn, FAIL + ": length of instruction != 1\n")
			break
		}

		marshalledTxs, err := fetchFa12Storage()
		if err != nil {
			io.WriteString(conn, FAIL + ": unable to fetch storage\n")
		}
		fmt.Println("\u2714 Relaying 1 Operation-Result packet [tezos] ---> [bifrostzone] @via {{bifrost}}...")
		io.WriteString(conn, marshalledTxs + "\n")

	default:
		fmt.Println("\u2714 Relaying 1 Operation-Result packet [tezos] ---> [bifrostzone] @via {{bifrost}}...")
		io.WriteString(conn, "Invalid Operation\n")
	}
}

func unlockTezos(address string, amount uint64) (bool, string, string) {
	var out1, err1 bytes.Buffer
	txcmd := commands.UnlockTokens(commands.ContractAddress, "alice", address, amount)
	cmd := exec.Command(txcmd[0], txcmd[1:]...)
	cmd.Stdout = &out1
	cmd.Stderr = &err1
	cmd.Run()

	// fmt.Println(out1.String())
	// fmt.Println(err1.String())


	if &err1 != nil && (strings.Contains(err1.String(), BURNFAILFLAG1) || strings.Contains(err1.String(), BURNFAILFLAG2)) {
	 	return false, out1.String(), err1.String()
	}

	if &out1 != nil && strings.Contains(out1.String(), BURNSUCCESSFLAG) {
		return true, out1.String(), err1.String()
	}

	return false, out1.String(), err1.String()
}

func mintFA12(address string, amount uint64) (bool, string, string) {
	var out1, err1 bytes.Buffer
	mintCmd := commands.MintFa12Cosmos(commands.FA12, "alice", address, amount)
	cmd := exec.Command(mintCmd[0], mintCmd[1:]...)
	cmd.Stdout = &out1
	cmd.Stderr = &err1
	cmd.Run()

	// fmt.Println(out1.String())
	// fmt.Println()
	// fmt.Println(err1.String())


	if &err1 != nil && (strings.Contains(err1.String(), MINTFAILFLAG1) || strings.Contains(err1.String(), MINTFAILFLAG2)) {
		return false, out1.String(), err1.String()
    }

    if &out1 != nil && strings.Contains(out1.String(), MINTSUCCESSFLAG) {
	   return true, out1.String(), err1.String()
    }

    return false, out1.String(), err1.String()


}

func fetchFa12Storage() (string, error) {
	fetchStorage := commands.FetchTezosStorage(commands.FA12)
	var outf, errf bytes.Buffer
	fcmd := exec.Command(fetchStorage[0], fetchStorage[1:]...)
	fcmd.Stdout = &outf
	fcmd.Stderr = &errf
	fcmd.Run()
	// fmt.Println("out...")
	// fmt.Println(outf.String())
	// fmt.Println("Error...")
	// fmt.Println(errf.String())
	// step2
	x := tzpengine.GetTxInSlice(outf.String())
	txs, err := tzpengine.GetFa12Txs(x)
	if err != nil {
		return "", err
	}


	// json marshall the tx
	bs, err := json.Marshal(txs)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}


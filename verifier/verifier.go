package verifier

import (
	"io"
	"net"
	"fmt"
	"os"
	"os/exec"
	"github.com/sap200/bifrost-relayer/tzpengine"
	"github.com/sap200/bifrost-relayer/commands"
	"encoding/json"
)

func LaunchVerificationChannel(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		// handle error
		fmt.Println("unable to launch verification channel...")
		os.Exit(1)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			fmt.Println("\u2718 Relaying 1 Verification-Failed packet [tezos] <--- [bifrostzone] @via {{bifrost}}...")
			fmt.Println(err)
		}
	go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("\u2714 \u27A2 Relaying 1 Verification packet [tezos] <--- [bifrostzone] @via {{bifrost}}..")
	// fetch the storage and relay it to bifrost-zone
	fetchStorage := commands.FetchTezosStorage(commands.ContractAddress)
	out, err := exec.Command(fetchStorage[0], fetchStorage[1:]...).Output()
	if err != nil {
		fmt.Println(err)
	}

	x := tzpengine.GetTxInSlice(string(out))
	txs := tzpengine.GetCosmosTxs(x)

	bs, err := json.Marshal(txs)
	if err != nil {
		fmt.Println(err)
	}

	io.WriteString(conn, string(bs) + "\n")
	fmt.Println("\u2714 \u27A2 Relaying 1 Verification packet [tezos ---> [bifrostzone] @via {{bifrost}}..")
}
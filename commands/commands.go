package commands

import (
	"github.com/sap200/bifrost-relayer/tzpengine"
	"fmt"
)

func GetContractStorageCommand(contractAddress string) []string {
	return []string{"tezos-client", "get", "contract", "storage", "for", contractAddress}
}

func RelayPacketCommand(tx tzpengine.TxsCosmos, sender string) []string {
	amtStr := fmt.Sprintf("%v", tx.Amount)
	txidStr := fmt.Sprintf("%v", tx.TxId)
	cmd := []string{PegZone,  "tx", "bifrost", "create-receivedTxs", 
			 txidStr, tx.TezosSender,tx.CosmosReceiver, amtStr, tx.Denom, tx.SrcChain, tx.DestChain, 
			 tx.TxStatus, tx.Timestamp, tx.RelayTimestamp, "--from", sender, "-y"}

	return cmd
}

func RelayFa12PacketCommand(tx tzpengine.Fa12CosmosTxs, sender string) []string {
	amtStr := fmt.Sprintf("%v", tx.Amount)
	txidStr := fmt.Sprintf("%v", tx.TxId)
	isApprovalStr := fmt.Sprintf("%v", tx.IsApproval)
	cmd := []string{"bifrostd", "tx", "bifrost", "create-receivedFa12Txs", txidStr, tx.TezosSender, 
				tx.Denom, tx.CosmosReceiver, amtStr, tx.SrcChain, tx.DestChain, tx.TxStatus, 
				isApprovalStr, tx.Approver, tx.Timestamp, tx.RelayTimestamp, "--from", sender, "-y"}

	return cmd
}

func QueryTxStatus(txid uint64) []string {
	txd := fmt.Sprintf("%v", txid)
	qcmd := []string{PegZone, "q", "bifrost", "show-receivedTxs", txd}
	return qcmd
}

func FetchTezosStorage(contract string) []string {
	cmd1 := []string{"tezos-client", "get", "contract", "storage", "for", contract}
	return cmd1
}

func UpdateStatus(contract string, key string, status string, txid uint64) []string {
	arg := fmt.Sprintf("(Right (Right (Left (Pair \"%s\" %v))))", status, txid)
	cmd1 := []string{"tezos-client", "transfer", "0", "from", key, "to", contract, "--fee", "1", "--arg", arg}
	return cmd1
}

func UpdateFa12Status(contract string, key string, status string, txid uint64) []string {
	arg := fmt.Sprintf("(Right (Right (Right (Right (Pair %v \"%s\")))))", txid, status)
	cmd1 := []string{"tezos-client", "transfer", "0", "from", key, "to", contract, "--fee", "1", "--arg", arg}
	return cmd1
}

func QueryCosmosKeyPair() []string {
	cmd := []string{PegZone, "keys", "show", KeyCosmos}	
	return cmd
}

func CreateKeyPair() []string {
	cmd := []string{PegZone, "keys", "add", KeyCosmos}
	return cmd
}

func UnlockTokens(contract string, key string, address string, amount uint64) []string {
	arg := fmt.Sprintf("(Right (Left (Pair \"%s\" %v)))", address, amount)
	cmd1 := []string{"tezos-client", "transfer", "0", "from", key, "to", contract, "--fee", "1", "--arg", arg, "--burn-cap", "0.00025"}
	return cmd1
}

func MintFa12Cosmos(contract string, key string, address string, amount uint64) []string {
	arg := fmt.Sprintf("(Right (Left (Left (Pair \"%s\" %v))))", address, amount)
	cmd1 := []string{"tezos-client", "-w", FinalityThreshold, "transfer", "0", "from", key, "to", contract, "--fee", "1", "--arg", arg}
	return cmd1
}
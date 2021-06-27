package tzpengine

import (
	"errors"
	"strconv"
	// "fmt"
)

func GetFa12Txs(slice []string) ([]Fa12CosmosTxs, error) {
	fa12TxSlice := []Fa12CosmosTxs{}
	for i, v := range slice {
		// fmt.Println(i, v)
		if v == "Elt" {
			amountStr := slice[i+5]
			amount, err := strconv.ParseUint(amountStr, 10, 64)
			if err != nil {
				return nil, err
			}
			isApprovalStr := slice[i+6]
			isApprovalStr = isApprovalStr[:len(isApprovalStr)-1]
			isApproval, err := strconv.ParseBool(isApprovalStr)
			if err != nil {
				return nil, err
			}
			approver := slice[i+7]
			cosmosReceiver := slice[i+8]
			destChain := slice[i+9]
			srcChain := slice[i+12]
			timeStamp := slice[i+13]
			txid := slice[i+15]
			txId, err := strconv.ParseUint(txid, 10, 64)
			if err != nil {
				return nil, err
			}
			status := slice[i+16]
			tezosSender := slice[i+17]

			fa12Tx := NewFa12CosmosTxs(txId, tezosSender, cosmosReceiver, amount, srcChain, destChain, status, isApproval, approver, timeStamp)
			fa12TxSlice = append(fa12TxSlice, fa12Tx)
		}
	}

	return fa12TxSlice, nil
}

func GetTotalSupply(slice []string) (uint64, error) {
	x, err := strconv.ParseUint(slice[len(slice)-1], 10, 64)
	if err != nil {
		return 0, errors.New("Unable to parse uint")
	}
	return x, nil
}

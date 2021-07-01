package tzpengine

import (
	"regexp"
	"strconv"
)

const (
	Amount_Offset         = 5
	CosmosReceiver_Offset = 6
	DestChain_Offset      = 8
	SrcChain_Offset       = 9
	Timestamp_Offset      = 12
	Txid_Offset           = 13
	TxStatus_Offset       = 14
	TezosSender_Offset    = 15
)

func GetTxInSlice(storage string) []string {
	re := regexp.MustCompile("[ ,\t,\n,\r,\"]+")
	newArr := re.Split(storage, -1)
	return newArr
}

func getTezosOwner(slice []string) string {
	return slice[len(slice)-2]
}

func GetCosmosTxs(slice []string) []TxsCosmos {
	txSlice := []TxsCosmos{}
	for i, v := range slice {
		if v == "Elt" {
			amount := slice[i+Amount_Offset]
			cosmosReceiver := slice[i+CosmosReceiver_Offset]
			destChain := slice[i+DestChain_Offset]
			srcChain := slice[i+SrcChain_Offset]
			timestamp := slice[i+Timestamp_Offset]
			txId := slice[i+Txid_Offset]
			txId = txId[:len(txId)-1]
			txStatus := slice[i+TxStatus_Offset]
			tezosSender := slice[i+TezosSender_Offset]
			newAmount, err := strconv.ParseUint(amount, 10, 64)
			if err != nil {
				txStatus = Failed
			}
			newTxid, err := strconv.ParseUint(txId, 10, 64)
			if err != nil {
				txStatus = Failed
			}

			// form a new cosmos tx
			newTx := NewTxCosmos(newTxid, tezosSender, cosmosReceiver, newAmount, srcChain, destChain, txStatus, timestamp, getTezosOwner(slice))

			txSlice = append(txSlice, newTx)

		}
	}

	return txSlice
}

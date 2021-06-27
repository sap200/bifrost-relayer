package tzpengine

import (
	"time"
)

type Fa12CosmosTxs struct {
	TxId           uint64 `json:"txid"`
	TezosSender    string `json:"tezos_sender"`
	Denom		   string `json:"denom"`
	CosmosReceiver string `json:"cosmos_receiver"`
	Amount         uint64 `json:"amount"`
	SrcChain       string `json:"src_chain"`
	DestChain      string `json:"dest_chain"`
	TxStatus       string `json:"tx_status"`
	IsApproval     bool   `json:"is_approval"`
	Approver       string `json:"approver"`
	Timestamp      string `json:"timestamp"`
	RelayTimestamp string `json:"relay_timestamp"`
}

func NewFa12CosmosTxs(
	txId uint64,
	tezosSender, cosmosReceiver string,
	amount uint64, srcChain,
	destChain,
	txStatus string, isApproval bool,
	approver, timestamp string,
) Fa12CosmosTxs {

	txs := Fa12CosmosTxs{
		TxId:           txId,
		TezosSender:    tezosSender,
		Denom:			"bifrost/token",
		CosmosReceiver: cosmosReceiver,
		Amount:         amount,
		SrcChain:       srcChain,
		DestChain:      destChain,
		TxStatus:       txStatus,
		IsApproval:     isApproval,
		Approver:       approver,
		Timestamp:      timestamp,
		RelayTimestamp: time.Now().String(),
	}

	return txs
}

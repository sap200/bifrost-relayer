package tzpengine

import (
	"encoding/json"
	"time"
)

// constants marking the transaction status
const (
	Initiated = "Initiated"
	Pending   = "Pending"
	Failed    = "Failed"
	Success   = "Success"
)

// Struct to denote tx sent to cosmos from tezos
type TxsCosmos struct {
	TxId           uint64 `json:"tx_id"`
	TezosSender    string `json:"tezos_sender"`
	CosmosReceiver string `json:"cosmos_receiver"`
	Amount         uint64 `json:"amount"`
	Denom          string `json:"denom"`
	SrcChain       string `json:"src_chain"`
	DestChain      string `json:"dest_chain"`
	TxStatus       string `json:"tx_status"`
	Timestamp      string `json:"timestamp"`
	RelayTimestamp string `json:"relay_timestamp"`
	TezosOwner     string `json:"tzeos_owner"`
}

// get a new txcosmos directly by calling this function
func NewTxCosmos(txId uint64, tzSender, cosmosReceiver string, amount uint64, srcChain, destChain, txStatus, timeStamp, tezosOwner string) TxsCosmos {
	return TxsCosmos{
		TxId:           txId,
		TezosSender:    tzSender,
		CosmosReceiver: cosmosReceiver,
		Amount:         amount,
		Denom:          "mutez",
		SrcChain:       srcChain,
		DestChain:      destChain,
		TxStatus:       txStatus,
		Timestamp:      timeStamp,
		RelayTimestamp: time.Now().String(),
		TezosOwner:     tezosOwner,
	}
}

func (t TxsCosmos) String() string {
	bs, err := json.Marshal(t)
	if err != nil {
		return err.Error()
	}

	return string(bs)
}

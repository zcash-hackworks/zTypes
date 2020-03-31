package zTypes

import (
	"encoding/json"
	"io/ioutil"
)

type BlockMetric struct {
	Height               int     `json:"height"`
	NumberofTransactions int     `json:"number_of_transactions"`
	SaplingValuePool     float64 `json:"sapling_value_pool"`
	SproutValuePool      float64 `json:"sprout_value_pool"`
	Size                 int     `json:"size"`
	Time                 int64   `json:"time"`
	NumberofTransparent  int     `json:"number_of_transparent_transactions"`
	NumberofShielded     int     `json:"number_of_shielded_transactions"`
	NumberofMixed        int     `json:"number_of_mixed_transactions"`
}

// GetBlockchainInfo return the zcashd rpc `getblockchaininfo` status
// https://zcash-rpc.github.io/getblockchaininfo.html
type GetBlockchainInfo struct {
	Chain                string     `json:"chain"`
	Blocks               int        `json:"blocks"`
	Headers              int        `json:"headers"`
	BestBlockhash        string     `json:"bestblockhash"`
	Difficulty           float64    `json:"difficulty"`
	VerificationProgress float64    `json:"verificationprogress"`
	SizeOnDisk           float64    `json:"size_on_disk"`
	SoftForks            []SoftFork `json:"softforks"`
}

type SoftFork struct {
	ID      string `json:"id"`
	Version int    `json:"version"`
}

type Block struct {
	Hash          string `json:"hash"`
	Confirmations int    `json:"confirmations"`
	Size          int    `json:"size"`
	Height        int    `json:"height"`
	Version       int    `json:"version"`
	//MerkleRoot        string        `json:"merkleroot"`
	//FinalSaplingRoot  string        `json:"finalsaplingroot"`
	TX   []Transaction `json:"tx"`
	Time int64         `json:"time"`
	//Nonce             string        `json:"nonce"`
	Difficulty        float64     `json:"difficulty"`
	PreviousBlockHash string      `json:"previousblockhash"`
	NextBlockHash     string      `json:"nextblockhash"`
	ValuePools        []ValuePool `json:"valuePools"`
}

func (b Block) TransactionTypes() (tTXs, sTXs int) {
	for _, tx := range b.TX {
		// If all 3 fields are empty, the transaction is transparent
		if len(tx.VJoinSplit) > 0 ||
			len(tx.VShieldedOutput) > 0 ||
			len(tx.VShieldedSpend) > 0 {
			tTXs++
		} else {
			// Otherwise, it's a shielded transaction
			sTXs++
		}
	}
	return tTXs, sTXs
}

func (b Block) WritetoFile(blockFile string) error {
	blockJSON, err := json.MarshalIndent(b, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(blockFile, blockJSON, 0644)
}

func (b Block) SaplingValuePool() float64 {
	for _, pool := range b.ValuePools {
		if pool.ID == "sapling" {
			return pool.ChainValue
		}
	}
	return 0
}

func (b Block) SproutValuePool() float64 {
	for _, pool := range b.ValuePools {
		if pool.ID == "sprout" {
			return pool.ChainValue
		}
	}
	return 0
}

func (b Block) NumberofTransactions() int {
	return len(b.TX)
}

type ValuePool struct {
	ID            string  `json:"id"`
	Monitored     bool    `json:"monitored"`
	ChainValue    float64 `json:"chainValue"`
	ChainValueZat float64 `json:"chainValueZat"`
	ValueDelta    float64 `json:"valueDelta"`
	ValueDeltaZat float64 `json:"valueDeltaZat"`
}

type Transaction struct {
	Hex             string                   `json:"hex"`
	Txid            string                   `json:"txid"`
	Version         int                      `json:"version"`
	Locktime        int                      `json:"locktime"`
	ExpiryHeight    int                      `json:"expirtheight"`
	VIn             []VInTX                  `json:"vin"`
	VOut            []VOutTX                 `json:"vout"`
	VJoinSplit      []VJoinSplitTX           `json:"vjoinsplit"`
	ValueBalance    float64                  `json:"valueBalance"`
	VShieldedSpend  []map[string]interface{} `json:"vShieldedSpend"`
	VShieldedOutput []map[string]interface{} `json:"vShieldedOutput"`
}

// TransparentInAndOut return if there are transparent
// inputs and outputs
func (t Transaction) TransparentInAndOut() bool {
	return len(t.VIn) > 0 && len(t.VOut) > 0
}

// IsTransparent returns if the transaction contains
// no shielded addresses
func (t Transaction) IsTransparent() bool {
	return t.TransparentInAndOut() &&
		len(t.VJoinSplit) == 0 &&
		t.ValueBalance == 0 &&
		len(t.VShieldedSpend) == 0 &&
		len(t.VShieldedSpend) == 0
}

// ContainsSprout returns if a transaction contains
// sprout transaction data
func (t Transaction) ContainsSprout() bool {
	return len(t.VJoinSplit) > 0
}

// ContainsSapling returns if a transaction contains
// sapling transaction data
// Check that there is a valueBalance value (positive or negative)
// Check that there is data for either VShieldedSpend or VShieldedOutput
func (t Transaction) ContainsSapling() bool {
	return t.ValueBalance != 0 && (len(t.VShieldedSpend) > 0 ||
		len(t.VShieldedOutput) > 0)
}

// IsShielded returns if the transaction contains
// no transparent addresses
func (t Transaction) IsShielded() bool {
	return !t.TransparentInAndOut() &&
		(t.ContainsSprout() || t.ContainsSapling())
}

// IsMixed returns if the transaction contains
// transparent addresses and shielded transaction data
func (t Transaction) IsMixed() bool {
	tInOrOut := len(t.VIn) > 0 || len(t.VOut) > 0
	return tInOrOut &&
		(t.ContainsSprout() || t.ContainsSapling())
}

type VInTX struct {
	TxID      string `json:"txid"`
	VOut      int    `json:"vout"`
	ScriptSig ScriptSig
	Sequence  int `json:"sequemce"`
}
type ScriptSig struct {
	//Asm string `json:"asm"`
	//Hex string `json:"hex"`
}
type VOutTX struct {
	Value        float64
	N            int
	ScriptPubKey ScriptPubKey
}
type ScriptPubKey struct {
	//Asm       string   `json:"asm"`
	//Hex       string   `json:"hex"`
	//ReqSigs   int      `json:"reqSigs`
	Type      string   `json:"type"`
	Addresses []string `json:"addresses"`
}
type VJoinSplitTX struct {
	VPubOldld float64 `json:"vpub_old"`
	VPubNew   float64 `json:"vpub_new"`
}

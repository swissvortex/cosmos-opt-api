package models

import "encoding/json"

func UnmarshalBlockchain(response []byte) (Blockchain, error) {
	var blockchain Blockchain
	err := json.Unmarshal(response, &blockchain)
	return blockchain, err
}

func (r *Blockchain) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Blockchain struct {
	Jsonrpc          string           `json:"jsonrpc"`
	ID               int64            `json:"id"`
	BlockchainResult BlockchainResult `json:"result"`
}

type BlockchainResult struct {
	LastHeight string      `json:"last_height"`
	BlockMetas []BlockMeta `json:"block_metas"`
}

type BlockMeta struct {
	BlockID   BlockID `json:"block_id"`
	BlockSize string  `json:"block_size"`
	Header    Header  `json:"header"`
	NumTxs    string  `json:"num_txs"`
}

type ChainID string

const (
	Cosmoshub4 ChainID = "cosmoshub-4"
)

type ConsensusHash string

const (
	The0F2908883A105C793B74495Eb7D6Df2Eea479Ed7Fc9349206A65Cb0F9987A0B8 ConsensusHash = "0F2908883A105C793B74495EB7D6DF2EEA479ED7FC9349206A65CB0F9987A0B8"
)

type EvidenceHash string

const (
	E3B0C44298Fc1C149Afbf4C8996Fb92427Ae41E4649B934Ca495991B7852B855 EvidenceHash = "E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855"
)

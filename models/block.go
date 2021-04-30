package models

import "encoding/json"

func UnmarshalBlock(response []byte) (BlockResponse, error) {
	var block BlockResponse
	err := json.Unmarshal(response, &block)
	return block, err
}

func (r *BlockResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type BlockResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int64  `json:"id"`
	Result  Result `json:"result"`
}

type Result struct {
	BlockID BlockID `json:"block_id"`
	Block   Block   `json:"block"`
}

type Block struct {
	Header     Header     `json:"header"`
	Data       Data       `json:"data"`
	Evidence   Evidence   `json:"evidence"`
	LastCommit LastCommit `json:"last_commit"`
}

type Data struct {
	Txs []string `json:"txs"`
}

type Evidence struct {
	Evidence []interface{} `json:"evidence"`
}

type Header struct {
	Version            Version `json:"version"`
	ChainID            string  `json:"chain_id"`
	Height             string  `json:"height"`
	Time               string  `json:"time"`
	LastBlockID        BlockID `json:"last_block_id"`
	LastCommitHash     string  `json:"last_commit_hash"`
	DataHash           string  `json:"data_hash"`
	ValidatorsHash     string  `json:"validators_hash"`
	NextValidatorsHash string  `json:"next_validators_hash"`
	ConsensusHash      string  `json:"consensus_hash"`
	AppHash            string  `json:"app_hash"`
	LastResultsHash    string  `json:"last_results_hash"`
	EvidenceHash       string  `json:"evidence_hash"`
	ProposerAddress    string  `json:"proposer_address"`
}

type BlockID struct {
	Hash  string `json:"hash"`
	Parts Parts  `json:"parts"`
}

type Parts struct {
	Total int64  `json:"total"`
	Hash  string `json:"hash"`
}

type Version struct {
	Block string `json:"block"`
}

type LastCommit struct {
	Height     string      `json:"height"`
	Round      int64       `json:"round"`
	BlockID    BlockID     `json:"block_id"`
	Signatures []Signature `json:"signatures"`
}

type Signature struct {
	BlockIDFlag      int64   `json:"block_id_flag"`
	ValidatorAddress string  `json:"validator_address"`
	Timestamp        string  `json:"timestamp"`
	Signature        *string `json:"signature"`
}

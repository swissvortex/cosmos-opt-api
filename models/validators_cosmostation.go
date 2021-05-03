package models

import "encoding/json"

type Validators []Validator

func UnmarshalValidators(data []byte) (Validators, error) {
	var r Validators
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Validators) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Validator struct {
	Rank              int64  `json:"rank"`
	AccountAddress    string `json:"account_address"`
	OperatorAddress   string `json:"operator_address"`
	ConsensusPubkey   string `json:"consensus_pubkey"`
	Jailed            bool   `json:"jailed"`
	Status            int64  `json:"status"`
	Tokens            string `json:"tokens"`
	DelegatorShares   string `json:"delegator_shares"`
	Moniker           string `json:"moniker"`
	Identity          string `json:"identity"`
	Website           string `json:"website"`
	Details           string `json:"details"`
	UnbondingHeight   string `json:"unbonding_height"`
	UnbondingTime     string `json:"unbonding_time"`
	Rate              string `json:"rate"`
	MaxRate           string `json:"max_rate"`
	MaxChangeRate     string `json:"max_change_rate"`
	UpdateTime        string `json:"update_time"`
	Uptime            Uptime `json:"uptime"`
	MinSelfDelegation string `json:"min_self_delegation"`
	KeybaseURL        string `json:"keybase_url"`
}

type Uptime struct {
	Address      string `json:"address"`
	MissedBlocks int64  `json:"missed_blocks"`
	OverBlocks   int64  `json:"over_blocks"`
}

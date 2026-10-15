package main

type Genesis struct {
	GenesisTime string           `json:"genesis_time"`
	ChainID     string           `json:"chain_id"`
	Balances    map[string]int64 `json:"balances"`
}
type Transaction struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Value int64  `json:"value"`
}

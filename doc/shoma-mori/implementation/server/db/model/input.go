package model

// Input represents transaction inputs
type Input struct {
	UTXO       Output   `json:"utxo"`
	Signature1 string   `json:"signature1"`
	Signature2 string   `json:"signature2"`
	Siblings   []Output `json:"siblings"`
}

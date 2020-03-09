package model

// Input represents transaction inputs
type Input struct {
	UTXO       Output `json:"utxo"`
	Signature1 string `json:"sig1"`
	Signature2 string `json:"sig2"`
}

package model

// Transaction is
type Transaction struct {
	Inputs  []Input  `json:"inputs"`
	Outputs []Output `json:"outputs"`
}

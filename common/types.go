package common

type Address []byte

// Defines an Input / Output tuple
type Tuple struct {
	Address		Address		`json: address`
	Value		int			`json: value`
}

// defines the payload of the message which we want to have signed
type TransactionSignReq struct {
	Inputs 	[]Tuple
	Outputs []Tuple
}

type TransactionSignRes struct {
	Outputs	[]Tuple
}
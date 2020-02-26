package common

type Address []byte

// Defines an Input / Output tuple
type Tuple struct {
	Address		Address
	Value		int
	Id			int
}

// defines the payload of the message which we want to have signed
type Transaction struct {
	Inputs 	[]Tuple
	Outputs []Tuple
}

type TransactionSignRes struct {
	Outputs	[]Tuple
}
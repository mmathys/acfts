package common

const (
	AddressLength = 1
)

type Address [AddressLength]byte

// Defines an Input / Output tuple; with extra fields
type Tuple struct {
	Address		Address
	Value		int
	Id			int
}

// Defines the payload of the message which we want to have signed
type Transaction struct {
	Inputs 		[]Tuple
	Outputs 	[]Tuple
}

type TransactionSignRes struct {
	Outputs		[]Tuple
}

type Wallet struct {
	Address 	Address
	UTXO 		map[int]Tuple
}
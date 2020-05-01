package common

import "crypto/rand"

func RandomIdentifier() Identifier {
	slice := make([]byte, IdentifierLength)
	rand.Read(slice)
	array := [IdentifierLength]byte{}
	copy(array[:], slice[:IdentifierLength])
	return array
}

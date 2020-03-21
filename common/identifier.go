package common

import "crypto/rand"

func RandomIdentifier() Identifier {
	slice := make([]byte, IdentifierLength)
	rand.Read(slice)

	result := [IdentifierLength]byte{}
	copy(result[:], slice[:IdentifierLength])

	return result
}
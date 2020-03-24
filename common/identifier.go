package common

import "crypto/rand"

func RandomIdentifier() Identifier {
	slice := make([]byte, IdentifierLength)
	rand.Read(slice)

	return slice
}
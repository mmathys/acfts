package common

import (
	"crypto/rand"
	"github.com/oasislabs/ed25519"
)

func GenerateKey() *Identity {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	return &Identity{
		Address:    pub,
		PrivateKey: priv,
	}
}

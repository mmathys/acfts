package common

import (
	"bytes"
	"math"
	"testing"
)

func TestServerQuorum(t *testing.T) {
	ServerAddresses = []Address{}
	for i := 0; i < 10; i++ {
		addr := GenerateKey(ModeEdDSA, i).GetAddress()
		ServerAddresses = append(ServerAddresses, addr)
	}

	quorumSize := int(math.Ceil(2.0 / 3.0 * float64(10)))

	q1 := ServerQuorum()
	if len(q1) != quorumSize {
		t.Fatal("wrong size")
	}

	q2 := ServerQuorum()

	if bytes.Equal(q1[0], q2[0]) {
		t.Fatal("quorum must be random")
	}
}

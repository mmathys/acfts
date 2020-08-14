package common

import (
	"bytes"
	"fmt"
	"math"
	"testing"
)

func TestServerQuorum(t *testing.T) {
	numV := 4

	ServerAddresses = []Address{}
	for i := 0; i < numV; i++ {
		addr := GenerateKey(ModeNaive, i).GetAddress()
		ServerAddresses = append(ServerAddresses, addr)
	}

	quorumSize := int(math.Ceil(2.0 / 3.0 * float64(numV)))
	fmt.Println("quorum size =", quorumSize)

	q1 := ServerQuorum()
	if len(q1) != quorumSize {
		t.Fatal("wrong size")
	}

	q2 := ServerQuorum()

	fmt.Println(q1)
	fmt.Println(q2)

	if bytes.Equal(q1[0], q2[0]) {
		t.Fatal("quorum must be random")
	}
}

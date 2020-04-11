package common

import (
	"bytes"
	"testing"
)

func TestIdentifier(t *testing.T) {
	i1 := RandomIdentifier()
	i2 := RandomIdentifier()

	if bytes.Equal(i1, i2) {
		t.Error("two identifiers are the same")
	}
}

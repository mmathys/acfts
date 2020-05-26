package ed25519

import (
	"crypto"
	"crypto/rand"
	"github.com/oasislabs/ed25519"
	"io"
	"testing"
)

const (
	batchCount    = 64
)

func testBatchInit(tb testing.TB, r io.Reader, batchSize int, opts *ed25519.Options) ([]ed25519.PublicKey, [][]byte, [][]byte) {
	sks := make([]ed25519.PrivateKey, batchSize)
	pks := make([]ed25519.PublicKey, batchSize)
	sigs := make([][]byte, batchSize)
	messages := make([][]byte, batchSize)

	// generate keys
	for i := 0; i < batchSize; i++ {
		pub, priv, err := ed25519.GenerateKey(r)
		if err != nil {
			tb.Fatalf("failed to generate key #%d: %v", i, err)
		}

		sks[i], pks[i] = priv, pub
	}

	// generate messages
	for i := 0; i < batchSize; i++ {
		// Yes, this generates too much, but the amount read from r needs
		// to match what was used to generate the good final y coord.
		m := make([]byte, 128)
		if _, err := io.ReadFull(r, m); err != nil {
			tb.Fatalf("failed to generate message #%d: %v", i, err)
		}
		mLen := (i & 127) + 1
		messages[i] = m[:mLen]

		// Pre-hash the message if required.
		if opts.Hash != crypto.Hash(0) {
			h := opts.Hash.New()
			_, _ = h.Write(messages[i])
			messages[i] = h.Sum(nil)
		}
	}

	// sign messages
	for i := 0; i < batchSize; i++ {
		sig, err := sks[i].Sign(nil, messages[i], opts)
		if err != nil {
			tb.Fatalf("failed to generate signature #%d: %v", i, err)
		}
		sigs[i] = sig
	}

	return pks, sigs, messages
}

func BenchmarkVerifyBatch64(b *testing.B) {
	var opts ed25519.Options
	pks, sigs, messages := testBatchInit(b, rand.Reader, batchCount, &opts)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ok, _, _ := ed25519.VerifyBatch(nil, pks[:], messages[:], sigs[:], &opts)
		if !ok {
			b.Fatalf("unexpected batch verification failure!")
		}
	}
}

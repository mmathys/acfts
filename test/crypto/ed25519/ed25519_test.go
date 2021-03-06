package ed25519

import (
	"crypto"
	"crypto/rand"
	"github.com/mmathys/acfts/common"
	"github.com/oasislabs/ed25519"
	"io"
	"testing"
)

const (
	batchCount = 64
)

type zeroReader struct{}

func (zeroReader) Read(buf []byte) (int, error) {
	for i := range buf {
		buf[i] = 0
	}
	return len(buf), nil
}

func BenchmarkSignEd25519(b *testing.B) {
	var zero zeroReader
	_, priv, err := ed25519.GenerateKey(zero)
	if err != nil {
		b.Fatal(err)
	}
	message := []byte("Hello, world!")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ed25519.Sign(priv, message)
	}
}

func BenchmarkVerifyEd25519(b *testing.B) {
	var zero zeroReader
	pub, priv, err := ed25519.GenerateKey(zero)
	if err != nil {
		b.Fatal(err)
	}
	message := []byte("Hello, world!")
	signature := ed25519.Sign(priv, message)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ed25519.Verify(pub, message, signature)
	}
}

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

func TestVerifyBatch64(t *testing.T) {
	hash := make([]byte, 64) // random hash
	rand.Read(hash)
	numSigs := 64
	var sigs []common.Signature
	for i := 0; i < numSigs; i++ {
		key := common.GenerateKey(common.ModeEdDSA, 0)
		sig := key.SignHash(hash)
		sigs = append(sigs, *sig)
	}

	ok, err := common.VerifyEdDSABatch(sigs, hash)
	if err != nil {
		panic(err)
	}
	if !ok {
		panic("batch verification failed")
	}

	sigs[0].Signature[0]++
	ok2, err2 := common.VerifyEdDSABatch(sigs, hash)
	if err2 != nil {
		panic(err2)
	}
	if ok2 {
		panic("batch verification should have failed")
	}

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

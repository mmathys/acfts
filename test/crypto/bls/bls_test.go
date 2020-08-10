package bls

import (
	"fmt"
	"github.com/herumi/bls-eth-go-binary/bls"
	"testing"
)

func TestPrintGeneratedTOfNBLSKeys(t *testing.T) {
	bls.Init(bls.BLS12_381)
	bls.SetETHmode(bls.EthModeDraft07)

	// 3-of-4 scheme; 3f+1=n, f=1
	n := 4
	k := 3

	// k secret keys erstellen
	secs := make([]bls.SecretKey, k)
	for i := 0; i < k; i++ {
		secs[i].SetByCSPRNG()
	}

	// n shares aus k keys erstellen
	ids := make([]bls.ID, n)
	shares := make([]bls.SecretKey, n)
	for i := 0; i < n; i++ {
		ids[i].SetLittleEndian([]byte{uint8(i + 1)})
		shares[i].Set(secs, &ids[i])
	}

	// master public key
	mpk := secs[0].GetPublicKey()

	// generate public keys for each user
	pubs := make([]*bls.PublicKey, n)
	for i := 0; i < n; i++ {
		pubs[i] = shares[i].GetPublicKey()
	}

	fmt.Println("master key:")
	fmt.Printf("{\"%x\",\"%x\"},\n", mpk.Serialize(), secs[0].Serialize())

	fmt.Printf("%v shares:\n", n)
	for i := 0; i < n; i++ {
		fmt.Printf("{\"%x\",\"%x\"},\n", pubs[i].Serialize(), shares[i].Serialize())
	}
}

// Below: benchmarks taken from BLS library
func getSecPubHash() (*bls.SecretKey, *bls.PublicKey, []byte) {
	var sec bls.SecretKey
	sec.SetByCSPRNG()
	pub := sec.GetPublicKey()
	var x bls.Fp2
	x.D[0].SetByCSPRNG()
	x.D[1].SetByCSPRNG()
	hash := x.Serialize()
	return &sec, pub, hash
}

func BenchmarkSignHash(b *testing.B) {
	b.StopTimer()
	err := bls.Init(bls.BLS12_381)
	if err != nil {
		b.Fatal(err)
	}
	sec, _, hash := getSecPubHash()
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		sec.SignHash(hash)
	}
	b.StopTimer()
}

func BenchmarkVerifyHash(b *testing.B) {
	b.StopTimer()
	err := bls.Init(bls.BLS12_381)
	if err != nil {
		b.Fatal(err)
	}
	sec, pub, hash := getSecPubHash()
	sig := sec.SignHash(hash)
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		sig.VerifyHash(pub, hash)
	}
	b.StopTimer()
}

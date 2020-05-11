package funset

const arrayLength = 10000000

type FunSet struct {
	data [arrayLength]*LinkedList
}

func NewFunSet() *FunSet {
	table := [arrayLength]*LinkedList{}

	for i := 0; i < arrayLength; i++ {
		table[i] = NewList()
	}
	return &FunSet{
		table,
	}
}

// Use some hash. Copy from normal go map?
// I don't know what an "identifier" is, but ultimately,
// we're storing hashes of stuff that we signed, right?
// Can we get away without hashing here again at all?
func hash(t [32]byte) int {
	h := 0
	for i := 0; i < 32; i++ {
		h *= 256
		h += int(t[i])
		h %= arrayLength
	}
	return h
}

func index(hash int) int {
	return hash % arrayLength
}

func (h *FunSet) Insert(k [32]byte) bool {
	index := index(hash(k))
	return h.data[index].Add(&k)
}

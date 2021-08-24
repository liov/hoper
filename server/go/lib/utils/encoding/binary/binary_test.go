package binary

import (
	"encoding/binary"
	"testing"
)

func TestIntFromBinary(t *testing.T) {
	b := []byte{0, 1, 0, 0, 0, 0, 0, 0}
	t.Log(IntFromBinary(b))
	t.Log(binary.BigEndian.Uint64(b))
	b = IntToBinary(15)
	t.Log(IntFromBinary(b))
}

func BenchmarkIntFromBinary(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bb := UIntToBinary(15)
		UIntFromBinary(bb)
	}
}

func BenchmarkBigEndianUint64(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bb := make([]byte, 8, 8)
		binary.BigEndian.PutUint64(bb, 15)
		binary.BigEndian.Uint64(bb)
	}
}

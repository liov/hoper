package gocall

import (
	"testing"
)

func BenchmarkEmptyCgoCalls(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Cempty()
	}
}

func BenchmarkEmptyGoCalls(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Empty()
	}
}

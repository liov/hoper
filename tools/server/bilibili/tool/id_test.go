package tool

import "testing"

func TestBv2av(t *testing.T) {
	t.Log(Bv2av("BV15G411h7gY"))
}

func TestAv2bv(t *testing.T) {
	t.Log(Av2bv(15576810))
}

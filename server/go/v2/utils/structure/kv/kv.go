package kv


type KVs struct {
	noCopy noCopy //nolint:unused,structcheck

	args []KV
	buf  []byte
}

type KV struct {
	key     []byte
	value   []byte
	noValue bool
}

// noCopy may be embedded into structs which must not be copied
// after the first use.
//
// See https://golang.org/issues/8005#issuecomment-190753527
// for details.
type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

package gen

import (
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"sync"
	"sync/atomic"
)

// 随机id
func GenID() uint64 {
	return IDGenerator.NewID()
}

var currentID uint64 = 0

// 单机程顺序id
func GenOrderID() uint64 {
	return atomic.AddUint64(&currentID, 1)
}

var IDGenerator = &defaultIDGenerator{}

func init() {
	IDGenerator.init()
}

type defaultIDGenerator struct {
	sync.Mutex

	IDAdd  [2]uint64
	IDRand *rand.Rand

	initOnce sync.Once
}

// init initializes the generator on the first call to avoid consuming entropy
// unnecessarily.
func (gen *defaultIDGenerator) init() {
	gen.initOnce.Do(func() {
		// initialize traceID and spanID generators.
		var rngSeed int64
		for _, p := range []interface{}{&rngSeed, &gen.IDAdd} {
			binary.Read(crand.Reader, binary.LittleEndian, p)
		}
		gen.IDRand = rand.New(rand.NewSource(rngSeed))

	})
}

func (gen *defaultIDGenerator) NewID() uint64 {
	gen.Lock()
	id := gen.IDRand.Uint64()
	gen.Unlock()
	return id
}

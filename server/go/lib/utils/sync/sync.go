package synci

import (
	"sync"
	"unsafe"
)

type WaitGroup struct {
	noCopy noCopy
	state1 uint64
	state2 uint32
}

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

func WaitGroupState(wg *sync.WaitGroup) (counter int32, wcounter uint32) {
	wgc := (*WaitGroup)(unsafe.Pointer(wg))
	statep, _ := wgc.state()
	return int32(*statep >> 32), uint32(*statep)
}

func (wg *WaitGroup) state() (statep *uint64, semap *uint32) {
	if unsafe.Alignof(wg.state1) == 8 || uintptr(unsafe.Pointer(&wg.state1))%8 == 0 {
		// state1 is 64-bit aligned: nothing to do.
		return &wg.state1, &wg.state2
	} else {
		// state1 is 32-bit aligned but not 64-bit aligned: this means that
		// (&state1)+4 is 64-bit aligned.
		state := (*[3]uint32)(unsafe.Pointer(&wg.state1))
		return (*uint64)(unsafe.Pointer(&state[1])), &state[0]
	}
}

func WaitGroupStopWait(wg *sync.WaitGroup) {
	state, _ := WaitGroupState(wg)
	wg.Add(int(-state))
}

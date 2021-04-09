package tdragonboat

import (
	"encoding/binary"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/lni/dragonboat/v3/statemachine"
	"io"
	"io/ioutil"
)

type ExampleStateMachine struct {
	ClusterID uint64
	NodeID    uint64
	Count     uint64
}

// NewExampleStateMachine creates and return a new ExampleStateMachine object.
func NewExampleStateMachine(clusterID uint64,
	nodeID uint64) statemachine.IStateMachine {
	return &ExampleStateMachine{
		ClusterID: clusterID,
		NodeID:    nodeID,
		Count:     0,
	}
}

// Lookup performs local lookup on the ExampleStateMachine instance. In this example,
// we always return the Count value as a little endian binary encoded byte
// slice.
func (s *ExampleStateMachine) Lookup(query []byte) []byte {
	result := make([]byte, 8)
	binary.LittleEndian.PutUint64(result, s.Count)
	return result
}

// Update updates the object using the specified committed raft entry.
func (s *ExampleStateMachine) Update(data []byte) uint64 {
	// in this example, we print out the following message for each
	// incoming update request. we also increase the counter by one to remember
	// how many updates we have applied
	s.Count++
	log.Infof("from ExampleStateMachine.Update(), msg: %s, count:%d\n",
		string(data), s.Count)
	return uint64(len(data))
}

// SaveSnapshot saves the current IStateMachine state into a snapshot using the
// specified io.Writer object.
func (s *ExampleStateMachine) SaveSnapshot(w io.Writer,
	fc statemachine.ISnapshotFileCollection,
	done <-chan struct{}) (uint64, error) {
	// as shown above, the only state that can be saved is the Count variable
	// there is no external file in this IStateMachine example, we thus leave
	// the fc untouched
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, s.Count)
	_, err := w.Write(data)
	if err != nil {
		return 0, err
	}
	return uint64(len(data)), nil
}

// RecoverFromSnapshot recovers the state using the provided snapshot.
func (s *ExampleStateMachine) RecoverFromSnapshot(r io.Reader,
	files []statemachine.SnapshotFile,
	done <-chan struct{}) error {
	// restore the Count variable, that is the only state we maintain in this
	// example, the input files is expected to be empty
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	v := binary.LittleEndian.Uint64(data)
	s.Count = v
	return nil
}

// CloseDao closes the IStateMachine instance. There is nothing for us to cleanup
// or release as this is a pure in memory data store. Note that the CloseDao
// method is not guaranteed to be called as node can crash at any time.
func (s *ExampleStateMachine) Close() {}

// GetHash returns a uint64 representing the current object state.
func (s *ExampleStateMachine) GetHash() uint64 {
	// the only state we have is that Count variable. that uint64 value pretty much
	// represents the state of this IStateMachine
	return s.Count
}

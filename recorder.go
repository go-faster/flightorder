package flightorder

import (
	"sync"
)

// Recorder records internal route events.
// Used for testing purposes.
type recorder struct {
	takeCalls []*Ticket
	compCalls []*Ticket
	completed []*Ticket
	mux       sync.Mutex
}

func newRecorder() *recorder {
	return &recorder{}
}

func (r *recorder) takeCall(t *Ticket) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.takeCalls = append(r.takeCalls, t)
}

func (r *recorder) completeCall(t *Ticket) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.compCalls = append(r.compCalls, t)
}

func (r *recorder) recordCompleted(t *Ticket) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.completed = append(r.completed, t)
}

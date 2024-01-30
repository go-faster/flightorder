package flightorder

import "sync"

// TicketAllocator is responsible for ticket allocation.
type TicketAllocator interface {
	AcquireTicket() *Ticket
	ReleaseTicket(t *Ticket)
}

var (
	_ TicketAllocator = (StdAllocator{})
	_ TicketAllocator = (*SyncpoolAllocator)(nil)
	_ TicketAllocator = (*testAllocator)(nil)
)

// StdAllocator is a standard ticket allocator without any memory reuse.
type StdAllocator struct{}

// AcquireTicket acquires a new ticket.
func (StdAllocator) AcquireTicket() *Ticket {
	return newTicket()
}

// ReleaseTicket does nothing. Let GC erase ticket for us.
func (StdAllocator) ReleaseTicket(t *Ticket) { t.reset() }

// SyncpoolAllocator uses sync.Pool under the hood to reuse allocated tickets.
type SyncpoolAllocator struct {
	pool *sync.Pool
}

// NewSyncpoolAllocator creates new SyncpoolAllocator.
func NewSyncpoolAllocator() *SyncpoolAllocator {
	return &SyncpoolAllocator{
		pool: &sync.Pool{
			New: func() any {
				return newTicket()
			},
		},
	}
}

// AcquireTicket acquires a new ticket from the pool.
func (a *SyncpoolAllocator) AcquireTicket() *Ticket {
	return a.pool.Get().(*Ticket)
}

// ReleaseTicket releases ticket to the pool.
func (a *SyncpoolAllocator) ReleaseTicket(t *Ticket) {
	t.reset()
	a.pool.Put(t)
}

// testAllocator is a test allocator.
type testAllocator struct {
	released []*Ticket
	mux      sync.Mutex
}

func newTestAllocator() *testAllocator {
	return &testAllocator{}
}

func (a *testAllocator) AcquireTicket() *Ticket {
	return newTicket()
}

func (a *testAllocator) ReleaseTicket(t *Ticket) {
	a.mux.Lock()
	defer a.mux.Unlock()
	a.released = append(a.released, t)
}

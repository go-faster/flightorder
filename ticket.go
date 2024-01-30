package flightorder

// Ticket is a route ticket.
type Ticket struct {
	// previous ticket of a route, may be nil
	prev *Ticket
	// holds ticket completion status
	completed chan struct{}
}

func newTicket() *Ticket {
	return &Ticket{
		completed: make(chan struct{}, 1),
	}
}

func (f *Ticket) reset() {
	f.prev = nil
	if len(f.completed) != 0 {
		// Means that the ticket processing code is broken.
		panic("unreachable")
	}
}

package flightorder

import (
	"context"
	"fmt"
	"sync"
)

// Route is responsible for tickets processing.
type Route struct {
	allocator TicketAllocator
	last      *Ticket
	recorder  *recorder
	mux       sync.Mutex
}

// RouteParams sets route parameters.
type RouteParams struct {
	// TicketAllocator sets custom ticket allocator.
	// Optional.
	TicketAllocator TicketAllocator
}

// NewRoute creates new route for tickets processing.
func NewRoute(params RouteParams) *Route {
	if params.TicketAllocator == nil {
		params.TicketAllocator = StdAllocator{}
	}

	return &Route{
		allocator: params.TicketAllocator,
	}
}

// Ticket takes a new ticket.
func (r *Route) Ticket() *Ticket {
	r.mux.Lock()
	defer r.mux.Unlock()

	ticket := r.allocator.AcquireTicket()
	ticket.prev = r.last
	r.last = ticket

	if r.recorder != nil {
		r.recorder.takeCall(ticket)
	}

	return ticket
}

// CompleteTicket completes a ticket.
// Waits for previous taken tickets to complete first, if any.
// Completion function is optional.
func (r *Route) CompleteTicket(ctx context.Context, t *Ticket, completion func(ctx context.Context) error) error {
	if completion == nil {
		completion = func(ctx context.Context) error { return nil }
	}

	if r.recorder != nil {
		r.recorder.completeCall(t)
	}

	if t.prev == nil {
		return r.completeTail(ctx, t, completion)
	}

	if err := r.waitFor(ctx, t.prev); err != nil {
		return fmt.Errorf("wait for previous ticket: %w", err)
	}

	r.allocator.ReleaseTicket(t.prev)
	t.prev = nil

	return r.completeTail(ctx, t, completion)
}

func (r *Route) completeTail(ctx context.Context, t *Ticket, f func(ctx context.Context) error) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	// Last ticket on a trip. No tickets ahead.
	if r.last == t {
		r.last = nil
		r.allocator.ReleaseTicket(t)
		if r.recorder != nil {
			r.recorder.recordCompleted(t)
		}

		return f(ctx)
	}

	if err := f(ctx); err != nil {
		return err
	}

	// There is a ticket ahead.
	// Mark current ticket as completed.
	t.completed <- struct{}{}
	if r.recorder != nil {
		r.recorder.recordCompleted(t)
	}

	return nil
}

func (r *Route) waitFor(ctx context.Context, f *Ticket) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-f.completed:
		return nil
	}
}

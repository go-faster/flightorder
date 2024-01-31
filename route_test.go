package flightorder

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRoute(t *testing.T) {
	t.Run("1 ticket", func(t *testing.T) {
		alloc := newTestAllocator()
		route := NewRoute(RouteParams{
			TicketAllocator: alloc,
		})

		t1 := route.Ticket()
		require.NoError(t, route.CompleteTicket(context.TODO(), CompleteTicketParams{Ticket: t1}))
		require.Nil(t, route.last)
		require.Equal(t, []*Ticket{t1}, alloc.released)
	})

	t.Run("2 tickets, completion t1 t2", func(t *testing.T) {
		alloc := newTestAllocator()
		route := NewRoute(RouteParams{
			TicketAllocator: alloc,
		})

		t1 := route.Ticket()
		t2 := route.Ticket()
		require.NoError(t, route.CompleteTicket(context.TODO(), CompleteTicketParams{Ticket: t1}))
		require.NoError(t, route.CompleteTicket(context.TODO(), CompleteTicketParams{Ticket: t2}))
		require.Nil(t, route.last)
		require.Equal(t, []*Ticket{t1, t2}, alloc.released)
	})

	t.Run("3 tickets, completion t3 t2 t1", func(t *testing.T) {
		alloc := newTestAllocator()
		route := NewRoute(RouteParams{
			TicketAllocator: alloc,
		})

		rec := newRecorder()
		route.recorder = rec

		t1 := route.Ticket()
		t2 := route.Ticket()
		t3 := route.Ticket()

		var wg sync.WaitGroup
		wg.Add(3)
		go func() {
			time.Sleep(time.Millisecond * 10)
			require.NoError(t, route.CompleteTicket(context.TODO(), CompleteTicketParams{Ticket: t3}))
			wg.Done()
		}()

		go func() {
			time.Sleep(time.Millisecond * 20)
			require.NoError(t, route.CompleteTicket(context.TODO(), CompleteTicketParams{Ticket: t2}))
			wg.Done()
		}()

		go func() {
			time.Sleep(time.Millisecond * 30)
			require.NoError(t, route.CompleteTicket(context.TODO(), CompleteTicketParams{Ticket: t1}))
			wg.Done()
		}()

		wg.Wait()
		require.Equal(t, rec.compCalls, []*Ticket{t3, t2, t1})
		require.Equal(t, rec.takeCalls, rec.completed)
		require.Nil(t, route.last)
		require.Equal(t, rec.completed, alloc.released)
	})

	t.Run("3 tickets, completion t2 t3 t1", func(t *testing.T) {
		alloc := newTestAllocator()
		route := NewRoute(RouteParams{
			TicketAllocator: alloc,
		})

		rec := newRecorder()
		route.recorder = rec

		t1 := route.Ticket()
		t2 := route.Ticket()
		t3 := route.Ticket()

		var wg sync.WaitGroup
		wg.Add(3)
		go func() {
			time.Sleep(time.Millisecond * 10)
			require.NoError(t, route.CompleteTicket(context.TODO(), CompleteTicketParams{Ticket: t2}))
			wg.Done()
		}()

		go func() {
			time.Sleep(time.Millisecond * 20)
			require.NoError(t, route.CompleteTicket(context.TODO(), CompleteTicketParams{Ticket: t3}))
			wg.Done()
		}()

		go func() {
			time.Sleep(time.Millisecond * 30)
			require.NoError(t, route.CompleteTicket(context.TODO(), CompleteTicketParams{Ticket: t1}))
			wg.Done()
		}()

		wg.Wait()
		require.Equal(t, rec.compCalls, []*Ticket{t2, t3, t1})
		require.Equal(t, rec.takeCalls, rec.completed)
		require.Nil(t, route.last)
		require.Equal(t, rec.completed, alloc.released)
	})
}

func BenchmarkRoute(b *testing.B) {
	ctx := context.Background()
	b.Run("alloc-std", func(b *testing.B) {
		route := NewRoute(RouteParams{})
		for i := 0; i < b.N; i++ {
			var (
				t1 = route.Ticket()
				t2 = route.Ticket()
				t3 = route.Ticket()
			)

			require.NoError(b, route.CompleteTicket(ctx, CompleteTicketParams{Ticket: t1}))
			require.NoError(b, route.CompleteTicket(ctx, CompleteTicketParams{Ticket: t2}))
			require.NoError(b, route.CompleteTicket(ctx, CompleteTicketParams{Ticket: t3}))
		}
	})

	b.Run("alloc-syncpool", func(b *testing.B) {
		route := NewRoute(RouteParams{
			TicketAllocator: NewSyncpoolAllocator(),
		})
		for i := 0; i < b.N; i++ {
			var (
				t1 = route.Ticket()
				t2 = route.Ticket()
				t3 = route.Ticket()
			)

			require.NoError(b, route.CompleteTicket(ctx, CompleteTicketParams{Ticket: t1}))
			require.NoError(b, route.CompleteTicket(ctx, CompleteTicketParams{Ticket: t2}))
			require.NoError(b, route.CompleteTicket(ctx, CompleteTicketParams{Ticket: t3}))
		}
	})
}

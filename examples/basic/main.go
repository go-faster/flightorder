package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-faster/flightorder"
)

func main() {
	// Create a route.
	route := flightorder.NewRoute(flightorder.RouteParams{})

	// Take some tickets.
	t1 := route.Ticket()
	t2 := route.Ticket()
	t3 := route.Ticket()

	// Perform parallel processing.
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond * 20)
		fmt.Println("Task 1 started")
		route.CompleteTicket(context.TODO(), flightorder.CompleteTicketParams{
			Ticket: t1,
			Completion: func(ctx context.Context) error {
				fmt.Println("Task 1 completed")
				return nil
			},
		})
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond * 30)
		fmt.Println("Task 2 started")
		route.CompleteTicket(context.TODO(), flightorder.CompleteTicketParams{
			Ticket: t2,
			Completion: func(ctx context.Context) error {
				fmt.Println("Task 2 completed")
				return nil
			},
		})
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond * 10)
		fmt.Println("Task 3 started")
		route.CompleteTicket(context.TODO(), flightorder.CompleteTicketParams{
			Ticket: t3,
			Completion: func(ctx context.Context) error {
				fmt.Println("Task 3 completed")
				return nil
			},
		})
	}()

	wg.Wait()
}

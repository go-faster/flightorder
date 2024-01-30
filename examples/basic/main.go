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
	t1 := route.TakeTicket()
	t2 := route.TakeTicket()
	t3 := route.TakeTicket()

	// Perform parallel processing.
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		time.Sleep(time.Millisecond * 20)
		fmt.Println("Task 1 started")
		route.CompleteTicket(context.TODO(), t1)
		fmt.Println("Task 1 completed")
		wg.Done()
	}()

	go func() {
		time.Sleep(time.Millisecond * 30)
		fmt.Println("Task 2 started")
		route.CompleteTicket(context.TODO(), t2)
		fmt.Println("Task 2 completed")
		wg.Done()
	}()

	go func() {
		time.Sleep(time.Millisecond * 10)
		fmt.Println("Task 3 started")
		route.CompleteTicket(context.TODO(), t3)
		fmt.Println("Task 3 completed")
		wg.Done()
	}()

	wg.Wait()
}

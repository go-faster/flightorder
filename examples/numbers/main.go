package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/go-faster/flightorder"
)

func main() {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	processingOrder, output := processInput(context.TODO(), input)
	fmt.Printf("input:     %v\n", input)
	fmt.Printf("processed: %v\n", processingOrder)
	fmt.Printf("output:    %v\n", output)
}

func processInput(ctx context.Context, input []int) (processing, output []int) {
	route := flightorder.NewRoute(flightorder.RouteParams{})

	var (
		mux sync.Mutex
		wg  sync.WaitGroup
	)

	wg.Add(len(input))
	for _, v := range input {
		ticket := route.Ticket()
		go func(t *flightorder.Ticket, v int) {
			defer wg.Done()
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))

			mux.Lock()
			processing = append(processing, v)
			mux.Unlock()

			_ = route.CompleteTicket(ctx, t, func(context.Context) error {
				mux.Lock()
				output = append(output, v)
				mux.Unlock()
				return nil
			})
		}(ticket, v)
	}

	wg.Wait()
	return
}

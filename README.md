# flightorder [![GoDoc](https://godoc.org/github.com/go-faster/flightorder?status.svg)](https://godoc.org/github.com/go-faster/flightorder)

This package allows to do _[ordered input] -> [parallel processing] -> [ordered output]_ in a streaming manner.

The name was inspired by [golang.org/x/sync/singleflight](https://pkg.go.dev/golang.org/x/sync/singleflight) package.

## Motivation

Sending logs from a single file to multiple kafka brokers concurrently while preserving at-least-once delivery guarantees:
* Logs are sent to multiple kafka brokers in parallel to enhance throughput.
* File offsets are committed in the exact order they are read to ensure at-least-once delivery guarantees and prevent data loss in case of shipper or broker failures.

## Installation

```
go get github.com/go-faster/flightorder@latest
```

## Example

```go
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
	processingOrder, output := processInput(input)
	fmt.Printf("input:     %v\n", input)
	fmt.Printf("processed: %v\n", processingOrder)
	fmt.Printf("output:    %v\n", output)
}

func processInput(input []int) (processing, output []int) {
	route := flightorder.NewRoute(flightorder.RouteParams{})

	var (
		mux sync.Mutex
		wg  sync.WaitGroup
	)

	wg.Add(len(input))
	for _, v := range input {
		ticket := route.TakeTicket()
		go func(t *flightorder.Ticket, v int) {
			defer wg.Done()
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))

			mux.Lock()
			processing = append(processing, v)
			mux.Unlock()

			_ = route.CompleteTicket(context.TODO(), t)

			mux.Lock()
			output = append(output, v)
			mux.Unlock()
		}(ticket, v)
	}

	wg.Wait()
	return
}
```

Output:
```
input:     [1 2 3 4 5 6 7 8 9]
processed: [3 1 9 7 6 2 5 4 8]
output:    [1 2 3 4 5 6 7 8 9]
```

## License

Source code is available under Apache License 2.0
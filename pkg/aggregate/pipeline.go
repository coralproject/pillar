package aggregate

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"sync"

	"golang.org/x/net/context"
)

// Accumultor is an object that can accumulate statistics based on individual
// objects within a context.
type Accumulator interface {
	Accumulate(context.Context, interface{})
	Combine(interface{})
}

// Pipeline runs a concurrent processing pipeline that executes a given
// function for every value read off of an input channel. The KeyValue object
// provided to the processor function is only accessable to one of the go
// routines within the pipeline.
func Pipeline(
	ctx context.Context,
	in chan interface{},

	// newAccumulator will be used to produce a single accumulator object for
	// each Go routine withine the pipeline.
	newAccumulator func() Accumulator,
) Accumulator {
	uid := fmt.Sprintf("%8X", rand.Uint32())

	// Keep track of GOMAXPROCS accumulators using a slice.
	gomaxprocs := runtime.GOMAXPROCS(-1)
	accumulators := make([]Accumulator, gomaxprocs)

	// Start GOMAXPROCS Go routines to read values from the input channel; we'll
	// track their execution using a sync.WaitGroup instance.
	var waitGroup sync.WaitGroup
	waitGroup.Add(gomaxprocs)
	for i := 0; i < gomaxprocs; i++ {

		// Use the provided newAccumulator function to produce an accumulator
		// value for this Go routine and add it to a new context.Context instance.
		accumulator := newAccumulator()
		accumulators[i] = accumulator

		// Start a new Go routine passing the index value for logging.
		go func(i int) {

			// Process values until the input channel is closed, then signal this go
			// routine has finished processing.
			for object := range in {
				accumulator.Accumulate(ctx, object)
			}
			waitGroup.Done()
		}(i)
	}

	// Wait for the processors to exit.
	waitGroup.Wait()

	// Combine the accumulators.
	accumulator := accumulators[0]
	for i := 1; i < len(accumulators); i++ {
		accumulator.Combine(accumulators[i])
	}

	// Return a single accumulator.
	return accumulator
}

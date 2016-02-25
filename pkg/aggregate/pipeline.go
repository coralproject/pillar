package aggregate

import (
	"log"
	"runtime"
	"sync"

	"golang.org/x/net/context"
)

// Pipeline runs a concurrent processing pipeline that executes a given
// function for every value read off of an input channel. The KeyValue object
// provided to the processor function is only accessable to one of the go
// routines within the pipeline.
func Pipeline(
	ctx context.Context,
	in chan interface{},

	// newAccumulator will be used to produce a single accumulator object for
	// each Go routine withine the pipeline.
	newAccumulator func(context.Context) interface{},

	// processor is a method that will be called for each object read from the
	// input channel. An accumulator will be available under the key
	// "accumulator" in the provided context.Context instance.
	processor func(context.Context, interface{}),
) []interface{} {

	// Keep track of GOMAXPROCS accumulators using a slice.
	gomaxprocs := runtime.GOMAXPROCS(-1)
	accumulators := make([]interface{}, gomaxprocs)

	// Start GOMAXPROCS Go routines to read values from the input channel; we'll
	// track their execution using a sync.WaitGroup instance.
	var waitGroup sync.WaitGroup
	waitGroup.Add(gomaxprocs)
	for i := 0; i < gomaxprocs; i++ {

		// Use the provided newAccumulator function to produce an accumulator
		// value for this Go routine and add it to a new context.Context instance.
		accumulators[i] = newAccumulator(ctx)
		accumulatorCtx := context.WithValue(ctx, "accumulator", accumulators[i])

		// Start a new Go routine passing the index value for logging.
		go func(i int) {

			// Process values until the input channel is closed, then signal this go
			// routine has finished processing.
			total := 0
			for object := range in {
				processor(accumulatorCtx, object)
				total++
				if total%1000 == 0 {
					log.Printf("routine[%d] processed %d", i, total)
				}
			}
			waitGroup.Done()
		}(i)
	}

	// Wait for the processors to exit, then return the accumulators.œ
	waitGroup.Wait()
	return accumulators
}

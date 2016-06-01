package calc

import (
	"log"

	"golang.org/x/net/context"

	"github.com/coralproject/pillar/pkg/aggregate"
	"github.com/coralproject/pillar/pkg/backend"
	"github.com/coralproject/pillar/pkg/backend/iterator"
	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/model/statistics"
)

// CalculateUserStatistics calculates User Statistics , creating a collection user_statistics with a document per user
func CalculateUserStatistics(ctx context.Context) error {

	// Look for a backen in the context and return an error if one is not
	// present.
	b, ok := ctx.Value("backend").(backend.Backend)
	if !ok {
		return backend.BackendNotInitializedError
	}

	// Get the users iterator.
	iter, err := b.Find("users", nil)
	if err != nil {
		return err
	}

	// Pipeline expects a generic input channel.
	in := make(chan interface{})

	go func() {
		defer close(in)
		if err := iterator.Each(iter, func(doc interface{}) error {

			// Assert that the document is the type we expect.
			user, ok := doc.(*model.User)
			if !ok {
				return backend.BackendTypeError
			}

			in <- user
			return nil
		}); err != nil {
			log.Println("CalculateUserStatistics error:", err)
			return
		}
	}()

	accumulator := aggregate.Pipeline(ctx, in, func() aggregate.Accumulator {
		return statistics.NewUserAccumulator()
	})

	if userAccumulator, ok := accumulator.(*statistics.UserAccumulator); ok {
		for _, dimension := range userAccumulator.Dimensions() {
			if err := b.Upsert("dimensions", map[string]interface{}{"name": dimension.Name}, dimension); err != nil {
				log.Println("CalculateUserStatistics error:", err)
			}
		}
	}

	return nil
}

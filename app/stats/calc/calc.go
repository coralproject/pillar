package calc

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/coralproject/pillar/pkg/aggregate"
	"github.com/coralproject/pillar/pkg/backend"
	"github.com/coralproject/pillar/pkg/backend/iterator"
	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/model/statistics"
)

func CalculateCommentStatistics(ctx context.Context, collection, foreignKeyField string) error {

	b, ok := ctx.Value("backend").(backend.Backend)
	if !ok {
		return backend.BackendNotInitializedError
	}

	iter, err := b.Find(collection, nil)
	if err != nil {
		return err
	}

	in := make(chan interface{})

	go func() {
		defer close(in)

		// Get the unique values of the foreign keys.
		if err := iterator.Each(iter, func(doc interface{}) error {

			// Assert that the document is the type we expect.
			user, ok := doc.(*model.User)
			if !ok {
				return backend.BackendTypeError
			}

			in <- user
			return nil
		}); err != nil {
			fmt.Println("User error:", err)
			return
		}
	}()

	aggregate.Pipeline(ctx, in, func() aggregate.Accumulator {
		return statistics.NewUserAccumulator()
	})

	return nil
}

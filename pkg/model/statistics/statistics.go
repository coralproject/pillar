package statistics

import (
	"golang.org/x/net/context"
)

const (
	referenceOnlyKey = "referenceonly"
)

func NewReferenceOnlyContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, referenceOnlyKey, true)
}

func ReferenceOnlyFromContext(ctx context.Context) bool {
	value, ok := ctx.Value(referenceOnlyKey).(bool)
	return ok && value
}

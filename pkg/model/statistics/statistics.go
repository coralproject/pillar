package statistics

import (
	"golang.org/x/net/context"
)

const (
	referenceKey = "referenceonly"
)

func NewOmitReferencesContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, referenceKey, true)
}

func OmitReferencesFromContext(ctx context.Context) bool {
	value, ok := ctx.Value(referenceKey).(bool)
	return ok && value
}

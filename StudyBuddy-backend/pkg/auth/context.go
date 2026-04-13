package auth

import "context"

type contextKey string

const userContextKey contextKey = "user"

// WithUserID stores user ID in context (for use after JWT validation).
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userContextKey, userID)
}

// UserIDFromContext returns the user ID from context, or "" if missing.
func UserIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(userContextKey).(string)
	return v
}

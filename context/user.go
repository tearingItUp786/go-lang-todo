package context

import (
	"context"
	"errors"

	"github.com/tearingItUp786/go-lang-todo/models"
)

// need a type alias to avoid collisions with other keys in context
type key string

// use your own key type to avoid collisions with other context values
const (
	userKey key = "user"
)

// WithUser returns a new context with the given user attached.
func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	val := ctx.Value(userKey)
	user, ok := val.(*models.User)
	if !ok {
		// The most likely case is that nothing was ever stored in the context,
		// so it doesn't have a type of *models.User. It is also possible that
		// other code in this package wrote an invalid value using the user key.
		return nil
	}
	return user
}

func GetUserId(ctx context.Context) (int, error) {
	user := User(ctx)
	if user == nil {
		return -1, errors.New("user not found")
	}
	return user.ID, nil
}

package userservice

import (
	"context"
)

type UserServiceI interface {
	ModerateUser(ctx context.Context)
	ViewBlockedUsers(ctx context.Context)
}

package ports

import (
	"context"
	"tweetschallenge/internal/domain"
)

type FollowRepo interface {
	Create(ctx context.Context, f *domain.Follow) error
	Unfollow(ctx context.Context, followerID, followeeID string) error
	FollowingIDs(ctx context.Context, followerID string) ([]string, error)
	FollowersIDs(ctx context.Context, followeeID string) ([]string, error)
}

package ports

import (
	"context"

	"tweetschallenge/internal/domain"
)

type TweetRepo interface {
	Create(ctx context.Context, t *domain.Tweet) error
	Timeline(ctx context.Context, userID string, limit, offset int) ([]domain.Tweet, error)
	TimelineForUsers(ctx context.Context, userIDs []string, limit, offset int) ([]domain.Tweet, error)
}

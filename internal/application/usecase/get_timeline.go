package usecase

import (
	"context"
	"fmt"

	"tweetschallenge/internal/domain"
	"tweetschallenge/internal/ports"
)

type GetTimeline struct {
	Tweets  ports.TweetRepo
	Follows ports.FollowRepo
}

type GetTimelineInput struct {
	UserID string
	Limit  int
	Offset int
}

func (uc GetTimeline) Exec(ctx context.Context, in GetTimelineInput) ([]domain.Tweet, error) {
	if uc.Follows == nil {
		return nil, fmt.Errorf("timeline: follow repo not wired")
	}
	ids, err := uc.Follows.FollowingIDs(ctx, in.UserID)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return []domain.Tweet{}, nil
	}
	return uc.Tweets.TimelineForUsers(ctx, ids, in.Limit, in.Offset)
}

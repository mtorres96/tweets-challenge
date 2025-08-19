package usecase

import (
	"context"
	"tweetschallenge/internal/domain"
	"tweetschallenge/internal/ports"
)

type FollowUser struct {
	Follows ports.FollowRepo
	Clock   ports.Clock
	IDGen   ports.IDGen
}

type FollowUserInput struct{ FollowerID, FolloweeID string }

func (uc FollowUser) Exec(ctx context.Context, in FollowUserInput) (domain.Follow, error) {
	f, err := domain.NewFollow(uc.IDGen.NewID(), in.FollowerID, in.FolloweeID, uc.Clock.NowUnix())
	if err != nil {
		return domain.Follow{}, err
	}
	if err := uc.Follows.Create(ctx, &f); err != nil {
		return domain.Follow{}, err
	}
	return f, nil
}

package usecase

import (
	"context"
	"tweetschallenge/internal/ports"
)

type UnfollowUser struct{ Follows ports.FollowRepo }

type UnfollowUserInput struct{ FollowerID, FolloweeID string }

func (uc UnfollowUser) Exec(ctx context.Context, in UnfollowUserInput) error {
	return uc.Follows.Unfollow(ctx, in.FollowerID, in.FolloweeID)
}

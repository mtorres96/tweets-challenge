package usecase

import (
	"context"

	"tweetschallenge/internal/domain"
	"tweetschallenge/internal/ports"
)

type PostTweet struct {
	Tweets ports.TweetRepo
	Clock  ports.Clock
	IDGen  ports.IDGen
}

type PostTweetInput struct{ UserID, Text string }

func (uc PostTweet) Exec(ctx context.Context, in PostTweetInput) (domain.Tweet, error) {
	id := uc.IDGen.NewID()
	now := uc.Clock.NowUnix()
	tw, err := domain.NewTweet(id, in.UserID, in.Text, now)
	if err != nil {
		return domain.Tweet{}, err
	}
	if err := uc.Tweets.Create(ctx, &tw); err != nil {
		return domain.Tweet{}, err
	}
	return tw, nil
}

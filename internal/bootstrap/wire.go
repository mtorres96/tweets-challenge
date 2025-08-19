package bootstrap

import (
	"fmt"

	"github.com/gin-gonic/gin"

	adapterclock "tweetschallenge/internal/adapters/clock"
	adaptersdb "tweetschallenge/internal/adapters/db"
	adaptershttp "tweetschallenge/internal/adapters/http"
	adapterid "tweetschallenge/internal/adapters/id"
	app "tweetschallenge/internal/application/usecase"
)

func BuildHTTPServer() (*gin.Engine, func(), error) {
	// DB (in-memory SQLite)
	db, err := adaptersdb.NewInMemoryGorm()
	if err != nil {
		return nil, nil, fmt.Errorf("db: %w", err)
	}
	if err := adaptersdb.AutoMigrate(db); err != nil {
		return nil, nil, fmt.Errorf("migrate: %w", err)
	}
	if err := adaptersdb.AutoMigrateFollow(db); err != nil {
		return nil, nil, fmt.Errorf("migrate follow: %w", err)
	}

	// Adapters
	tweetRepo := adaptersdb.NewTweetRepoGorm(db)
	followRepo := adaptersdb.NewFollowRepoGorm(db)
	clock := adapterclock.SystemClock{}
	idgen := adapterid.ULID{}
	limiter := adaptershttp.NewRateLimiterFromEnv()

	// Use cases
	postTweet := app.PostTweet{Tweets: tweetRepo, Clock: clock, IDGen: idgen}
	getTimeline := app.GetTimeline{Tweets: tweetRepo, Follows: followRepo}
	followUser := app.FollowUser{Follows: followRepo, Clock: clock, IDGen: idgen}
	unfollowUser := app.UnfollowUser{Follows: followRepo}

	// HTTP
	h := adaptershttp.BuildHandlers(postTweet, getTimeline, followUser, unfollowUser, limiter)
	r := adaptershttp.NewRouter(h)

	shutdown := func() {}
	return r, shutdown, nil
}

package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"tweetschallenge/internal/application/usecase"
)

type Handlers struct {
	Tweet  TweetHandler
	Follow FollowHandler
}

func NewRouter(h Handlers) *gin.Engine {
	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/v1")
	{
		// Tweets
		api.POST("/tweets", h.Tweet.Create)
		api.GET("/timeline/:userID", h.Tweet.Timeline)

		// Follows (solo follow/unfollow)
		api.POST("/follows", h.Follow.Create)
		api.DELETE("/follows", h.Follow.Delete)
	}
	return r
}

func BuildHandlers(
	postTweet usecase.PostTweet,
	getTimeline usecase.GetTimeline,
	followUser usecase.FollowUser,
	unfollowUser usecase.UnfollowUser,
	limiter *RateLimiter,
) Handlers {
	return Handlers{
		Tweet:  TweetHandler{PostTweet: postTweet, GetTimeline: getTimeline, Limiter: limiter},
		Follow: FollowHandler{FollowUser: followUser, UnfollowUser: unfollowUser},
	}
}

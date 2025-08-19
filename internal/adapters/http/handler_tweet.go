package http

import (
	"net/http"
	"strconv"

	"tweetschallenge/internal/application/usecase"

	"github.com/gin-gonic/gin"
)

type TweetHandler struct {
	PostTweet   usecase.PostTweet
	GetTimeline usecase.GetTimeline
	Limiter     *RateLimiter // rate limit por user_id
}

type CreateTweetReq struct {
	UserID string `json:"user_id" binding:"required"`
	Text   string `json:"text"   binding:"required"`
}

// @Summary Create tweet
// @Tags tweets
// @Accept json
// @Produce json
// @Param payload body CreateTweetReq true "payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 429 {object} map[string]string
// @Router /v1/tweets [post]
func (h TweetHandler) Create(c *gin.Context) {
	var req CreateTweetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	// Rate limit por usuario
	if h.Limiter != nil && !h.Limiter.Allow(req.UserID) {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
		return
	}

	tw, err := h.PostTweet.Exec(c, usecase.PostTweetInput{UserID: req.UserID, Text: req.Text})
	if err != nil {
		c.JSON(422, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": tw})
}

// @Summary Timeline
// @Tags tweets
// @Produce json
// @Param userID path string true "user id"
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Success 200 {object} map[string]interface{}
// @Router /v1/timeline/{userID} [get]
func (h TweetHandler) Timeline(c *gin.Context) {
	userID := c.Param("userID")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	list, err := h.GetTimeline.Exec(c, usecase.GetTimelineInput{
		UserID: userID, Limit: limit, Offset: offset,
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": list})
}

package http

import (
	"net/http"

	"tweetschallenge/internal/application/usecase"

	"github.com/gin-gonic/gin"
)

type FollowHandler struct {
	FollowUser   usecase.FollowUser
	UnfollowUser usecase.UnfollowUser
}

type FollowReq struct {
	FollowerID string `json:"follower_id" binding:"required"`
	FolloweeID string `json:"followee_id" binding:"required"`
}

// @Summary Follow user
// @Tags follows
// @Accept json
// @Produce json
// @Param payload body FollowReq true "payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Router /v1/follows [post]
func (h FollowHandler) Create(c *gin.Context) {
	var req FollowReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	f, err := h.FollowUser.Exec(c, usecase.FollowUserInput{
		FollowerID: req.FollowerID, FolloweeID: req.FolloweeID,
	})
	if err != nil {
		c.JSON(422, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": f})
}

// @Summary Unfollow user
// @Tags follows
// @Accept json
// @Produce json
// @Param payload body FollowReq true "payload"
// @Success 204 {string} string ""
// @Failure 400 {object} map[string]string
// @Router /v1/follows [delete]
func (h FollowHandler) Delete(c *gin.Context) {
	var req FollowReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	if err := h.UnfollowUser.Exec(c, usecase.UnfollowUserInput{
		FollowerID: req.FollowerID, FolloweeID: req.FolloweeID,
	}); err != nil {
		c.JSON(422, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

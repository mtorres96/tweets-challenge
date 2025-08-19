package domain

import "errors"

type Follow struct {
	ID         string `json:"id"`
	FollowerID string `json:"follower_id"`
	FolloweeID string `json:"followee_id"`
	CreatedAt  int64  `json:"created_at"`
}

func NewFollow(id, followerID, followeeID string, createdAt int64) (Follow, error) {
	if followerID == "" || followeeID == "" {
		return Follow{}, errors.New("both ids required")
	}
	if followerID == followeeID {
		return Follow{}, errors.New("cannot follow self")
	}
	return Follow{ID: id, FollowerID: followerID, FolloweeID: followeeID, CreatedAt: createdAt}, nil
}

package db

import (
	"context"

	"tweetschallenge/internal/domain"

	"gorm.io/gorm"
)

type FollowModel struct {
	ID         string `gorm:"primaryKey"`
	FollowerID string `gorm:"index:idx_follow_pair,unique"`
	FolloweeID string `gorm:"index:idx_follow_pair,unique"`
	CreatedAt  int64  `gorm:"index"`
}

type FollowRepoGorm struct{ db *gorm.DB }

func AutoMigrateFollow(db *gorm.DB) error          { return db.AutoMigrate(&FollowModel{}) }
func NewFollowRepoGorm(db *gorm.DB) FollowRepoGorm { return FollowRepoGorm{db: db} }

func (r FollowRepoGorm) Create(ctx context.Context, f *domain.Follow) error {
	var m FollowModel
	err := r.db.WithContext(ctx).
		Where(&FollowModel{FollowerID: f.FollowerID, FolloweeID: f.FolloweeID}).
		First(&m).Error
	if err == nil {
		return nil
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	m = FollowModel{ID: f.ID, FollowerID: f.FollowerID, FolloweeID: f.FolloweeID, CreatedAt: f.CreatedAt}
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r FollowRepoGorm) Unfollow(ctx context.Context, followerID, followeeID string) error {
	return r.db.WithContext(ctx).
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Delete(&FollowModel{}).Error
}

func (r FollowRepoGorm) FollowingIDs(ctx context.Context, followerID string) ([]string, error) {
	var rows []FollowModel
	if err := r.db.WithContext(ctx).Where("follower_id = ?", followerID).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]string, 0, len(rows))
	for _, m := range rows {
		out = append(out, m.FolloweeID)
	}
	return out, nil
}

func (r FollowRepoGorm) FollowersIDs(ctx context.Context, followeeID string) ([]string, error) {
	var rows []FollowModel
	if err := r.db.WithContext(ctx).Where("followee_id = ?", followeeID).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]string, 0, len(rows))
	for _, m := range rows {
		out = append(out, m.FollowerID)
	}
	return out, nil
}

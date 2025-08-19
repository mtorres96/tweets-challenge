package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"tweetschallenge/internal/domain"
)

func sqliteDSN() string {
	if dsn := os.Getenv("SQLITE_DSN"); dsn != "" {
		return dsn
	}
	// nombre único; cache=shared permite que el pool interno comparta la misma DB
	name := fmt.Sprintf("mem_%d", time.Now().UnixNano())
	return fmt.Sprintf("file:%s?mode=memory&cache=shared&_fk=1", name)
}

func NewInMemoryGorm() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(sqliteDSN()), &gorm.Config{})
}

// ----------------------------------------------------------------------------
// Model & Repo
// ----------------------------------------------------------------------------

type TweetModel struct {
	ID        string `gorm:"primaryKey"`
	UserID    string `gorm:"index"`
	Text      string `gorm:"size:280"`
	CreatedAt int64  `gorm:"index"`
}

type TweetRepoGorm struct{ db *gorm.DB }

func AutoMigrate(db *gorm.DB) error              { return db.AutoMigrate(&TweetModel{}) }
func NewTweetRepoGorm(db *gorm.DB) TweetRepoGorm { return TweetRepoGorm{db: db} }

func (r TweetRepoGorm) toDomain(m TweetModel) domain.Tweet {
	return domain.Tweet{
		ID:        m.ID,
		UserID:    m.UserID,
		Text:      m.Text,
		CreatedAt: m.CreatedAt,
	}
}

func (r TweetRepoGorm) Create(ctx context.Context, t *domain.Tweet) error {
	m := TweetModel{ID: t.ID, UserID: t.UserID, Text: t.Text, CreatedAt: t.CreatedAt}
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r TweetRepoGorm) Timeline(ctx context.Context, userID string, limit, offset int) ([]domain.Tweet, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	var rows []TweetModel
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]domain.Tweet, 0, len(rows))
	for _, m := range rows {
		out = append(out, r.toDomain(m))
	}
	return out, nil
}

func (r TweetRepoGorm) TimelineForUsers(ctx context.Context, userIDs []string, limit, offset int) ([]domain.Tweet, error) {
	if len(userIDs) == 0 {
		return []domain.Tweet{}, nil
	}
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	var rows []TweetModel
	err := r.db.WithContext(ctx).
		Where("user_id IN ?", userIDs).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]domain.Tweet, 0, len(rows))
	for _, m := range rows {
		out = append(out, r.toDomain(m))
	}
	return out, nil
}

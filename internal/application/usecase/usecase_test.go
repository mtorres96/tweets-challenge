package usecase

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"tweetschallenge/internal/domain"
	"tweetschallenge/internal/ports"
)

// ===== fakes =====

type fakeClock struct{ now int64 }

func (f fakeClock) NowUnix() int64 { return f.now }

type fakeID struct{ id string }

func (f fakeID) NewID() string { return f.id }

type memTweetRepo struct {
	created []domain.Tweet
	byUser  map[string][]domain.Tweet
}

func (m *memTweetRepo) Create(ctx context.Context, t *domain.Tweet) error {
	m.created = append(m.created, *t)
	m.byUser[t.UserID] = append(m.byUser[t.UserID], *t)
	return nil
}
func (m *memTweetRepo) Timeline(ctx context.Context, userID string, limit, offset int) ([]domain.Tweet, error) {
	return m.byUser[userID], nil
}
func (m *memTweetRepo) TimelineForUsers(ctx context.Context, userIDs []string, limit, offset int) ([]domain.Tweet, error) {
	res := []domain.Tweet{}
	for _, id := range userIDs {
		res = append(res, m.byUser[id]...)
	}
	return res, nil
}

type memFollowRepo struct {
	following  map[string][]string // follower -> followees
	unfollowed [][2]string
	errCreate  error
}

func (m *memFollowRepo) Create(ctx context.Context, f *domain.Follow) error {
	if m.errCreate != nil {
		return m.errCreate
	}
	m.following[f.FollowerID] = append(m.following[f.FollowerID], f.FolloweeID)
	return nil
}
func (m *memFollowRepo) Unfollow(ctx context.Context, followerID, followeeID string) error {
	m.unfollowed = append(m.unfollowed, [2]string{followerID, followeeID})
	return nil
}
func (m *memFollowRepo) FollowingIDs(ctx context.Context, followerID string) ([]string, error) {
	return m.following[followerID], nil
}
func (m *memFollowRepo) FollowersIDs(ctx context.Context, followeeID string) ([]string, error) {
	return nil, nil
}

// ===== tests =====

func TestPostTweet_OK(t *testing.T) {
	tr := &memTweetRepo{byUser: map[string][]domain.Tweet{}}
	uc := PostTweet{Tweets: tr, Clock: fakeClock{now: 100}, IDGen: fakeID{id: "T1"}}

	got, err := uc.Exec(context.Background(), PostTweetInput{UserID: "u1", Text: "hola"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got.ID != "T1" || got.UserID != "u1" || got.Text != "hola" || got.CreatedAt != 100 {
		t.Fatalf("unexpected tweet: %#v", got)
	}
	if len(tr.created) != 1 {
		t.Fatalf("expected 1 created, got %d", len(tr.created))
	}
}

func TestPostTweet_Validation(t *testing.T) {
	tr := &memTweetRepo{byUser: map[string][]domain.Tweet{}}
	uc := PostTweet{Tweets: tr, Clock: fakeClock{now: 100}, IDGen: fakeID{id: "T1"}}
	if _, err := uc.Exec(context.Background(), PostTweetInput{UserID: "", Text: "x"}); err == nil {
		t.Fatal("expected error for empty user")
	}
	if _, err := uc.Exec(context.Background(), PostTweetInput{UserID: "u1", Text: ""}); err == nil {
		t.Fatal("expected error for empty text")
	}
}

func TestGetTimeline_OnlyFollowing(t *testing.T) {
	tr := &memTweetRepo{byUser: map[string][]domain.Tweet{
		"u1": {{ID: "A", UserID: "u1", Text: "mine", CreatedAt: 1}},
		"u2": {{ID: "B", UserID: "u2", Text: "from u2", CreatedAt: 2}},
	}}
	fr := &memFollowRepo{following: map[string][]string{"u1": {"u2"}}}
	uc := GetTimeline{Tweets: tr, Follows: fr}

	got, err := uc.Exec(context.Background(), GetTimelineInput{UserID: "u1", Limit: 50, Offset: 0})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	want := []domain.Tweet{{ID: "B", UserID: "u2", Text: "from u2", CreatedAt: 2}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v want %#v", got, want)
	}
}

func TestFollowUser_OK(t *testing.T) {
	fr := &memFollowRepo{following: map[string][]string{}}
	uc := FollowUser{Follows: fr, Clock: fakeClock{now: 50}, IDGen: fakeID{id: "F1"}}
	f, err := uc.Exec(context.Background(), FollowUserInput{FollowerID: "u1", FolloweeID: "u2"})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if f.ID != "F1" || f.FollowerID != "u1" || f.FolloweeID != "u2" || f.CreatedAt != 50 {
		t.Fatalf("unexpected follow: %#v", f)
	}
}

func TestFollowUser_DomainError(t *testing.T) {
	fr := &memFollowRepo{}
	uc := FollowUser{Follows: fr, Clock: fakeClock{now: 50}, IDGen: fakeID{id: "F1"}}
	if _, err := uc.Exec(context.Background(), FollowUserInput{FollowerID: "u1", FolloweeID: "u1"}); err == nil {
		t.Fatal("expected error for self follow")
	}
}

func TestFollowUser_RepoError(t *testing.T) {
	fr := &memFollowRepo{errCreate: errors.New("boom")}
	uc := FollowUser{Follows: fr, Clock: fakeClock{now: 50}, IDGen: fakeID{id: "F1"}}
	if _, err := uc.Exec(context.Background(), FollowUserInput{FollowerID: "u1", FolloweeID: "u2"}); err == nil {
		t.Fatal("expected repo error")
	}
}

func TestUnfollowUser_OK(t *testing.T) {
	fr := &memFollowRepo{following: map[string][]string{}}
	uc := UnfollowUser{Follows: fr}
	if err := uc.Exec(context.Background(), UnfollowUserInput{FollowerID: "u1", FolloweeID: "u2"}); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(fr.unfollowed) != 1 || fr.unfollowed[0][0] != "u1" || fr.unfollowed[0][1] != "u2" {
		t.Fatalf("unexpected unfollow calls: %#v", fr.unfollowed)
	}
}

// Interface assertions (por si cambiamos firmas sin querer)
var _ ports.TweetRepo = (*memTweetRepo)(nil)
var _ ports.FollowRepo = (*memFollowRepo)(nil)

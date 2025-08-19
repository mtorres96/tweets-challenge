package domain

import "testing"

func TestNewFollow_OK(t *testing.T) {
	f, err := NewFollow("f1", "u1", "u2", 42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.FollowerID != "u1" || f.FolloweeID != "u2" {
		t.Fatalf("unexpected follow: %#v", f)
	}
}

func TestNewFollow_Self(t *testing.T) {
	_, err := NewFollow("f1", "u1", "u1", 1)
	if err == nil {
		t.Fatal("expected error for self-follow")
	}
}

func TestNewFollow_Missing(t *testing.T) {
	if _, err := NewFollow("f1", "", "u2", 1); err == nil {
		t.Fatal("expected error for missing follower")
	}
	if _, err := NewFollow("f1", "u1", "", 1); err == nil {
		t.Fatal("expected error for missing followee")
	}
}

package domain

import "testing"

func TestNewTweet_OK(t *testing.T) {
	tw, err := NewTweet("id1", "u1", "hola", 123)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tw.ID != "id1" || tw.UserID != "u1" || tw.Text != "hola" || tw.CreatedAt != 123 {
		t.Fatalf("unexpected tweet: %#v", tw)
	}
}

func TestNewTweet_EmptyUser(t *testing.T) {
	_, err := NewTweet("id1", "", "hola", 1)
	if err == nil {
		t.Fatal("expected error for empty user_id")
	}
}

func TestNewTweet_TextLength(t *testing.T) {
	_, err := NewTweet("id1", "u1", "", 1)
	if err == nil {
		t.Fatal("expected error for empty text")
	}
	long := make([]byte, 281)
	for i := range long {
		long[i] = 'a'
	}
	_, err = NewTweet("id1", "u1", string(long), 1)
	if err == nil {
		t.Fatal("expected error for >280 chars")
	}
}

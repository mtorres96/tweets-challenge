package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"tweetschallenge/internal/bootstrap"
)

func doReq(r http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	var buf *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		buf = bytes.NewBuffer(b)
	} else {
		buf = &bytes.Buffer{}
	}
	req := httptest.NewRequest(method, path, buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestTimeline_OnlyFollowing(t *testing.T) {
	// Desactivar rate limit para no interferir
	os.Setenv("RATE_LIMIT_ENABLED", "false")
	defer os.Unsetenv("RATE_LIMIT_ENABLED")

	router, shutdown, err := bootstrap.BuildHTTPServer()
	if err != nil {
		t.Fatalf("build server: %v", err)
	}
	defer shutdown()

	// u1 y u2 postean, pero u1 sigue a u2; timeline de u1 debe traer SOLO de u2
	w := doReq(router, http.MethodPost, "/v1/follows", map[string]string{
		"follower_id": "u1", "followee_id": "u2",
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("follow status %d body %s", w.Code, w.Body.String())
	}

	w = doReq(router, http.MethodPost, "/v1/tweets", map[string]string{"user_id": "u1", "text": "mine"})
	if w.Code != http.StatusCreated {
		t.Fatalf("tweet u1 status %d body %s", w.Code, w.Body.String())
	}
	w = doReq(router, http.MethodPost, "/v1/tweets", map[string]string{"user_id": "u2", "text": "from u2"})
	if w.Code != http.StatusCreated {
		t.Fatalf("tweet u2 status %d body %s", w.Code, w.Body.String())
	}

	w = doReq(router, http.MethodGet, "/v1/timeline/u1", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("timeline status %d body %s", w.Code, w.Body.String())
	}
	type resp struct {
		Data []struct {
			UserID string `json:"user_id"`
			Text   string `json:"text"`
		} `json:"data"`
	}
	var out resp
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode: %v body=%s", err, w.Body.String())
	}
	if len(out.Data) != 1 || out.Data[0].UserID != "u2" {
		t.Fatalf("want only tweets from u2, got: %+v", out.Data)
	}
}

func TestRateLimit_CreateTweet(t *testing.T) {
	os.Setenv("RATE_LIMIT_ENABLED", "true")
	os.Setenv("RATE_LIMIT_MAX_TWEETS", "1")
	os.Setenv("RATE_LIMIT_WINDOW_SEC", "60")
	defer func() {
		os.Unsetenv("RATE_LIMIT_ENABLED")
		os.Unsetenv("RATE_LIMIT_MAX_TWEETS")
		os.Unsetenv("RATE_LIMIT_WINDOW_SEC")
	}()

	router, shutdown, err := bootstrap.BuildHTTPServer()
	if err != nil {
		t.Fatalf("build server: %v", err)
	}
	defer shutdown()

	// Primer tweet: 201
	w := doReq(router, http.MethodPost, "/v1/tweets", map[string]string{"user_id": "u1", "text": "hola"})
	if w.Code != http.StatusCreated {
		t.Fatalf("first tweet status %d body %s", w.Code, w.Body.String())
	}
	// Segundo en misma ventana: 429
	w = doReq(router, http.MethodPost, "/v1/tweets", map[string]string{"user_id": "u1", "text": "spam"})
	if w.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d body %s", w.Code, w.Body.String())
	}
}

func TestFollow_Idempotent(t *testing.T) {
	router, shutdown, err := bootstrap.BuildHTTPServer()
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	defer shutdown()

	body := map[string]string{"follower_id": "u1", "followee_id": "u2"}

	w := doReq(router, http.MethodPost, "/v1/follows", body)
	if w.Code != http.StatusCreated {
		t.Fatalf("1st follow got %d", w.Code)
	}

	// repetir mismo follow: debe seguir OK (idempotente)
	w = doReq(router, http.MethodPost, "/v1/follows", body)
	if w.Code != http.StatusCreated {
		t.Fatalf("2nd follow got %d", w.Code)
	}
}

func TestUnfollow_Idempotent(t *testing.T) {
	router, shutdown, err := bootstrap.BuildHTTPServer()
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	defer shutdown()

	// hacer unfollow de un par inexistente: 204 igual (idempotente)
	w := doReq(router, http.MethodDelete, "/v1/follows", map[string]string{
		"follower_id": "uX", "followee_id": "uY",
	})
	if w.Code != http.StatusNoContent {
		t.Fatalf("unfollow non-existing got %d body=%s", w.Code, w.Body.String())
	}
}

func TestCreateTweet_BadPayload(t *testing.T) {
	router, shutdown, err := bootstrap.BuildHTTPServer()
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	defer shutdown()

	// falta user_id y/o text => 400
	w := doReq(router, http.MethodPost, "/v1/tweets", map[string]string{})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("bad payload got %d body=%s", w.Code, w.Body.String())
	}
}

func TestRateLimit_Disabled_AllowsMultiple(t *testing.T) {
	os.Setenv("RATE_LIMIT_ENABLED", "false")
	defer os.Unsetenv("RATE_LIMIT_ENABLED")

	router, shutdown, err := bootstrap.BuildHTTPServer()
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	defer shutdown()

	// con RL deshabilitado, dos tweets del mismo user en la misma ventana deben ser 201 ambos
	w := doReq(router, http.MethodPost, "/v1/tweets", map[string]string{"user_id": "u1", "text": "a"})
	if w.Code != http.StatusCreated {
		t.Fatalf("1st tweet got %d", w.Code)
	}
	w = doReq(router, http.MethodPost, "/v1/tweets", map[string]string{"user_id": "u1", "text": "b"})
	if w.Code != http.StatusCreated {
		t.Fatalf("2nd tweet got %d", w.Code)
	}
}

func TestTimeline_EmptyWhenNoFollowing(t *testing.T) {
	os.Setenv("RATE_LIMIT_ENABLED", "false")
	defer os.Unsetenv("RATE_LIMIT_ENABLED")

	router, shutdown, err := bootstrap.BuildHTTPServer()
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	defer shutdown()

	w := doReq(router, http.MethodGet, "/v1/timeline/u1", nil)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d body=%s", w.Code, w.Body.String())
	}
	var out struct {
		Data []any `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(out.Data) != 0 {
		t.Fatalf("expected empty timeline, got %v", out.Data)
	}
}

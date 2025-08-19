package id

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

type ULID struct{}

func (ULID) NewID() string {
	entropy := ulid.Monotonic(rand.Reader, 0)
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}

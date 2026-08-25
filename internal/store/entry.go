package store

import "time"

type Entry struct {
	Value    interface{}
	ExpireAt time.Time
}

func (e *Entry) IsExpired() bool {
	if e.ExpireAt.IsZero() == false && e.ExpireAt.Before(time.Now()) {
		return true
	}
	return false
}

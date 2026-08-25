package command

import (
	"time"

	"github.com/surajsencha/redis-clone/internal/resp"
	"github.com/surajsencha/redis-clone/internal/store"
)

func TTL(s *store.Store) Handler {
	return func(args []resp.Value) (resp.Value, error) {
		if len(args) != 1 {
			return resp.Value{Type: resp.Error, Str: "ERR wrong number of arguments for 'ttl' command"}, nil
		}
		entry, ok := s.Get(args[0].Str)
		if !ok {
			return resp.Value{Type: resp.Integer, Num: -2}, nil
		}

		if entry.ExpireAt.IsZero() {
			return resp.Value{Type: resp.Integer, Num: -1}, nil
		}
		ttl := int64(time.Until(entry.ExpireAt).Seconds())
		return resp.Value{Type: resp.Integer, Num: ttl}, nil
	}
}

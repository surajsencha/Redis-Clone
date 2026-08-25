package command

import (
	"strconv"
	"strings"
	"time"

	"github.com/surajsencha/redis-clone/internal/resp"
	"github.com/surajsencha/redis-clone/internal/store"
)

func Set(s *store.Store) Handler {
	return func(args []resp.Value) (resp.Value, error) {
		var expireAt time.Time
		if len(args) < 2 {
			return resp.Value{Type: resp.Error, Str: "ERR wrong number of arguments for 'set' command"}, nil
		}
		if len(args) >= 4 && strings.ToUpper(args[2].Str) == "EX" {
			seconds, err := strconv.Atoi(args[3].Str)
			if err != nil {
				return resp.Value{Type: resp.Error, Str: "ERR value is not an integer or out of range"}, nil
			}
			expireAt = time.Now().Add(time.Duration(seconds) * time.Second)
		} else {
			expireAt = time.Time{}
		}

		s.Set(args[0].Str, args[1].Str, expireAt)
		return resp.Value{Type: resp.SimpleString, Str: "OK"}, nil
	}
}

func Get(s *store.Store) Handler {
	return func(args []resp.Value) (resp.Value, error) {
		if len(args) != 1 {
			return resp.Value{Type: resp.Error, Str: "ERR wrong number of arguments for 'get' command"}, nil
		}
		entry, ok := s.Get(args[0].Str)
		if !ok {
			return resp.Value{Type: resp.Null, Str: ""}, nil
		}
		return resp.Value{Type: resp.BulkString, Str: entry.Value.(string)}, nil
	}
}

func Del(s *store.Store) Handler {
	return func(args []resp.Value) (resp.Value, error) {
		if len(args) == 0 {
			return resp.Value{Type: resp.Error, Str: "ERR wrong number of arguments for 'del' command"}, nil
		}

		keys := make([]string, len(args))
		for i, arg := range args {
			keys[i] = arg.Str
		}
		deleted := s.Del(keys...)
		return resp.Value{Type: resp.Integer, Num: int64(deleted)}, nil
	}
}

package command

import (
	"github.com/surajsencha/redis-clone/internal/resp"
	"github.com/surajsencha/redis-clone/internal/store"
)

func Set(s *store.Store) Handler {
	return func(args []resp.Value) (resp.Value, error) {
		length := len(args)
		if length != 2 {
			return resp.Value{Type: resp.Error, Str: "ERR wrong number of arguments for 'set' command"}, nil
		}
		s.Set(args[0].Str, args[1].Str)
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

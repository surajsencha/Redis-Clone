package command

import (
	"strconv"

	"github.com/surajsencha/redis-clone/internal/resp"
	"github.com/surajsencha/redis-clone/internal/store"
)

func RPush(s *store.Store) Handler {
	return func(args []resp.Value) (resp.Value, error) {
		if len(args) < 2 {
			return resp.Value{Type: resp.Error, Str: "ERR wrong number of arguments for 'rpush' command"}, nil
		}
		values := []string{}
		for _, v := range args[1:] {
			values = append(values, v.Str)
		}
		count := s.RPush(args[0].Str, values...)
		return resp.Value{Type: resp.Integer, Num: int64(count)}, nil
	}
}

func LRange(s *store.Store) Handler {
	return func(args []resp.Value) (resp.Value, error) {
		if len(args) != 3 {
			return resp.Value{Type: resp.Error, Str: "ERR wrong number of arguments for 'lrange' command"}, nil
		}
		start, err1 := strconv.Atoi(args[1].Str)
		stop, err2 := strconv.Atoi(args[2].Str)
		if err1 != nil || err2 != nil {
			return resp.Value{Type: resp.Error, Str: "ERR value is not an integer or out of range"}, nil
		}

		list, ok := s.LRange(args[0].Str, int(start), int(stop))
		if !ok {
			return resp.Value{Type: resp.Array, Elems: []resp.Value{}}, nil
		}
		var listResp []resp.Value
		for _, item := range list {
			listResp = append(listResp, resp.Value{Type: resp.BulkString, Str: item})
		}
		return resp.Value{Type: resp.Array, Elems: listResp}, nil
	}
}

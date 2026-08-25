package command

import "github.com/surajsencha/redis-clone/internal/resp"

func Ping(args []resp.Value) (resp.Value, error) {
	if len(args) == 0 {
		return resp.Value{Type: resp.SimpleString, Str: "PONG"}, nil
	}
	return resp.Value{Type: resp.SimpleString, Str: args[0].Str}, nil
}

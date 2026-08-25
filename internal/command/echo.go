package command

import "github.com/surajsencha/redis-clone/internal/resp"

func Echo(args []resp.Value) (resp.Value, error) {
	if len(args) == 0 {
		return resp.Value{Type: resp.Error, Str: "ERR wrong number of arguments for 'echo' command"}, nil
	}
	return args[0], nil
}

package command

import (
	"strings"

	"github.com/surajsencha/redis-clone/internal/resp"
)

type Handler func(args []resp.Value) (resp.Value, error)

type Registry struct {
	handlers map[string]Handler
}

func NewRegistry() *Registry {
	return &Registry{
		handlers: make(map[string]Handler),
	}
}

func (r *Registry) Register(name string, handler Handler) {
	r.handlers[strings.ToUpper(name)] = handler
}

func (r *Registry) Execute(name string, args []resp.Value) resp.Value {
	h, ok := r.handlers[strings.ToUpper(name)]
	if !ok {
		return resp.Value{Type: resp.Error, Str: "ERR unknown command '" + name + "'"}
	}
	res, err := h(args)
	if err != nil {
		return resp.Value{Type: resp.Error, Str: "ERR " + err.Error()}
	}
	return res
}

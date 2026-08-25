package server

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/surajsencha/redis-clone/internal/command"
	"github.com/surajsencha/redis-clone/internal/config"
	"github.com/surajsencha/redis-clone/internal/resp"
	"github.com/surajsencha/redis-clone/internal/store"
)

type Server struct {
	config   config.Config
	listener net.Listener
	registry *command.Registry
	store    *store.Store
}

func NewServer(cfg *config.Config) *Server {
	store := store.NewStore()
	// 1. Create the registry
	reg := command.NewRegistry()

	// 2. Register the commands!
	reg.Register("PING", command.Ping)
	reg.Register("ECHO", command.Echo)
	reg.Register("SET", command.Set(store))
	reg.Register("GET", command.Get(store))
	reg.Register("DEL", command.Del(store))
	reg.Register("TTL", command.TTL(store))
	store.StartActiveExpiry()
	return &Server{
		config:   *cfg,
		registry: reg,
		store:    store,
	}
}

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen(
		"tcp",
		fmt.Sprintf(":%d", s.config.Port),
	)
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	s.listener = listener

	fmt.Printf("Listening on port %d\n", s.config.Port)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil // Normal shutdown, not an error
			}
			return fmt.Errorf("failed to accept connection: %w", err)
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	r := resp.NewReader(conn)
	w := resp.NewWriter(conn)
	for {
		value, err := resp.Read(r)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected")
			} else {
				fmt.Printf("Error reading from Connection: %v\n", err)
			}
			return
		}
		if value.Type != resp.Array || len(value.Elems) == 0 {
			continue
		}

		commandName := value.Elems[0].Str
		args := value.Elems[1:]

		reply := s.registry.Execute(commandName, args)
		err = w.Write(reply)
		if err != nil {
			fmt.Printf("Error writing to Connection: %v\n", err)
		}
		w.Flush()
	}
}
func (s *Server) Shutdown() error {
	return s.listener.Close()
}

package server

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/surajsencha/redis-clone/internal/config"
	"github.com/surajsencha/redis-clone/internal/resp"
)

type Server struct {
	config   config.Config
	listener net.Listener
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		config: *cfg,
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
		if value.Type == resp.Array && len(value.Elems) > 0 {
			commandName := value.Elems[0].Str
			if strings.ToUpper(commandName) == "PING" {
				reply := resp.Value{Type: resp.SimpleString, Str: "PONG"}
				err = w.Write(reply)
				w.Flush()
			}
		}
		if err != nil {
			fmt.Printf("Error writing to Connection: %v\n", err)
			return
		}
	}
}
func (s *Server) Shutdown() error {
	return s.listener.Close()
}

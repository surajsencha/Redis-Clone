package server

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/surajsencha/redis-clone/internal/config"
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
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected")
			} else {
				fmt.Printf("Error reading from Connection: %v\n", err)
			}
			return
		}
		fmt.Println("Received Data : ", string(buffer[:n]))
		_, err = conn.Write(buffer[:n])
		if err != nil {
			fmt.Printf("Error writing to Connection: %v\n", err)
			return
		}
	}
}
func (s *Server) Shutdown() error {
	return s.listener.Close()
}

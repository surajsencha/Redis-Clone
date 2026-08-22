package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/surajsencha/redis-clone/internal/config"
	"github.com/surajsencha/redis-clone/internal/server"
)

func main() {
	port := flag.Int("port", 6379, "Port for Redis server")
	flag.Parse()
	cfg := &config.Config{
		Host: "",
		Port: *port,
	}
	srv := server.NewServer(cfg)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Printf("Server stopped: %v\n", err)
		}
	}()
	sig := <-sigChan
	fmt.Printf("\nReceived %v. Shutting down...\n", sig)
	if err := srv.Shutdown(); err != nil {
		fmt.Printf("Error shutting down server: %v\n", err)
	}
}

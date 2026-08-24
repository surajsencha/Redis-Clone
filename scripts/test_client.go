package main

import (
	"fmt"
	"net"
)

func main() {
	// 1. Connect to your server (make sure the port matches!)
	conn, err := net.Dial("tcp", "localhost:6381")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// 2. Send the exact RESP bytes for the PING command
	fmt.Println("Sending PING...")
	conn.Write([]byte("*1\r\n$4\r\nPING\r\n"))

	// 3. Read the response from your server
	buffer := make([]byte, 1024)
	n, _ := conn.Read(buffer)

	// 4. Print it out! (We use %q to show the hidden \r\n characters)
	fmt.Printf("Server replied: %q\n", string(buffer[:n]))
}

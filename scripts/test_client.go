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
	fmt.Println("Sending just ECHO...")
	// *1 (array of 1) -> $4 ECHO
	conn.Write([]byte("*1\r\n$4\r\nECHO\r\n"))

	// 3. Read the response from your server
	buffer := make([]byte, 1024)
	n, _ := conn.Read(buffer)

	// 4. Print it out! (We use %q to show the hidden \r\n characters)
	fmt.Printf("Server replied: %q\n", string(buffer[:n]))
}

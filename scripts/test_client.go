package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:6381")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	buffer := make([]byte, 1024)

	// --- COMMENT THIS BLOCK OUT AFTER RUNNING ONCE ---

	fmt.Println("\n--- Test 1: SET name suraj ---")
	conn.Write([]byte("*3\r\n$3\r\nSET\r\n$4\r\nname\r\n$5\r\nsuraj\r\n"))
	n, _ := conn.Read(buffer)
	fmt.Printf("Expected: \"+OK\\r\\n\"\nGot:      %q\n", string(buffer[:n]))

	// ------------------------------------------------

	// --- TEST 2: GET the key ---
	fmt.Println("\n--- Test 2: GET name ---")
	conn.Write([]byte("*2\r\n$3\r\nGET\r\n$4\r\nname\r\n"))

	// BUG FIX: You commented out conn.Read! The client wasn't waiting for the server's reply!
	// n, _ := conn.Read(buffer)

	fmt.Printf("Expected: \"$5\\r\\nsuraj\\r\\n\"\nGot:      %q\n", string(buffer[:n]))
}

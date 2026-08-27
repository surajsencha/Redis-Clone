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

	// --- TEST 1: RPUSH ---
	fmt.Println("\n--- Test 1: RPUSH mylist a b c ---")
	conn.Write([]byte("*5\r\n$5\r\nRPUSH\r\n$6\r\nmylist\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n"))
	n, _ := conn.Read(buffer)
	fmt.Printf("Expected: \":3\\r\\n\"\nGot:      %q\n", string(buffer[:n]))

	// --- TEST 2: LRANGE ---
	fmt.Println("\n--- Test 2: LRANGE mylist 0 -1 ---")
	conn.Write([]byte("*4\r\n$6\r\nLRANGE\r\n$6\r\nmylist\r\n$1\r\n0\r\n$2\r\n-1\r\n"))
	n, _ = conn.Read(buffer)
	fmt.Printf("Got:      %q\n", string(buffer[:n]))
}

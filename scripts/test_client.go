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

	// --- TEST 1: SET ---
	fmt.Println("\n--- Test 1: SET ---")
	// *3 (array of 3) -> $3 SET -> $4 name -> $5 suraj
	conn.Write([]byte("*5\r\n$3\r\nSET\r\n$4\r\nname\r\n$5\r\nsuraj\r\n"))
	n, _ := conn.Read(buffer)
	fmt.Printf("Expected: \"+OK\\r\\n\"\nGot:      %q\n", string(buffer[:n]))

	// --- TEST 2: GET ---
	fmt.Println("\n--- Test 2: GET ---")
	// *2 (array of 2) -> $3 GET -> $4 name
	conn.Write([]byte("*2\r\n$3\r\nGET\r\n$4\r\nname\r\n"))
	n, _ = conn.Read(buffer)
	fmt.Printf("Expected: \"$5\\r\\nsuraj\\r\\n\"\nGot:      %q\n", string(buffer[:n]))

	// --- TEST 3: DEL ---
	fmt.Println("\n--- Test 3: DEL ---")
	// *2 (array of 2) -> $3 DEL -> $4 name
	conn.Write([]byte("*2\r\n$3\r\nDEL\r\n$4\r\nname\r\n"))
	n, _ = conn.Read(buffer)
	fmt.Printf("Expected: \":1\\r\\n\"\nGot:      %q\n", string(buffer[:n]))

	// --- TEST 4: GET (Should be missing now) ---
	fmt.Println("\n--- Test 4: GET (After Delete) ---")
	conn.Write([]byte("*2\r\n$3\r\nGET\r\n$4\r\nname\r\n"))
	n, _ = conn.Read(buffer)
	fmt.Printf("Expected: \"$-1\\r\\n\"\nGot:      %q\n", string(buffer[:n]))
}

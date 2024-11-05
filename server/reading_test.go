package main

import (
	"net"
	"os"
	"testing"
	"time"
)

// TestTCPClient tests the TCP server by connecting and sending messages
func TestTCPClient(t *testing.T) {
	// Connect to the TCP server
	conn, err := net.Dial("tcp", "localhost:8989")
	if err != nil {
		t.Fatalf("Error connecting to server: %v", err)
	}
	defer conn.Close()

	// Send the first message
	if _, err := conn.Write([]byte("Yenis\n")); err != nil {
		t.Fatalf("Error sending first message: %v", err)
	}

	// Send the second message
	if _, err := conn.Write([]byte("Hi (test msg 1)\n")); err != nil {
		t.Fatalf("Error sending second message: %v", err)
	}

	// Case 1: Close the connection immediately after sending messages
	t.Log("Closing connection...")
	if err := conn.Close(); err != nil {
		t.Fatalf("Error closing connection: %v", err)
	}

	// Reconnect for further tests
	conn, err = net.Dial("tcp", "localhost:8989")
	if err != nil {
		t.Fatalf("Error reconnecting to server: %v", err)
	}
	defer conn.Close()

	// Case 2: Send data and then close the connection
	if _, err := conn.Write([]byte("Yenis\n")); err != nil {
		t.Fatalf("Error sending message: %v", err)
	}
	time.Sleep(1 * time.Second) // Wait a bit before closing
	t.Log("Closing connection after sending Yenis...")
	if err := conn.Close(); err != nil {
		t.Fatalf("Error closing connection: %v", err)
	}

	// Reconnect for more tests
	conn, err = net.Dial("tcp", "localhost:8989")
	if err != nil {
		t.Fatalf("Error reconnecting to server: %v", err)
	}
	defer conn.Close()

	// Case 3: Send data and simulate a hang
	if _, err := conn.Write([]byte("Yenis\n")); err != nil {
		t.Fatalf("Error sending message: %v", err)
	}
	t.Log("Hanging connection... (not closing or sending more data)")
	time.Sleep(10 * time.Second) // Wait for some time; server will likely timeout

	// Case 4: Send data and then abruptly terminate the client
	t.Log("Abruptly terminating the client...")
	os.Exit(0) // This will cause an error on the server when it tries to read
}

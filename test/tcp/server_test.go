package tcp_test

import (
	"net"
	"testing"
	"time"

	"github.com/appyzdl/Netrunner/pkg/tcp"
)

func TestTCPServer(t *testing.T) {
	// start the server
	go tcp.StartServer("localhost:6969")

	time.Sleep(100 * time.Millisecond)

	// connect to the server
	conn, err := net.Dial("tcp", "localhost:6969")
	if err != nil {
		t.Fatalf("Could not connect to server: %v", err)
	}
	defer conn.Close()

	t.Log("Successfully connected to the server")
}

func TestEchoFunctionality(t *testing.T) {
	go tcp.StartServer("localhost:6970")
	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("tcp", "localhost:6970")
	if err != nil {
		t.Fatalf("Could not connect to server: %v", err)
	}
	defer conn.Close()

	testMessage := "Hello, Netrunner!"
	_, err = conn.Write([]byte(testMessage))
	if err != nil {
		t.Fatalf("Could not send message to server: %v", err)
	}

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		t.Fatalf("Could not read from server: %v", err)
	}

	response := string(buffer[:n])
	if response != testMessage {
		t.Fatalf("Expected response %q, but got %q", testMessage, response)
	}

	t.Log("Echo functionality working correctly")
}

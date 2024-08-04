package tcp

import (
	"fmt"
	"io"
	"net"
)

func StartServer(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}
	defer listener.Close()

	fmt.Printf("Server listening on %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Failed to accept connection: %v\n", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("New connection from %s\n", conn.RemoteAddr().String())

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Error reading from connection: %v\n", err)
			}
			return
		}

		fmt.Printf("Received: %s", string(buffer[:n]))

		_, err = conn.Write(buffer[:n])
		if err != nil {
			fmt.Printf("Error writing to connection: %v\n", err)
			return
		}
	}
}

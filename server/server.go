package server

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

type Message struct {
	Data         string
	DateSend     time.Time
	DateReceived time.Time
	From         string
}

type Server struct {
	Port     string
	Clients  map[string]net.Conn
	Messages chan Message
}

func NewServer(port string) *Server {
	return &Server{
		Port:     port,
		Clients:  make(map[string]net.Conn),
		Messages: make(chan Message),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", ":"+s.Port)
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Println("TCP server started on port", s.Port)

	// Goroutine to handle incoming messages
	go func() {
		for {
			msg := <-s.Messages
			fmt.Println("Msg: ", msg)
			fmt.Printf("Received from client %s at %s sent at %s data %s\n", msg.From, msg.DateReceived, msg.DateSend, msg.Data)

		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	s.Clients[conn.RemoteAddr().String()] = conn
	fmt.Println("Client connected:", conn.RemoteAddr())

	go func() {
		defer func() {
			delete(s.Clients, conn.RemoteAddr().String())
			fmt.Println("Client disconnected:", conn.RemoteAddr())
			conn.Close()
		}()

		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			receivedMessage := scanner.Text()

			msg := Message{
				Data:         receivedMessage,
				DateSend:     time.Now(),
				DateReceived: time.Now(),
				From:         conn.RemoteAddr().String(),
			}

			s.Messages <- msg

			for addr, clientConn := range s.Clients {
				if addr != conn.RemoteAddr().String() {
					_, err := clientConn.Write([]byte(receivedMessage + "\n"))
					if err != nil {
						fmt.Println("Error writing to client:", err)
					}
				}
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading from client:", err)
		}
	}()
}

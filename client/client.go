package client

import (
	"fmt"
	"net"
)

type TcpClient struct {
	Conn        net.Conn
	SendData    chan []byte
	ReceiveData chan []byte
}

func NewTcpClient(conn net.Conn) *TcpClient {
	return &TcpClient{
		Conn:        conn,
		SendData:    make(chan []byte),
		ReceiveData: make(chan []byte),
	}
}

func (tcp *TcpClient) Close() {

	defer tcp.Conn.Close()

	fmt.Printf("Connection Closed by the client!")
}

func (tcp *TcpClient) Start() {
	go tcp.sendDataRoutine()
	go tcp.receiveDataRoutine()
}

func (tcp *TcpClient) sendDataRoutine() {
	for {
		select {
		case data := <-tcp.SendData:
			_, err := tcp.Conn.Write(data)
			if err != nil {
				fmt.Println("Error sending data:", err)
				return
			}

			fmt.Println("Data Send: ", string(data))
		}
	}
}

func (tcp *TcpClient) receiveDataRoutine() {
	buffer := make([]byte, 1024)
	for {
		n, err := tcp.Conn.Read(buffer)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			return
		}
		tcp.ReceiveData <- buffer[:n]

		fmt.Println("Data receive: ", string(buffer[:n]))
	}
}

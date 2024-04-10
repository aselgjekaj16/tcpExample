package main

import (
	"fmt"
	"net"
	cli "tcpExample/client"
	ser "tcpExample/server"
)

func main() {

	go func() {
		server := ser.NewServer("8080")

		if err := server.Start(); err != nil {
			fmt.Println("Error starting server:", err)
		}
	}()

	clientConn, err := net.Dial("tcp", ":8080")

	defer clientConn.Close()

	if err != nil {
		fmt.Println("Erro conenction from the client: ", err)
		return
	}

	// Start first client
	client := cli.NewTcpClient(clientConn)

	client.Start()

	message := "Hello Server"

	client.SendData <- []byte(message)

	<-make(chan struct{})

}

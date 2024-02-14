package main

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
)

const (
	connHost = "0.0.0.0"
	connPort = "8998"
	connType = "tcp"
)

func main() {
	fmt.Println("Starting " + connType + " server on " + connHost + ":" + connPort)
	l, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return
		}
		fmt.Println("Client connected.")

		fmt.Println("Client " + c.RemoteAddr().String() + " connected.")

		go handleConnection(c)
	}
}

func handleConnection(conn net.Conn) {
	// Read from the client
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}
	fmt.Println("Received from client:", string(buffer[:n]))

	// Send the image to Gotify
	err = sendTextToGotify(string(buffer[:n]))
	if err != nil {
		fmt.Println("Error sending image to Gotify:", err)
		conn.Write([]byte("Error sending image to Gotify."))
		conn.Close()
		return
	}

	// Send a response back to the client
	conn.Write([]byte("Message received."))
	// Close the connection when you're done with it.

	conn.Close()
}

func sendTextToGotify(text string) error {
	// Create the Gotify API URL
	gotifyURL := "https://noti.ratchaphon1412.co/message?token=AJOHE8Y-vlWqt7j"

	// Create the request body
	requestBody := fmt.Sprintf(`{"message": "%s", "priority": 5}`, text)

	// Send the request to Gotify
	resp, err := http.Post(gotifyURL, "application/json", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send text to Gotify: %s", resp.Status)
	}

	return nil
}

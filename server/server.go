package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
)

type Message struct {
	sender  int
	message string
}

func handleError(err error) {
	// TODO: all
	// Deal with an error event.
	if err == io.EOF {

	} else {
		fmt.Println(err)
	}
}

func acceptConns(ln net.Listener, conns chan net.Conn) {
	// TODO: all
	// Continuously accept a network connection from the Listener
	// and add it to the channel for handling connections.

	for {
		conn, _ := ln.Accept()
		conns <- conn
	}
}

func handleClient(client net.Conn, clientid int, msgs chan Message) {
	// TODO: all
	// So long as this connection is alive:
	// Read in new messages as delimited by '\n's
	// Tidy up each message and add it to the messages channel,
	// recording which client it came from.

	reader := bufio.NewReader(client)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			handleError(err)
			break
		}
		msgs <- Message{
			clientid,
			msg,
		}
	}
}

func main() {
	// Read in the network port we should listen on, from the commandline argument.
	// Default to port 8030
	portPtr := flag.String("port", ":8030", "port to listen on")
	flag.Parse()

	//TODO Create a Listener for TCP connections on the port given above.
	ln, _ := net.Listen("tcp", *portPtr)

	//Create a channel for connections
	conns := make(chan net.Conn)
	//Create a channel for messages
	msgs := make(chan Message)
	//Create a mapping of IDs to connections
	clients := make(map[int]net.Conn)

	//Start accepting connections
	go acceptConns(ln, conns)
	for {
		select {
		case conn := <-conns:
			//TODO Deal with a new connection
			// - assign a client ID
			// - add the client to the clients channel
			// - start to asynchronously handle messages from this client

			clientID := len(clients)
			fmt.Fprintln(conn, "Welcome. Your id is ", clientID)
			clients[clientID] = conn
			go handleClient(conn, clientID, msgs)

		case msg := <-msgs:
			//TODO Deal with a new message
			// Send the message to all clients that aren't the sender

			formattedMsg := fmt.Sprintf("%d: %s", msg.sender, msg.message)
			for id, client := range clients {
				if id != msg.sender {
					fmt.Fprintf(client, formattedMsg)
				}
			}
		}
	}
}

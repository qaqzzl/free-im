package main

import (
	"free-im/internal/ws_conn"
	"net/http"
)

func main() {
	// Configure websocket route
	http.HandleFunc("/", ws_conn.HandleConnections)

	// Start listening for incoming chat messages
	//go handleMessages()

	// Start the server on localhost port 8000 and log any errors
	err := http.ListenAndServe(":8989", nil)
	if err != nil {
		panic(err)
	}
}
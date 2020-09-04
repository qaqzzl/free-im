package main

import (
	"free-im/internal/im/tcp_conn"
	"net"
)

func main()  {
	addr, err := net.ResolveTCPAddr("tcp", ":1208")
	if err != nil {
		panic(err)
	}
	server, err := net.ListenTCP("tcp", addr)
	if err != nil {
		print("Fail to start server, %s\n", err)
	}
	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			print("Fail to connect, %s\n", err)
			break
		}
		err = conn.SetKeepAlive(true)
		if err != nil {
			print("Fail to connect, %s\n", err)
			break
		}
		go tcp_conn.ConnSocketHandler(conn)
	}
}
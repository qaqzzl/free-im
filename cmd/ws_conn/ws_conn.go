package main

import (
	"free-im/config"
	"free-im/internal/ws_conn"
)

func main() {
	ws_conn.NewWebSocketServer(config.WSConnConf.WSListenAddr)
}

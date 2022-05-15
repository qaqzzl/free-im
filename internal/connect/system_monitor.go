package connect

import (
	"fmt"
	cmap "github.com/orcaman/concurrent-map"
	"time"
)

// 系统监听
func SystemMonitor() {
	go func() {
		ticker := time.NewTicker(time.Second * 3)
		for {
			<-ticker.C
			fmt.Println("-----------------------------------")
			fmt.Println("连接用户数: ", ConnPool.Count())
			for key, vo := range ConnPool.Items() {
				fmt.Println("--------------")
				fmt.Println("连接用户ID: ", key)
				ConcurrentMap := vo.(cmap.ConcurrentMap)
				for k, v := range ConcurrentMap.Items() {
					fmt.Println("连接设备类型: ", k)
					fmt.Println("连接设备ID: ", v.(Conn).DeviceID)
				}

			}
		}
	}()
}

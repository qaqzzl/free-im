package connect

import (
	"bufio"
	"free-im/pkg/logger"
	"free-im/pkg/protos/pbs"
	"free-im/pkg/service/user"
	"github.com/gorilla/websocket"
	cmap "github.com/orcaman/concurrent-map"
	"net"
	"strconv"
	"sync"
	"time"
)

type ConnType int8

const (
	ConnTypeTCP ConnType = 1 // TCP Conn
	ConnTypeWS  ConnType = 2 // WS Conn
)

// IM连接信息
type Conn struct {
	CoonType      ConnType
	TCP           net.Conn
	TCPReader     *bufio.Reader
	TCPReadTicker *time.Ticker
	WS            *websocket.Conn
	WSMutex       sync.Mutex // 避免重复关闭管道
	Version       int32
	DeviceID      string // 设备id 简写 DID
	UserID        int64  // 用户id 简写 UID
	DeviceType    string // 设备类型, 移动端:mobile , PC端:pc
	ClientType    string // 客户端类型, (android, ios) | (windows, mac, linux)
	IsAuth        bool   // 是否认证(登录)
}

// 链接池
// sync.Map
// [user_id][DeviceType]Context
// var ConnPool cmap.ConcurrentMap
var ConnPool = cmap.New()

// HandleConnect 建立连接
func (conn *Conn) HandleConnect() {
	logger.Logger.Info("connect ")
}

// Write 写入
func (conn *Conn) Write(mp pbs.MsgPackage) error {
	if conn.CoonType == ConnTypeTCP {
		return conn.tcpWrite(mp)
	} else if conn.CoonType == ConnTypeWS {
		return conn.wsWrite(mp)
	}
	return nil
}
func (conn *Conn) tcpWrite(mp pbs.MsgPackage) error {
	if b, err := Protocol.Encode(mp); err != nil {
		return err
	} else {
		_, err := conn.TCP.Write(b)
		if err != nil {
			return err
		}
		return nil
	}
}
func (conn *Conn) wsWrite(mp pbs.MsgPackage) error {
	err := conn.WS.WriteJSON(mp)
	return err
}

// Close 关闭
func (conn *Conn) Close() {
	if conn.IsAuth {
		// 删除连接
		DeleteConn(conn)
		// 用户在线状态
		user.User.SetUserOnline(conn.UserID, false, conn.DeviceType)
	}
	if conn.CoonType == ConnTypeTCP {
		conn.tcpClose()
	} else if conn.CoonType == ConnTypeWS {
		conn.wsClose()
	}
}
func (conn *Conn) tcpClose() {
	conn.TCPReadTicker.Stop()
	conn.TCP.Close()
}
func (conn *Conn) wsClose() {
	conn.WSMutex.Lock()
	defer conn.WSMutex.Unlock()
	conn.WS.Close()
}

// HandlePackage 处理消息包
func (conn *Conn) HandlePackage(mp pbs.MsgPackage) {
	Handler.Handler(conn, mp)
}

// load 获取链接
func LoadConn(UserID string, DeviceID string) (conn *Conn) {
	tmp, ok := ConnPool.Get(UserID)
	if ok && tmp.(cmap.ConcurrentMap).Count() > 0 {
		for _, vo := range tmp.(cmap.ConcurrentMap).Items() {
			conn := vo.(*Conn)
			if conn.DeviceID == DeviceID {
				break
			}
		}
	}
	return conn
}

// load 获取链接通过UID
func LoadConnsByUID(UserID int64) (conns []*Conn) {
	tmp, ok := ConnPool.Get(strconv.Itoa(int(UserID)))
	if ok && tmp.(cmap.ConcurrentMap).Count() > 0 {
		for _, vo := range tmp.(cmap.ConcurrentMap).Items() {
			conn := vo.(*Conn)
			conns = append(conns, conn)
		}
	}
	return conns
}

// store 存储
func StoreConn(conn *Conn) {
	key := strconv.Itoa(int(conn.UserID))
	if tmp, ok := ConnPool.Get(key); ok {
		device_map := tmp.(cmap.ConcurrentMap)
		device_map.Set(conn.DeviceType, conn)
		ConnPool.Set(key, device_map)
	} else {
		device_map := cmap.New()
		device_map.Set(conn.DeviceType, conn)
		ConnPool.Set(key, device_map)
	}
	// 用户在线状态
	user.User.SetUserOnline(conn.UserID, true, conn.DeviceType)
}

// 删除
func DeleteConn(conn *Conn) {
	key := strconv.Itoa(int(conn.UserID))
	if user_map, ok := ConnPool.Get(key); ok {
		user_map.(cmap.ConcurrentMap).Remove(conn.DeviceType)
		ConnPool.Set(key, user_map)
	}
}

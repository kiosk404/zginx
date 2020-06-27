/**
* @Author: kiosk
* @Mail: weijiaxiang007@foxmail.com
* @Date: 2020/6/20
**/
package ziface

import "net"

// define connection interface
type IConnection interface {
	// Create A Connection
	Start()
	// Stop a Connection
	Stop()
	// Accept socket
	GetTCPConnection() *net.TCPConn
	// Get Socket
	GetConnId() uint32
	// Get Remote Addr
	RemoteAdr() net.Addr
	// Send Message
	SendMsg(uint32, []byte) error
	RecvMsg() (IMessage, error)

	//设置链接属性
	SetProperty(key string, value interface{})
	//获取链接属性
	GetProperty(key string) (interface{}, error)
	//移除链接属性
	RemoveProperty(key string)
}

// Define Connection Func
type HandleFunc func(*net.TCPConn, []byte, int) error
/**
* @Author: kiosk
* @Mail: weijiaxiang007@foxmail.com
* @Date: 2020/6/20
**/
package znet

import (
	"errors"
	"io"
	"net"
	"sync"
	"zginx/global"
	"zginx/ziface"
	"zginx/zlog"
)

type Connection struct {
	// Connection Handler
	Conn 			*net.TCPConn
	// Unique connection id
	ConnId 			uint32
	//
	isClosed 		bool
	// Notify Connection to exit (Notify Close Signal From Reader Goroutine  )
	ExitChan 		chan bool
	// Read/Write Msg Chan
	MsgChan			chan []byte
	// Actually, Router is a callback function to do business logic
	Router 			ziface.IMsgHandler
	// Server
	TcpServer 		ziface.IServer
	// Connection Context
	property 		map[string]interface{}
	// 保护链接属性修改的锁
	propertyLock	sync.RWMutex
}


func NewConnection(server ziface.IServer,conn *net.TCPConn,connID uint32, router  ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer: 	server,
		Conn:       conn,
		ConnId:     connID,
		isClosed:   false,
		ExitChan:   make(chan bool),
		MsgChan:    make(chan []byte),
		Router:		router,
		property:   make(map[string]interface{}),
	}

	c.TcpServer.GetConnManager().Add(c)

	return c
}


func (c *Connection) StartReader() {
	var (
		err error
		msg ziface.IMessage
	)

	zlog.ConnLog(c.ConnId).Debugf("[ Reader Goroutine is running ... ]")
	defer zlog.ConnLog(c.ConnId).Debugf("Reader is exit,remote addr is %s ",c.RemoteAdr().String())
	defer c.Stop()

	for {
		msg, err = c.RecvMsg()
		if err != nil {
			break
		}

		// Do Router Func From Bind Connection
		req := Request{iconn: c, msg: msg}

		if global.ZConfig.Server.WorkerPoolSize > 0 {
			// WorkerPool Management
			c.Router.SendMsgToWorkerPool(&req)
		} else {
			// Goroutine Management (Each request has a goroutine)
			go c.Router.DoMessageHandler(&req)
		}

	}
}

func (c *Connection) StartWriter() {
	zlog.ConnLog(c.ConnId).Debugf("[ Writer Goroutine is running ... ]")
	defer zlog.ConnLog(c.ConnId).Debugf("Writer is exit,remote addr is %s ",c.RemoteAdr().String())

	for {
		select {
		case data := <- c.MsgChan:
			// Msg Write To Client
			if _, err := c.Conn.Write(data);err != nil {
				return
			}
		case <- c.ExitChan:
			// Reader is Exited
			return
		}
	}
}

func (c *Connection) Start() {
	defer zlog.ConnLog(c.ConnId).Infof("Conn Start() ... ")

	go c.StartReader()
	go c.StartWriter()

	c.TcpServer.CallOnConnStart(c)

}

func (c *Connection) Stop() {
	defer zlog.ConnLog(c.ConnId).Infof("Conn Stop() ...  ")
	if c.isClosed == true {
		return
	}
	// Notify Writer ,Socket is Closed ..
	c.ExitChan <- true
	c.isClosed = true

	_ = c.Conn.Close()
	c.TcpServer.GetConnManager().Remove(c)
	close(c.ExitChan)
	close(c.MsgChan)

	c.TcpServer.CallOnConnStop(c)

}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnId() uint32 {
	return c.ConnId
}

func (c *Connection) RemoteAdr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgId uint32,data []byte) error {
	var (
		binaryMsg 	[]byte
		err 		error
	)

	if c.isClosed {
		return errors.New("Connection closed when send msg ! ")
	}

	dataP := NewDataPack()

	if binaryMsg, err = dataP.Pack(NewMsgPacket(msgId,data));err != nil {
		return errors.New("Packet Error : " + err.Error())
	}

	c.MsgChan <- binaryMsg

	return nil
}

func (c *Connection) RecvMsg() (ziface.IMessage, error) {
	dataP := NewDataPack()

	headData := make([]byte, dataP.GetHeadLen())

	// read header
	if _, err :=  io.ReadFull(c.GetTCPConnection(), headData); err != nil {
		//zlog.Logger.Debugf("Read Msg Head error : %s", err.Error())
		return nil, errors.New("Read Msg Head error ")
	}

	msg, err := dataP.UnPack(headData)
	if err != nil {
		zlog.ConnLog(c.ConnId).Errorf("Unpack Error : %s", err.Error())
		return nil, errors.New("Unpack Error ")
	}

	var data []byte
	if msg.GetMsgLen() > 0 {
		data = make([]byte, msg.GetMsgLen())
		if _,err := io.ReadFull(c.GetTCPConnection(),data); err != nil {
			return nil, errors.New("Read Msg Error ")
		}
	}
	msg.SetData(data)

	return msg, nil
}

//
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

//获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

//移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
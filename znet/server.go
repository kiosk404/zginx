/**
* @Author: kiosk
* @Mail: weijiaxiang007@foxmail.com
* @Date: 2020/6/20
**/
package znet

import (
	"fmt"
	"net"
	"zginx/global"
	"zginx/ziface"
	"zginx/zlog"
)

// iServer interface
type Server struct {
	// Server Name
	IServerName string
	// Server Version
	Version 	string
	// IP Version
	IPVersion 	string
	// IP
	IP 			string
	// Port
	Port 		int
	// Router Obj Do CallBack Function
	MsgRouter 	ziface.IMsgHandler
	// Set Config
	SetConfig  	global.ZginxConfig
	// TCP Max Message Package Len
	MaxMessagePacket  uint32
	//
	ConnManager ziface.IConnManager
	//
	OnConnStart func(ziface.IConnection)
	//
	OnConnStop  func(ziface.IConnection)
}


func (s *Server) Start() {
	zlog.Logger.Infof("┌───────────────────────────────────────────────────┐")
	zlog.Logger.Infof("%s  [ Start Server Listener at IP %s:%d ]   %s","│",s.IP, s.Port,"│")
	zlog.Logger.Infof("%s  [ %s %s ] Engine  !                         %s","│" ,s.IServerName,global.ZConfig.Server.Version,"│")
	zlog.Logger.Infof("└───────────────────────────────────────────────────┘")

	go func() {
		if global.ZConfig.Server.WorkerPoolSize > 0 {
			// Start a Worker Pool to deal writer event
			s.MsgRouter.StartWorkerPool()
		}

		// Create Socket Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion,fmt.Sprintf("%s:%d",s.IP,s.Port))
		if err != nil {
			zlog.Logger.Panicf("Resolve addr error: %s", err.Error())
		}

		// Listen Socket Addr
		listener, err := net.ListenTCP(s.IPVersion, addr)

		if err != nil {
			zlog.Logger.Panicf("Listen %s:%d error",addr.IP,addr.Port)
		}

		zlog.Logger.Infof("Start Zginx server success, %s succ Listening ...", s.IServerName)

		// Init ConnId
		var cid uint32
		cid = 0

		// Accept Client TCP Requests & Loop
		for {
			conn, err := listener.AcceptTCP()
			_ = conn.SetKeepAlive(false)
			_ = conn.SetNoDelay(true)
			_ = conn.SetLinger(2)
			if err != nil {
				zlog.ConnLog(cid).Errorf("%s shutdown \n",conn.RemoteAddr().String())
				continue
			}

			if s.ConnManager.GetLen() >= global.ZConfig.TCP.MaxConn {
				zlog.ConnLog(cid).Warnf("Too many connections ... ")
				_ = conn.Close()
				continue
			}
			dealConn := NewConnection(s, conn, cid, s.MsgRouter)

			cid ++

			// Process business logic
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
 	// todo: stop server and resources recovery
	s.ConnManager.CleanAllConn()
}


func (s *Server) Serve() {
	s.Start()
	// todo: do something
	// Block
	select {}
}


func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgRouter.AddRouter(msgId, router)
	zlog.Logger.Infof("Add %d Router Succ !!", msgId)
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnManager
}

func (s *Server) SetOnConnStart(hookFunc func(ziface.IConnection)) {
	s.OnConnStart = hookFunc
}
// Hook Func When Conn Destroy
func (s *Server) SetOnConnStop(hookFunc func(ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// Call Func OnConnStart Hook
func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> CallOnConnStart....")
		s.OnConnStart(conn)
	}
}

// Call Func OnConnStop Hook
func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}


/*
  Init Server From Config
*/
func NewServer(config global.ZginxConfig) ziface.IServer {
	config.SetConfig()
	fmt.Println(global.ZConfig.TCP)
	return &Server {
		IServerName: global.ZConfig.Server.Name,
		Version: global.ZConfig.Server.Version,
		IPVersion: global.ZConfig.TCP.Version,
		IP: global.ZConfig.TCP.Host,
		Port: global.ZConfig.TCP.Port,
		MaxMessagePacket: global.ZConfig.TCP.MaxPacketSize,
		MsgRouter: NewMsgHandler(),
		ConnManager: NewConnManager(),
	}
}

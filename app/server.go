/**
* @Author: kiosk
* @Mail: weijiaxiang007@foxmail.com
* @Date: 2020/6/20
**/
package main

import (
	"fmt"
	"zginx/global"
	"zginx/ziface"
	"zginx/zlog"
	"zginx/znet"
)

// Powered by Zginx


type PingRouter struct {
	znet.BaseRouter
}


//Test MainRouter
func (c *PingRouter) MainHandle(request ziface.IRequest) {
	conn := request.GetConnection()
	recv := request.GetData()

	zlog.ConnLog(conn.GetConnId()).Infof("recv data from [remoteAddr: %s RouterId: %d  Data: %s] ",
		conn.RemoteAdr().String(), request.GetMsgID(), recv)

	pongStr := fmt.Sprintf("pong... : %s ", recv)

	err := conn.SendMsg(request.GetMsgID(), []byte(pongStr))

	if err != nil {
		conn.Stop()
	}
}


func main() {
	s := znet.NewServer(global.ZginxConfig{ConfFilePath:"../conf/server.yml"})

	s.AddRouter(11,&PingRouter{})

	s.SetOnConnStart(func(connection ziface.IConnection) {
		fmt.Printf("➜  connection %d jion in \n",connection.GetConnId())
		_ = connection.SendMsg(0, []byte("welcome to zginx ... "))
		connection.SetProperty("server","zginx context")
	})

	s.Serve()
}
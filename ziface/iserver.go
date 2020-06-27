/**
* @Author: kiosk
* @Mail: weijiaxiang007@foxmail.com
* @Date: 2020/6/20
**/
package ziface

type IServer interface {
	// Start Server
	Start()
	// Stop Server
	Stop()
	// Run Server
	Serve()
	// Router Model
	AddRouter(msgId uint32, router IRouter)
	// Get Conn Manager
	GetConnManager() IConnManager
	// Hook Func When Conn Create
	SetOnConnStart(func(IConnection))
	// Hook Func When Conn Destroy
	SetOnConnStop(func(IConnection))
	//
	CallOnConnStart(IConnection)
	//
	CallOnConnStop(IConnection)
}
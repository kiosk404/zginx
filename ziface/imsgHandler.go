/**
* @Author: kiosk
* @Mail: weijiaxiang007@foxmail.com
* @Date: 2020/6/26
**/
package ziface

/*
   Message Control Handler
*/

type IMsgHandler interface {
	// Do Router Message
	DoMessageHandler(request IRequest)
	// Add
	AddRouter(msgId uint32,router IRouter)
	// Create A Worker Pool
	StartWorkerPool()
	//
	SendMsgToWorkerPool(request IRequest)
}
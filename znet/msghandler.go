/**
* @Author: kiosk
* @Mail: weijiaxiang007@foxmail.com
* @Date: 2020/6/26
**/
package znet

import (
	"zginx/global"
	"zginx/ziface"
	"zginx/zlog"
)

type MsgHandle struct {
	// Router Api
	Apis 			map[uint32] ziface.IRouter
	// Worker Message Queue size
 	WorkerPollSize 	uint32
	// Worker Message Queue
	WorkerPool		[]chan ziface.IRequest
}

// create
func NewMsgHandler() *MsgHandle {
	if global.ZConfig.Server.WorkerPoolSize > 125 {
		zlog.Logger.Warnf("Too Many Worker Pool :%d , It's be DANGER", global.ZConfig.Server.WorkerPoolSize)
	}

	return &MsgHandle{
		Apis:make(map[uint32] ziface.IRouter),
		WorkerPollSize: global.ZConfig.Server.WorkerPoolSize,
		WorkerPool: make([]chan ziface.IRequest, global.ZConfig.Server.WorkerPoolSize),
	}
}

func (mh *MsgHandle) DoMessageHandler(request ziface.IRequest) {
	// 1. get msgId from request.message
	handler,ok := mh.Apis[request.GetMsgID()]
	if !ok {
		zlog.ConnLog(request.GetConnection().GetConnId()).
			Warnf("api msgId = %d, is not found! Need Register ", request.GetMsgID())
		_ = request.GetConnection().SendMsg(0, []byte("404, No Msg Router"))
		return
	}
	// 2.
	handler.PreHandle(request)
	handler.MainHandle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandle) AddRouter(msgId uint32,router ziface.IRouter) {
	// 1. Check Router Bind Already
	if _, ok := mh.Apis[msgId]; ok {
		zlog.Logger.Fatalf("repeat api, msgId = %d", msgId)
	}
	// 2. Add Router
	mh.Apis[msgId] = router
	zlog.Logger.Infof("Add api MsgID = %d succ !! ", msgId)
}

// Start Worker Pool
func (mh *MsgHandle) StartWorkerPool() {
	for i := 0 ; i < int(mh.WorkerPollSize); i++ {
		// Request memory space
		mh.WorkerPool[i] = make(chan ziface.IRequest, global.ZConfig.Server.MaxWorkerTaskLen)
		go mh.StartOneWorkerFlow(i, mh.WorkerPool[i])
	}
	zlog.Logger.Infof("%d workers are started ... ", global.ZConfig.Server.WorkerPoolSize)
}

// Create Workflow
func (mh *MsgHandle) StartOneWorkerFlow(workerId int, workerPool chan ziface.IRequest) {
	for {
		select {
		case request := <- workerPool:
			mh.DoMessageHandler(request)
		}
	}
}


func (mh *MsgHandle) SendMsgToWorkerPool(request ziface.IRequest)  {
	//1. Select a Worker
	// todo : create a balance func to select worker
	workerID := request.GetConnection().GetConnId() % mh.WorkerPollSize

	zlog.ConnLog(request.GetConnection().GetConnId()).
		Debugf("Add ConnId = %d request to WorkerID = %d",request.GetConnection().GetConnId(), workerID)

	//2. Send Msg to Worker
	mh.WorkerPool[workerID] <- request
}



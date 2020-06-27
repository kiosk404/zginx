/**
* @Author: kiosk
* @Mail: weijiaxiang007@foxmail.com
* @Date: 2020/6/24
**/
package znet

import (
	"zginx/ziface"
)

type Request struct {
	iconn ziface.IConnection
	msg   ziface.IMessage
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.iconn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}


func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
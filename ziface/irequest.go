/**
* @Author: kiosk
* @Mail: weijiaxiang007@foxmail.com
* @Date: 2020/6/24
**/
package ziface

type IRequest interface {
	// Get Connection Interface
	GetConnection() IConnection

	// Get Request Message Data
	GetData() []byte

	// Get Message Id
	GetMsgID() uint32
}



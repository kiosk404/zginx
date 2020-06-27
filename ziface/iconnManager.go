/**
* @Author: kiosk
* @Mail: weijiaxiang007@foxmail.com
* @Date: 2020/6/26
**/
package ziface

type IConnManager interface {
	// Add a Connection
	Add(conn IConnection)
	// Remove a Connection
	Remove(conn IConnection)
	// Get Connection By Id
	Get(connID uint32) (IConnection , error)
	// Get Connection Count
	GetLen() int
	// Clean ALL Connection
	CleanAllConn()
}


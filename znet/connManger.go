/**
* @Author: kiosk
* @Mail: weijiaxiang007@foxmail.com
* @Date: 2020/6/26
**/
package znet

import (
	"errors"
	"sync"
	"zginx/ziface"
	"zginx/zlog"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock 	sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (cm *ConnManager) Add(conn ziface.IConnection) {
	// Protect map shared resource ！！
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	cm.connections[conn.GetConnId()] = conn
	zlog.ConnLog(conn.GetConnId()).Infof("Connection add to ConnManager successfully ! ")
}

func (cm *ConnManager) Remove(conn ziface.IConnection) {
	// Protect map shared resource ！！
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	delete(cm.connections, conn.GetConnId())
}
func (cm *ConnManager) Get(connID uint32) (ziface.IConnection , error) {
	// Protect map shared resource ！！
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	if conn, ok := cm.connections[connID];ok {
		return conn,nil
	}
	return nil, errors.New("connection not found ... ")
}

func (cm *ConnManager) GetLen() int {
	return len(cm.connections)
}

// Clean ALL Connection
func (cm *ConnManager) CleanAllConn() {
	// Protect map shared resource ！！
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	for connId,conn := range cm.connections {
		// Stop
		conn.Stop()
		// Delete
		delete(cm.connections, connId)
	}
	zlog.Logger.Infof("Clean All Connection ... ")
}
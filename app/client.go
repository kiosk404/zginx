/**
* @Author: kiosk
* @Mail: weijiaxiang007@foxmail.com
* @Date: 2020/6/20
**/
package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
	"zginx/ziface"
	"zginx/zlog"
	"zginx/znet"
)

type ClientConnection struct {
	Conn 			*net.TCPConn
	MsgId			uint32
	ExitChan 		chan bool
}

func (cc *ClientConnection) RecvMsg() (ziface.IMessage, error) {
	dataP := znet.NewDataPack()
	headData := make([]byte, dataP.GetHeadLen())

	// read header
	if _, err :=  io.ReadFull(cc.Conn, headData); err != nil {
		zlog.Logger.Debugf("Read Msg Head error : %s", err.Error())
		return nil, errors.New("Read Msg Head error ")
	}

	dataBuff := bytes.NewReader(headData)

	msg := &znet.Message{}

	// Get dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// Get msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	var data []byte
	if msg.GetMsgLen() > 0 {
		data = make([]byte, msg.GetMsgLen())
		if _,err := io.ReadFull(cc.Conn,data); err != nil {
			zlog.Logger.Debugf("Read Msg Error : %s", err.Error())
			return nil,errors.New("Read Msg Error ")
		}
	}
	msg.SetData(data)

	return msg, nil
}

func (cc *ClientConnection) SendMsg(msgId uint32,data []byte) error {
	var (
		binaryMsg 	[]byte
		err 		error
	)

	dataP := znet.NewDataPack()

	if binaryMsg, err = dataP.Pack(znet.NewMsgPacket(msgId,data));err != nil {
		return errors.New("Packet Error : " + err.Error())
	}

	if _, err = cc.Conn.Write(binaryMsg);err != nil {
		return errors.New("conn Write error ")
	}

	return nil
}

func (cc *ClientConnection) Stop() {
	_ = cc.Conn.Close()
	close(cc.ExitChan)
}

func NewClientConnection(conn *net.TCPConn, msgId uint32) *ClientConnection {
	return &ClientConnection{
		Conn:		conn,
		MsgId: 		msgId,
		ExitChan: 	make(chan bool),
	}
}


func DoReader(c *ClientConnection) {
	for {
		msg, err := c.RecvMsg()
		if err != nil {
			fmt.Println("close read socket")
			c.ExitChan <- true
			break
		}
		fmt.Printf("%s \n",string(msg.GetData()))
	}
}

func DoWriter(c *ClientConnection) {
	for {
		var inputBuf []byte
		_, _ = fmt.Scanln(&inputBuf)

		if string(inputBuf) == "exit" {
			c.ExitChan <- true
			break
		}
		if err := c.SendMsg(c.MsgId,inputBuf);err != nil {
			fmt.Printf("write conn err :%s ", err)
			break
		}
	}
}

func main() {
	time.Sleep(1 * time.Second)
	// 1 直接连接远程服务器，得到一个conn连接
	// laddr, _ := net.ResolveTCPAddr("tcp4","127.0.0.1:44444")
	raddr, _ := net.ResolveTCPAddr("tcp4","127.0.0.1:5000")

	conn, err := net.DialTCP("tcp4",nil, raddr)
	if err != nil {
		panic("error connect")
	}

	_ = conn.SetKeepAlive(false)
	_ = conn.SetLinger(2)

	fmt.Println("Start A Client, Connection Successfully ...")

	c := NewClientConnection(conn,11)
	// 2
	go DoReader(c)
	go DoWriter(c)
	select {
	case <- c.ExitChan:
		c.Stop()
		os.Exit(0)
	}
}
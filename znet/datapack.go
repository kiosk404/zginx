/**
* @Author: kiosk
* @Mail: weijiaxiang007@foxmail.com
* @Date: 2020/6/25
**/
package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"zginx/global"
	"zginx/ziface"
)

type DataPack struct {}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (dp* DataPack) GetHeadLen() uint32 {
	/*
	DataLen uint32 4 Bytes + ID uint32 4 Bytes
	*/
	return 8
}

func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})

	// Set DataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen());err != nil {
		return nil, err
	}
	// Set DataId
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId());err != nil {
		return nil, err
	}
	// Set DataBuff
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData());err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (dp *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	dataBuff := bytes.NewReader(binaryData)

	msg := &Message{}

	// Get dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// Get msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	if global.ZConfig.TCP.MaxPacketSize > 0 && msg.GetMsgLen() > global.ZConfig.TCP.MaxPacketSize {
		return nil, errors.New("Too Large Message !!! ")
	}

	return msg, nil
}


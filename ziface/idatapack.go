/**
* @Author: kiosk
* @Mail: weijiaxiang007@foxmail.com
* @Date: 2020/6/25
**/
package ziface

type IDataPack interface {
	GetHeadLen() uint32

	Pack(msg IMessage)([]byte, error)

	UnPack([]byte)(IMessage, error)
}
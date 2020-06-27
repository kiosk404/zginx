/**
* @Author: kiosk
* @Mail: weijiaxiang007@foxmail.com
* @Date: 2020/6/24
**/
package ziface

type IRouter interface {
	PreHandle(request  IRequest)
	MainHandle(request IRequest)
	PostHandle(request IRequest)
}





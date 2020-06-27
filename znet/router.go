/**
* @Author: kiosk
* @Mail: weijiaxiang007@foxmail.com
* @Date: 2020/6/24
**/
package znet

import "zginx/ziface"

type BaseRouter struct {}

// Base Router is nil,
// Business model will be using all inherited methods from BaseRouter Model.

// Pre Hook
func (r *BaseRouter) PreHandle(request  ziface.IRequest) {}

// Do Hook
func (r *BaseRouter) MainHandle(request ziface.IRequest) {}

// Post Hook
func (r *BaseRouter) PostHandle(request ziface.IRequest) {}


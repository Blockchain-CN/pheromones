// Copyright 2018 Blockchain-CN . All rights reserved.
// https://github.com/Blockchain-CN

package pheromones

import "net"

// MsgPto 协议数据格式
type MsgPto struct {
	Name      string `json:"name"`
	Operation string `json:"operation"`
	// 子协议json
	Data      []byte `json:"data"`
}

// Protocal 路由数据解析协议
type Protocal interface {
	// 获取协议的链接类型
	GetConnType() ConnType

	// 解析请求通信内容,并返回数据,双工协议
	Handle(c net.Conn, msg []byte) ([]byte, error)
}

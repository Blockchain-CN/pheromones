// Copyright 2018 Blockchain-CN . All rights reserved.
// https://github.com/Blockchain-CN

package pheromones

// ConnType 连接类型
type ConnType int

const (
	// PersistentConnection 长连接方式
	PersistentConnection ConnType = iota
	// ShortConnection 短连接方式
	ShortConnection
)

// Router 路由接口
// 实现了长连接／短连接两种通信方式
type Router interface {
	// 添加路由：短链接传的是地址；长链接传的是net.Conn
	AddRoute(name string, addr interface{}) error

	// 删除路由
	Delete(name string) error

	// 获取连接类型
	GetConnType() ConnType

	// 广播发送信息
	DispatchAll(msg []byte) map[string][]byte

	// 单点发送信息
	Dispatch(name string, msg []byte) ([]byte, error)

	// 获取全部peer
	FetchPeers() map[string]interface{}
}

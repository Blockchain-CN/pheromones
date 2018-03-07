// Copyright 2018 Blockchain-CN . All rights reserved.
// https://github.com/Blockchain-CN

package pheromones

const (
	// ErrLocalSocketTimeout 读取数据超时
	ErrLocalSocketTimeout = 1001

	// ErrRemoteSocketEmpty 远程socket地址为空
	ErrRemoteSocketEmpty   = 2001
	// ErrRemoteSocketExist 远程socket地址已存在
	ErrRemoteSocketExist   = 2002
	// ErrRemoteSocketMisType 远程socket链接类型错误
	ErrRemoteSocketMisType = 2003
	// ErrRemoteSocketConnect 远程socket连接失败
	ErrRemoteSocketConnect = 2004

	// ErrUnKnownProtocal 未知的数据协议格式
	ErrUnKnownProtocal            = 3003
	// ErrMismatchProtocalReq 请求数据格式不匹配
	ErrMismatchProtocalReq        = 3101
	// ErrMismatchProtocalConnectReq 连接请求数据格式不匹配
	ErrMismatchProtocalConnectReq = 3102
	// ErrMismatchProtocalResp 返回数据格式不匹配
	ErrMismatchProtocalResp       = 3202

	// ErrUnknuowPeer peer未添加路由表
	ErrUnknuowPeer = 4001
)

// Error ...
type Error int

// Error ...
func (err Error) Error() string {
	return errMap[err]
}

var errMap = map[Error]string{
	ErrLocalSocketTimeout: "链接超时",

	ErrRemoteSocketEmpty:   "链接为空",
	ErrRemoteSocketExist:   "链接已存在",
	ErrRemoteSocketMisType: "链接类型错误",
	ErrRemoteSocketConnect: "远程链接失败",

	ErrUnKnownProtocal:            "未知的协议类型",
	ErrMismatchProtocalReq:        "请求协议数据类型不匹配",
	ErrMismatchProtocalConnectReq: "连接请求不合法",
	ErrMismatchProtocalResp:       "返回协议数据类型不匹配",

	ErrUnknuowPeer: "未知peer",
}

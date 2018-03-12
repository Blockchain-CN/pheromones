// Copyright 2018 Blockchain-CN . All rights reserved.
// https://github.com/Blockchain-CN

package pheromones

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

// 短连接对象
type EndPointS struct {
	Addr string
}

// SRouter 短连接路由
type SRouter struct {
	sync.RWMutex
	sync.WaitGroup
	to time.Duration
	// 短链接池
	Pool map[string]EndPointS
}

// NewSRouter 建立短连接路由
func NewSRouter(to time.Duration) *SRouter {
	var r SRouter
	r.to = to
	r.Pool = make(map[string]EndPointS, 0)
	return &r
}

// AddRoute 添加路由时，已添加或者地址为空是都返回有错误，防止收到请求和主动连接重复建立
// 如果名字相同地址不同，则将原来的地址删除
func (r *SRouter) AddRoute(name string, addr interface{}) error {
	if _, ok := addr.(string); !ok {
		return Error(ErrRemoteSocketMisType)
	}
	if addr.(string) == "" {
		return Error(ErrRemoteSocketEmpty)
	}
	r.RLock()
	if a, ok := r.Pool[name]; ok {
		if a.Addr == addr.(string) {
			return Error(ErrRemoteSocketExist)
		}
	}
	r.RUnlock()
	fmt.Printf("添加路由, peername=@%s@||peeraddress=%s\n", name, addr.(string))
	r.Lock()
	defer r.Unlock()
	r.Pool[name] = EndPointS{addr.(string)}
	return nil
}

// Delete 删除peer
func (r *SRouter) Delete(s string) error {
	fmt.Printf("删除节点：%v\n", s)
	r.Lock()
	defer r.Unlock()
	delete(r.Pool, s)
	return nil
}

// GetConnType 获取连接类型
func (r *SRouter) GetConnType() ConnType {
	return ShortConnection
}

// DispatchAll 广播消息
func (r *SRouter) DispatchAll(msg []byte) map[string][]byte {
	var l sync.Mutex
	peers := r.FetchPeers()
	resps := make(map[string][]byte)
	for k, v := range peers {
		r.Add(1)
		go func(name, addr string) {
			fmt.Printf("dispatchall||name=%s||addr=%s\n", name, addr)
			defer r.Done()
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("panic: %v", err)
				}
			}()
			resp, err := r.dispatch(addr, msg)
			if err != nil {
				return
			}
			fmt.Printf("dispatchall||msg=%s\n", string(resp))
			l.Lock()
			resps[name] = resp
			l.Unlock()
		}(k, v.(EndPointS).Addr)
	}
	r.Wait()
	return resps
}

// clone
func (r *SRouter) FetchPeers() map[string]interface{} {
	p2 := make(map[string]interface{})
	r.RLock()
	defer r.RUnlock()
	for k, v := range r.Pool {
		p2[k] = v
	}
	return p2
}

// Dispatch 单点传输
func (r *SRouter) Dispatch(name string, msg []byte) ([]byte, error) {
	peer, err := r.getPeer(name)
	if err != nil {
		return nil, err
	}
	return r.dispatch(peer.Addr, msg)
}

// clone
func (r *SRouter) getPeer(name string) (*EndPointS, error) {
	p2 := &EndPointS{}
	r.RLock()
	defer r.RUnlock()
	if _, ok := r.Pool[name]; !ok {
		return p2, Error(ErrUnknuowPeer)
	}
	p2.Addr = r.Pool[name].Addr
	return p2, nil
}

func (r *SRouter) dispatch(addr string, msg []byte) ([]byte, error) {
	var resp []byte
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	for i := 0; i < 3; i++ {
		_, err = c.Write(msg)
		if err != nil {
			continue
		}
		fmt.Printf("发送请求, localhost=%s||remotehost=%s||msg=%s\n", c.LocalAddr(), c.RemoteAddr(), string(msg))
		resp, err = r.read(c, r.to)
		if err != nil {
			continue
		}
		fmt.Printf("收到回复, localhost=%s||remotehost=%s||msg=%s\n", c.LocalAddr(), c.RemoteAddr(), string(resp))
		return resp, err
	}
	return resp, err
}

func (r *SRouter) read(io io.Reader, to time.Duration) ([]byte, error) {
	buf := make([]byte, defultByte)
	messnager := make(chan int)
	go func() {
		n, _ := io.Read(buf[:])
		messnager <- n
		close(messnager)
	}()
	select {
	case n := <-messnager:
		return buf[:n], nil
	case <-time.After(to):
		return nil, Error(ErrLocalSocketTimeout)
	}
	return buf, nil
}

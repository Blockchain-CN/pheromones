// Copyright 2018 Blockchain-CN . All rights reserved.
// https://github.com/Blockchain-CN

package pheromones

import (
	"fmt"
	"io"
	"net"
	"runtime"
	"strings"
	"time"
)

const defultByte = 10240

// Server p2p监听连接server
type Server struct {
	// 如果支持双链接，需要再建立一个proto对象
	proto Protocal
	to    time.Duration
}

// NewServer ...
func NewServer(p Protocal, to time.Duration) *Server {
	return &Server{p, to}
}

// ListenAndServe 监听peer的链接请求
func (s *Server) ListenAndServe(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		println(err.Error())
		return err
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				runtime.Gosched()
				continue
			}
			// theres no direct way to detect this error because it is not exposed
			if !strings.Contains(err.Error(), "use of closed network connection") {
			}
			break
		}
		go s.handler(c)
	}
	return nil
}

func (s *Server) handler(c net.Conn) {
	defer func() {
		if s.proto.GetConnType() == ShortConnection{
			c.Close()
		}
	}()
	msg, err := s.read(c, s.to)
	fmt.Printf("收到请求, localhost=%s||remotehost=%s||msg=%s\n", c.LocalAddr(), c.RemoteAddr(), string(msg))
	if err != nil {
		return
	}
	resp, err := s.proto.Handle(c, msg)
	if err != nil || resp == nil {
		resp = nil
	}
	c.SetWriteDeadline(time.Now().Add(s.to))
	for i := 0; i < 3; i++ {
		_, err = c.Write(resp)
		if err != nil {
			continue
		}
		fmt.Printf("发送回复, localhost=%s||remotehost=%s||msg=%s\n", c.LocalAddr(), c.RemoteAddr(), string(resp))
		return
	}
}

func (s *Server) read(r io.Reader, to time.Duration) ([]byte, error) {
	buf := make([]byte, defultByte)
	messnager := make(chan int)
	go func() {
		n, _ := r.Read(buf[:])
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

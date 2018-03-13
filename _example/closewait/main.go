package main

import (
	"net"
	"io"
	"fmt"
	"time"
)

func main() {
	go Serve()
	time.Sleep(time.Second)
	c, err := net.Dial("tcp", "127.0.0.1:12345")
	if err != nil {
	}
	n, err := c.Write([]byte("11111"))
	fmt.Printf("client say: n=%d||err=%v\n", n, err)
	time.Sleep(time.Second *3)
	c.Close()
	n, err = c.Write([]byte("22222"))
	fmt.Printf("client say: n=%d||err=%v\n", n, err)
	for {
		time.Sleep(time.Second)
	}

}

func Serve() {
	ln, err := net.Listen("tcp", "127.0.0.1:12345")
	if err != nil {
		println("listen tcp fail.")
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			println("accept connect fail.")
		}
		println("server say: ", c.LocalAddr().String(), "connect with ", c.RemoteAddr().String())
		go IOLoop(c)
	}
}

func IOLoop(c net.Conn) {
	for {
		msg, err := read(c)
		if err != nil {
			c.Close()
			fmt.Printf("长连接收到信息,连接关闭, localhost=%s||remotehost=%s||msg=%s||err=%v\n", c.LocalAddr(), c.RemoteAddr(), string(msg), err)
			return
		}
		fmt.Printf("长连接收到信息,继续读取, localhost=%s||remotehost=%s||msg=%s\n", c.LocalAddr(), c.RemoteAddr(), string(msg))
	}
}

func read(r io.Reader) ([]byte, error) {
	buf := make([]byte, 10000)
	n, err := r.Read(buf)
	if err != nil {
		return nil, err
	}
	// read读出来的是[]byte("abcdefg"+0x00)，带一个结束符，需要去掉
	return buf[:n], nil
}


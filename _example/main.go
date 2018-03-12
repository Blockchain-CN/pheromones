package main

import (
	"encoding/json"
	"fmt"
	"time"

	p2p "github.com/Blockchain-CN/pheromones"
	pto "github.com/Blockchain-CN/pheromones/_example/protocal"
)

var (
	hello1 = p2p.MsgPto{
		Name:      "luda",
		Operation: pto.ConnectReq,
	}
	hellomsg1 = pto.MsgGreetingReq{
		Addr:    "127.0.0.1:12345",
		Account: 11900,
	}

	hello2 = p2p.MsgPto{
		Name:      "yoghurt",
		Operation: pto.ConnectReq,
	}
	hellomsg2 = pto.MsgGreetingReq{
		Addr:    "127.0.0.1:12346",
		Account: 11900,
	}
)

const (
	timeout = time.Millisecond * 100
)

func main() {
	mainShort()
	//mainPeriment()
}

func mainShort() {
	r1 := p2p.NewSRouter(timeout)
	p1 := pto.NewProtocal("luda", r1, timeout)
	s1 := p2p.NewServer(p1, timeout)
	println("h1 监听 12345")
	go s1.ListenAndServe("127.0.0.1:12345")

	r2 := p2p.NewSRouter(timeout)
	p2 := pto.NewProtocal("yoghurt", r2, timeout)
	s2 := p2p.NewServer(p2, timeout)
	println("h2 监听 12345")
	go s2.ListenAndServe("127.0.0.1:12346")

	r3 := p2p.NewSRouter(timeout)
	p3 := pto.NewProtocal("diudiu", r3, timeout)
	s3 := p2p.NewServer(p3, timeout)
	println("h3 监听 12345")
	go s3.ListenAndServe("127.0.0.1:12347")

	time.Sleep(time.Second)
	p1.Add("yoghurt", "127.0.0.1:12346")
	j, _ := json.Marshal(hellomsg1)
	hello1.Data = j
	msg, _ := json.Marshal(hello1)
	for msg != nil {
		b, err := p1.Dispatch("yoghurt", msg)
		if err != nil {
			println("操作失败", err.Error())
			break
		}
		msg = nil
		msg, err = p1.Handle(nil, b)
		fmt.Println(string(msg), err)
	}
	fmt.Println("test1 done")

	j, _ = json.Marshal(hellomsg2)
	hello2.Data = j
	msg, _ = json.Marshal(hello2)
	for msg != nil {
		b, err := p2.Dispatch("luda", msg)
		if err != nil {
			println("操作失败", err.Error())
			break
		}
		msg = nil
		msg, err = p2.Handle(nil, b)
		fmt.Println(string(msg), err)
	}
	fmt.Println("test2 done")

	p2.Add("diudiu", "127.0.0.1:12347")
	j, _ = json.Marshal(hellomsg2)
	hello2.Data = j
	msg, _ = json.Marshal(hello2)
	b := p2.DispatchAll(msg)
	fmt.Println(b)
	fmt.Println("test3 done")

	p2.Delete("diudiu1")
	j, _ = json.Marshal(hellomsg2)
	hello2.Data = j
	msg, _ = json.Marshal(hello2)
	b = p2.DispatchAll(msg)
	fmt.Println(b)

	for {
		time.Sleep(time.Second)
	}
}

func mainPeriment() {
	r1 := p2p.NewPRouter(timeout)
	p1 := pto.NewProtocal("luda", r1, timeout)
	s1 := p2p.NewServer(p1, timeout)
	println("h1 监听 12345")
	go s1.ListenAndServe("127.0.0.1:12345")

	r2 := p2p.NewPRouter(timeout)
	p2 := pto.NewProtocal("yoghurt", r2, timeout)
	s2 := p2p.NewServer(p2, timeout)
	println("h2 监听 12345")
	go s2.ListenAndServe("127.0.0.1:12346")

	r3 := p2p.NewPRouter(timeout)
	p3 := pto.NewProtocal("diudiu", r3, timeout)
	s3 := p2p.NewServer(p3, timeout)
	println("h3 监听 12345")
	go s3.ListenAndServe("127.0.0.1:12347")

	time.Sleep(time.Second)
	fmt.Println("test1 start__________________________________________")
	p1.Add("yoghurt", "127.0.0.1:12346")
	j, _ := json.Marshal(hellomsg1)
	hello1.Data = j
	msg, _ := json.Marshal(hello1)
	_, err := p1.Dispatch("yoghurt", msg)
	if err != nil {
		println("操作失败", err.Error())
	}
	fmt.Println("test1 done*******************************************")

	time.Sleep(time.Second)
	fmt.Println("test2 start__________________________________________")
	j, _ = json.Marshal(hellomsg2)
	hello2.Data = j
	msg, _ = json.Marshal(hello2)
	_, err = p2.Dispatch("luda", msg)
	if err != nil {
		println("操作失败", err.Error())
	}
	fmt.Println("test2 done*******************************************")

	time.Sleep(time.Second * 1)
	fmt.Println("test3 start__________________________________________")
	p2.Add("diudiu", "127.0.0.1:12347")
	j, _ = json.Marshal(hellomsg2)
	hello2.Data = j
	msg, _ = json.Marshal(hello2)
	p2.DispatchAll(msg)
	fmt.Println("test3 done*******************************************")

	time.Sleep(time.Second * 1)
	fmt.Println("test4 start__________________________________________")
	p2.Delete("diudiu")
	j, _ = json.Marshal(hellomsg2)
	hello2.Data = j
	msg, _ = json.Marshal(hello2)
	p2.DispatchAll(msg)
	fmt.Println("test4 done_______________________________________")

	for {
		time.Sleep(time.Second)
	}
}

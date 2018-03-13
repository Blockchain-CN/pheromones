# pheromones
Peer to Peer network.

## 网络模型
  在P2P网络环境中，彼此连接的多台计算机之间都处于对等的地位，各台计算机有相同的功能，无主从之分，一台计算机既可作为服务器，设定共享资源供网络中其他计算机所使用，又可以作为工作站，整个网络一般来说不依赖专用的集中服务器，也没有专用的工作站。网络中的每一台计算机既能充当网络服务的请求者，又对其它计算机的请求做出响应，提供资源、服务和内容。通常这些资源和服务包括：信息的共享和交换、计算资源（如CPU计算能力共享）、存储共享（如缓存和磁盘空间的使用）、网络共享、打印机共享等。
  p2p网络的实现要基于传输层协议(TCP/UDP)，而使用TCP协议时又分为短连接和长连接：
![image](https://github.com/Blockchain-CN/pheromones/raw/master/readme_image/short.png)
  
  短连接实现中，每次由router层重新创建连接，进行通信，并等待回复。
![image](https://github.com/Blockchain-CN/pheromones/raw/master/readme_image/perminent.png)
  
  长连接实现中，应该由protocal层、更高层创建连接，并开启线程保持监听(因为要对监听结果进行处理，也就是调用protocal层实现的解析协议。因此监听线程必须由protocal层或更高层进行维护)
  由于在广播数据的时候，需要进行传递式的广播，会形成广播风暴
![image](https://github.com/Blockchain-CN/pheromones/raw/master/readme_image/broadcast.png)
  
  因此需要在protocal层实现的传输协议中保证，传输的内容幂等(同一个消息请求n次的结果一致)且转发的信息不会再二次转发，这部分具体在协议层的实现中来规定。

## 使用
### 短连接
#### 创建p2p节点

``` go
r1 := p2p.NewSRouter(timeout) // 短连接路由
p1 := pto.NewProtocal("luda", r1, timeout)
s1 := p2p.NewServer(p1, timeout)
println("h1 监听 12345")
go s1.ListenAndServe("127.0.0.1:12345")

r2 := p2p.NewSRouter(timeout) // 短连接路由
p2 := pto.NewProtocal("yoghurt", r2, timeout)
s2 := p2p.NewServer(p2, timeout)
println("h2 监听 12345")
go s2.ListenAndServe("127.0.0.1:12346")
```
#### 添加路由

``` go
p1.Add("yoghurt", "127.0.0.1:12346")
```

#### 发送数据
由于为了完成协议状态机，因此需要循环对返回结果进行协议解析，直到返回结果为空

``` go
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
```

### 长连接

#### 创建p2p节点

``` go
r1 := p2p.NewPRouter(timeout)  // 长连接路由
p1 := pto.NewProtocal("luda", r1, timeout)
s1 := p2p.NewServer(p1, timeout)
println("h1 监听 12345")
go s1.ListenAndServe("127.0.0.1:12345")

r2 := p2p.NewPRouter(timeout)  // 长连接路由
p2 := pto.NewProtocal("yoghurt", r2, timeout)
s2 := p2p.NewServer(p2, timeout)
println("h2 监听 12345")
go s2.ListenAndServe("127.0.0.1:12346")
```
#### 添加路由

``` go
p1.Add("yoghurt", "127.0.0.1:12346")
```

#### 发送数据
对于长连接，发送后没有返回值，对返回数据的监听与处理在protocal层维护的携程中实现

``` go
_, err := p1.Dispatch("yoghurt", msg)
```


## 实现
### Server层 通常意义上的Server功能
``` go
// 开启接口监听，将读到的数据传输给prtocal层解析
ListenAndServe(addr string) error
```

### Protocal层 handler功能
提供一个空的接口，需要用户来实现协议
``` go
type Protocal interface {
    // 解析请求通信内容,并返回数据,双工协议
    Handle(c net.Conn, msg []byte) ([]byte, error)
}
```
_example 中实现支持了一个长/短的协议状态机器
状态机为：
``` go
func (p *Protocal) Handle(c net.Conn, msg []byte) ([]byte, error) {
    cType := p.Router.GetConnType()
    req := &p2p.MsgPto{}
    resp := &p2p.MsgPto{}
    err := json.Unmarshal(msg, req)
    if err != nil {
        resp.Name = p.HostName
        resp.Operation = UnknownOp
        ret, _ := json.Marshal(resp)
        return ret, p2p.Error(p2p.ErrMismatchProtocalReq)
    }
    resp.Name = p.HostName
    switch req.Operation {
    case ConnectReq:
        subReq := &MsgGreetingReq{}
        err := json.Unmarshal(req.Data, subReq)
        if err != nil {
            return nil, p2p.Error(p2p.ErrMismatchProtocalResp)
        }
        if cType == p2p.ShortConnection {
            err = p.Router.AddRoute(req.Name, subReq.Addr)
        } else {
            if p.Router.AddRoute(req.Name, c) == nil {
                go p.IOLoop(c)
            }
        }
        if err != nil {
        }
        resp.Operation = ConnectResp
    case GetReq:
        resp.Operation = GetResp
    case FetchReq:
        resp.Operation = FetchResp
    case NoticeReq:
        resp.Operation = NoticeResp
    case ConnectResp:
        resp.Operation = GetReq
    case GetResp:
        resp.Operation = FetchReq
    case FetchResp:
        resp.Operation = NoticeReq
    case NoticeResp:
        return nil, nil
    default:
        resp.Operation = UnknownOp
    }
    ret, err := json.Marshal(resp)
    return ret, nil
}
```
长连接维护线程为：
``` go
// 长连接的话，需要在加入路由的时刻起携程 循环监控
func (p *Protocal) IOLoop(c net.Conn) {
    for {
        msg, err := p.read(c)
        if err != nil {
            // 在连接失败或者err=EOF(对方关闭连接之后，己方也要关闭)
            c.Close()
            return
        }
        resp, err := p.Handle(c, msg)
        if err != nil || resp == nil {
            continue
	    }
        c.SetWriteDeadline(time.Now().Add(p.to))
        _, err = c.Write(resp)
        if err != nil {
            return
        }
    }
}
```

### Router层 Client功能：发送数据
同样是接口，只不过这次提供了一套长/短连接的默认实现。
``` go
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
}
```

## TODO 需要解决的问题
1. 如何保存路由表？
允许重名，使用ip+server端口。
2. 对路由表更安全的操作？
目前用尽可能小的锁。
3. 在会话尚未结束的时候，修改conn，重新连接，之前未传完的协议状态机会失效。
允许失效，同时应该保证有有效conn存在的时候，且正在状态流转的时候不允许修改conn。
4. 同时添加对方路由，并同时向对方发送hello的时候，互换conn地址，无法统一。
没解决，打电话和视频聊天都会出现这种情况，他们也都没解决，让用户自己重试。
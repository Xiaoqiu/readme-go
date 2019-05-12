chapter12 网络编程基础在go
-
## 关于net/http,net,和 http.RoundTripper
- net/http package
- http.Get()
- http.ListenAndServer()
- chapter13 更详细介绍net package
- http.RoundTripper: 一个接口，可以执行HTTP交互。
  - 一个go元素可以根据一个http.Request,获取她的http.Response
### http.Response type
### http.Transport type
### http.Transport type

## 关于TCP/IP
## 关于IPV4 和 IPV6
## nc(1) 命令行基础
- 也叫做netcat(1) ： 是测试TCP/UDP服务器和客户端最有用的工具。
```bash
nc 10.10.1.123 1234

````
## 读取网络接口（网卡）配置
```go
package main
import (
  "fmt"
  "net"
)
func main()  {
  interfaces,err := net.Interfaces() // 返回所有的当前机器的网络接口，slice结构，元素类型是net.Interface
  if err != nil {
    fmt.Println(err)
    return
  }

  for _,i := range interfaces {
    fmt.Printf("Interface: %v\n", i.Name) // 每个网络接口的名称
    byName,err := net.InterfaceByName(i.Name)
    if err != nil {
      fmt.Println(err)
    }
  }

  address,err := byName.Addrs()
  for k,v := range address {
    // 一般有两个地址，一个是IPV4 一个是IPV6
    fmt.Printf("Interface Address #%v: %v\n", k, v.String()) // 每个网络接口的地址
  }
  fmt.Println()
}
```
- 不是所有的网络接口都有网络地址，
- 列举的网络接口不是一定都有硬件网络设备连接到这个接口
- 例如lOo接口，是回环设备。就是一个虚拟的网络接口，给本机和自己通信的。

```bash
$ go run netConfig.go
Interface: 本地连接* 1
Interface Address #0: fe80::44ab:ffcd:5085:26e/64
Interface Address #1: 169.254.2.110/16

Interface: 本地连接* 2
Interface Address #0: fe80::534:a603:7c5b:fb09/64
Interface Address #1: 169.254.251.9/16

Interface: VMware Network Adapter VMnet1
Interface Address #0: fe80::617a:6239:ad40:32d0/64
Interface Address #1: 192.168.226.1/24

Interface: VMware Network Adapter VMnet8
Interface Address #0: fe80::1c82:e338:4e72:905b/64
Interface Address #1: 192.168.247.1/24

Interface: WLAN
Interface Address #0: fe80::3dbd:6f47:9eec:dca/64
Interface Address #1: 192.168.8.167/24

Interface: 蓝牙网络连接
Interface Address #0: fe80::5c0f:f4ae:cda3:f331/64
Interface Address #1: 169.254.243.49/16

Interface: Loopback Pseudo-Interface 1
Interface Address #0: ::1/128
Interface Address #1: 127.0.0.1/8

```
- 下面的例子：获取每个网络接口的带宽
- net.Interface结构
```go
type Interface struct {
  Index int
  MTU int
  Name string
  HardwareAddr HardwareAddr
  Flags Flags
}
```
- 例子代码
```go
package main
import (
  "fmt"
  "net"
)
func main()  {
  interfaces,err := net.Interfaces()

  if err := nil {
    fmt.Println(err)
    return
  }

  for _,i := range interfaces {
    fmt.Printf()
    fmt.Println()
  }
}
```
## 执行DNS查询
## 创建web server使用go
## HTTP跟踪
## 创建一个web client使用go
## HTTP 链接超时
## wireshark 和 tshark 这两个工具
## 资源
## 练习
## 总结

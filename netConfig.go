package main
import (
  "fmt"
  "net"
)
func main()  {
  interfaces,err := net.Interfaces() // 返回所有的当前机器的所有网卡信息，slice结构，元素类型是net.Interface
  if err != nil {
    fmt.Println(err)
    return
  }

  for _,i := range interfaces {
    fmt.Printf("Interface: %v\n", i.Name) // 网络接口名称
    byName,err := net.InterfaceByName(i.Name)
    if err != nil {
      fmt.Println(err)
    }
    address,err := byName.Addrs()
    for k,v := range address { // 每个网络接口的地址
      fmt.Printf("Interface Address #%v: %v\n", k, v.String())
    }
    fmt.Println()

  }
}

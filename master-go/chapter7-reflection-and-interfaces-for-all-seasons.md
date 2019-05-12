chapter 7  反射和接口
## type method
## go 接口
## 关于类型断言-type assertion
## 开发自己的接口
## 反射
## go的面向对象编程
## 资源
## 练习
## 总结

## type method
- 方法名前面有一个参数，这个参数叫做方法的接受者，就是说这个方法发送一个消息到这个对象。一般使用一个字母定义这个参数。
- 类型方法和普通方法的调用方式不同
- 类型方法个接口是相关的。
```go
// methods.go
package main
import (
  "fmt"
)
type twoInts strut {
  X int64
  Y int64
}

func regularFunction(a,b twoInts) twoInts  {
  temp := twoInts{X:a.X + b.Y, Y:a.Y + b.X}
  return temp
}

func (a twoInts)method(b twoInts)twoInts  {
  // 我们能在这个方法里面用到这个对象a
  temp := twoInts{X:a.X + b.Y, Y:a.Y + b.X}
  return temp
}

func main()  {
  i := twoInts{X:1,Y:2}
  j := twoInts{X:5,Y:-2}
  // 普通方法直接调用
  // 类型方法是通过一个对象调用
  fmt.Println(regularFunction(i,j))
  fmt.Println(i.method(j))

}
```
## Go 接口
- 接口是一种抽象类型
- 接口定义了一些方法，用来定义其他类型的行为。（组合设计模式）
  - 把行为抽象出来，需要时候组合
  - 设计原则，针对接口编程
- 类型是接口的一个实例
- 类型满足一个接口必须实现接口定义的所有的方法。
- 接口的使用：方法以接口为参数，那么可以传入实现了这个接口的类型。（设计原则，针对接口编程）
- 如果你定义了一个接口和实现在同一个go package，那么你使用错了接口
```go
type Reader interface {
  Read(p []byte)(n int, err error)
}

type Writer interface {
  Write(p []byte)(n int, err error)
}
```
- 接口应该是基础的，实现一种行为

## 关于类型断言-type assertion
-

## 开发自己的接口
## 反射
## go的面向对象编程
## 资源
## 练习
## 总结

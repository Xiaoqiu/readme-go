chapter2 理解go的内部
- go编译器
- 垃圾回收原理
- 在go里面访问C语言
- 在C语言里面访问go函数
- panic() recover() 函数
- unsafe package
- defer关键字
- strace(1) linux 基础工具
- dtrace(1) FreeBSD系统
- go环境变量信息
- node tree
- go 汇编

## go 编译
```bash
## 编译成object file： 机器编码二进制文件
go tool compile unsafe.go
ls -l unsafe.o
file unsafe.o

## 编译成archive file : 多个文件的组合成的一个文件
go tool compile -pack unsafe.go
ls -l unsafe.a
file unsafe.a

## 查看archive file里面的内容
ar t unsafe.a

## 检查竞争条件，第十章并发的时候很有用
go tool compile -race

## 这个编译出来很多文件，只是一般会被隐藏
go tool compile -S unsafe.go
```
## 垃圾回收
- go pacakge 提供了函数允许你学习垃圾回收器
```go
// gColl.go
package main
import (
  "fmt"
  "runtime"
  "time"
)
func printStats(mem runtime.MemStats)  {
  runtime.ReadMemStats(&mem)// 获取最近垃圾回收的统计
  fmt.Println("mem.Alloc: ", mem.Alloc)
  fmt.Println("mem.TotalAlloc:",mem.TotalAlloc)
  fmt.Println("mem.HeapAlloc:",mem.HeapAlloc)
  fmt.Println("mem.NumGC:",mem.NumGC)
  fmt.Println("-----")
}

func main(){
  var mem runtime.MemStats
  printStats(mem)
  for i := 0, i < 10 i++ {
    s := make([]byte,50000000) // 创建一个大的slice，触发垃圾回收
    if s == nil {
      fmt.Println("Operation failed!")
    }
  }
  printStats(mem)

  // 最后一部分: 继续分配内存
  for i := 0， i < 10; i++ {
    s := make([]byte,100000000)
    if s == nil {
      fmt.Println("Operation failed!")
    }
    time.Sleep(5 * time.Second)
  }
  printStats(mem)
}
```
```bash
$ go run gColl.go
## 输出：
mem.Alloc: 66024
mem.TotalAlloc: 66024
mem.HeapAlloc: 66024
mem.NumGC: 0
-----
mem.Alloc: 50078496
mem.TotalAlloc: 500117056
mem.HeapAlloc: 50078496
mem.NumGC: 10
-----
mem.Alloc: 76712
mem.TotalAlloc: 1500199904
mem.HeapAlloc: 76712
mem.NumGC: 20
-----

## 收集更详细的垃圾回收数据, 会打印分析数据关于垃圾回收
GODEBUG=gctrace=1 go run gColl.go
## 输出：
gc 4 @0.025s 0%: 0.002+0.065+0.018 ms clock, 0.021+0.040/0.057/0.003+0.14 ms cpu, 47->47->0 MB, 48 MB goal, 8 P
gc 17 @30.103s 0%: 0.004+0.080+0.019 ms clock, 0.033+0/0.076/0.071+0.15 ms cpu, 95->95->0 MB, 96 MB goal, 8 P

## 47->47->0:
第一个参数：垃圾回收器运行之前堆大小，
第二个参数：垃圾回收器停止的时候堆大小
第三个参数：live heap的值
```
### triColor 算法
- 一个垃圾回收算法
-
### 更多关于垃圾回收操作
### unsafe code
- 不具备类型安全和内存安全的代码，不推荐使用,一般和指针相关
  - unsafe.Pointer类型可以重写go的类型系统。非常危险。
- 实例
```go
package main
import(
  "fmt"
  "unsafe"
)
func main()  {
  var value int64 = 5
  var p1 = &value
// 创建一个int32的指针，指向int64变量（变量名为value）
// 使用p1指针获取该变量value。
// 任何的指针都可以转换位unsafe.Pointer
  var p2 = (*int32) (unsafe.pointer(p1))

//第二部分
fmt.Println("*p1: ", *p1)
fmt.Println("*p2: ", *p2)
*p1 = 5434123412312431212 // 32位的指针不可以存64位的整数
fmt.Println(value)
fmt.Println("*p2: ", *p2)
*p1 = 54341234
fmt.Println(value)
fmt.Println("*p2: ", *p2)
//You can dereference a pointer and get, use, or set its value using the star character (*).

}
```

```bash
$ go run unsafe.go
*p1:  5
*p2:  5
5434123412312431212
*p2:  -930866580
54341234
*p2:  54341234

```
### 关于unsafe package
- 源码位置：/usr/local/Cellar/go/1.9.1/libexec/src/unsafe/unsafe.go
```bash
$ cd /usr/local/Cellar/go/1.9.1/libexec/src/unsafe/
$ grep  -v '^//' unsafe.go | grep -v '^$'
package unsafe
type ArbitraryType int
type Pointer *ArbitraryType
func Sizeof(x ArbitraryType) uintptr
func Offsetof(x ArbitraryType) uintptr
func Alignof(x ArbitraryType) uintptr

```
- 在import导入这个unsafe package的时候go编译器实现unsafe package
- 底层的package都会用到这个unsafe package : runtime, syscall,os


### 一个unsafe package 例子
- & : 取地址
- * : 解析地址(就是取值)
- 例子：使用指针获取素有的数组元素
```go
package main
import(
  "fmt"
  "unsafe"
)
func main()  {
  array := [...]int{0,1,-2,3,4}
  pointer := &array[0] // 取地址：array[0]的地址，第一个元素的地址
  fmt.Print(*pointer, " ") //打印值
  memoryAddress := uinptr(unsafe.Pointer(pointer)) + unsafe.Sizeof(arrray[0])

  for i:=0, i<len(array)-1; i++ {
    pointer = (*int)(unsafe.Pointer(memoryAddress))
    fmt.Print(*pointer," ") //取值
    memoryAddress = uinptr(unsafe.Pointer(poiter)) + unsafe.Sizeof(array[0]) //指针+一个偏移量
  }

// 最后我们使用指针和内存地址，获取一个数组里面不存在的元素。
// 编译器不会捕获这个逻辑错误，因为使用了unsafe package,
//这就会返回一些不安全的参数
  fmt.Println()
  pointer = (*int)(unsafe.Pointer(memoryAddress))
  fmt.Print("One more: ", *pointer, " ") //取值
  memoryAddress = uinptr(unsafe.Pointer(pointer)) + unsafe.Sizeof(array[0])
  fmt.Println()
}
```
```bash
$ go run moreUnsafe.go
0 1 -2 3 4
One more: 842350722816

```

## 在go里面访问C语言
### 在go里面访问C语言，两者同一个文件中
- 这种方式如果你经常用到，请重新考虑你选取的语言是否合适
- 例子：
- 注意，你的c代码是写在注释里面的，引入了c package后，这些代码可以被识别
```go
// cGo.go
package main

//#include <stdio.h>
//void callc() {
//  printf(“Calling C code!\n”);
// }
import "C"

import "fmt" // 其他的package和 c package 应该分库
func main()  {
  C.callC()
  fmt.Println("A Go statements!")
}
```

```bash
$ go run cGo.go
A Go statement!
Calling C code!
Another Go statement!
```

### 在go里面访问C语言，两者不同的文件中
### C code
```c
// callClib/callC.h
#ifndef CALLC_H
#define CALLC_H

void cHello()
void printMessage(char* message)

# endif
```
```c
// 路径：callClib/callc.h
#include <stdio.h>
#include "callC.h"
void cHello(){
  printf("Hello from C !\n");
}
void printMessage(char* message) {
  printf("Go send me %s\n", message);
}
```

### Go code
- c package告诉go build 去处理输入文件使用cgo 工具，在go 编译器处理go文件之前。
- 仍然需要使用注释去告诉Go程序，关于C代码。
- 这个例子，告诉go代码去找到两个c代码文件。

```go
package main
// #cgo CFLAGS: -I${SRCDIR}/callClib
// #cgo LDFLAGS: ${SRCDIR}/callC.a
// #include<stdlib.h>
// #include<callC.h>
import "C"
import (
  "fmt"
  "unsafe"
)
func main()  {
  fmt.Println("Go to call a C function")
  C.cHello()
  fmt.Println("Going to call another C function!")
  // 创建一个可以传递给C语言的字符串
  myMessage := C.CString("This is Mihalis!")
  //需要释放C string的内存，当这个内存不在需要
  //释放内存需要一个defer的声明，和一个unsafe.Pointer
  //unsafe.Pointer，是创建一个指针可以操控这个变量的内存
  defer C.free(unsafe.Pointer(myMessage))
  C.printMessage(myMessage)

  fmt.Println("All perfectly done!")
}
```

### Mixing Go and C code
- 执行Go代码来调用C语言
- 首先需要编译C代码，生产一个library
```bash
ls -l callClib/
## 把C文件编译为Object文件
gcc -c callClib/*.c
ls -l callC.o

## 把几个库文件生成一个archive file
ar rs callC.a *.o
ls -l callC.a

## 查看生产的的archive文件内容
file callC.a

## 最后删除原理的库文件
rm callC.o

## 编译Go文件
go build callC.go
ls -l callC

## 查看编译Go代码后生成的执行文件
file callC

## 执行Go生产的执行文件
./callC

```
- 如果是调用简单的C语言，把C和Go写在同一个文件中是合适的。
- 如果是实现复杂和高级的功能还是需要创建一个静态的C库
## 在C语言里面访问go函数
- 例子里，c语言调用了两个go的方法
- go package会转化为C共享库，以这种方式被C语言调用
### Go package

### C code

## defer关键字
## panic和recover
## 两个比较常用的Unix基础服务
## 设置go开发运行环境
## go汇编
# 节点树
# 学习更多关于go build
## 一般go编码建议
## 资源附件
## 练习
## 总结

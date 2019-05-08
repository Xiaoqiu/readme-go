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
// todo

```go
package main
import "C"

import (
  "fmt"
)

//export PrintMessage
func PrintMessage()  {
  fmt.Println("a go function!")
}

//export Multiplay
func Multiplay(a,b int) int  {
  return a*b
}

func main()  {

}
```

```bash
go build -o usedByC.o -buildmode=c-shared usedByC.go
## 生成usedByC.h usedByC.o
ls -l usedByC.*
file usedByC.o
## 显示是一个linked shared library

```
### C code
```c
// willUseGo.c
# include <stdio.h>
# include "usedByC.h" // 这是C语言的法则，可以知道一个库里面可以用的方法
// GoInt p是从go函数中获取整型，可以用(int)p 转换为C的整型
int main(int argc,char **argv) {
  GoInt x = 12;
  GoInt y = 23;
  printf("about to call a go function \n", );
  PrintMessage();


  GoInt p = Multiplay(x,y);
  printf("product: %d\n", (int)p);
  printf("it worked!\n", );
  return 0;
}
```
```bash
gcc -o willUseGo willUseGo.c ./usedByC.o
./willUseGo

```

## defer关键字
- 延迟一个函数，等到外层的函数返回后。
- 一个函数里有三个defer,是按照“后进先出”的顺序执行，就是倒着执行
- 例如：三个顺序defer函数：f1(), f2(), f3()
- 执行顺序：f3()->f2()->f3()
- 使用
  - 关闭文件
  - panic()/recover()
  - 这里比较特殊的例子
 ```go
 // defer.go
 package main
 import (
   "fmt"
 )
 func d1()  {
    for i := 3; i > 0; i--{
      defer fmt.Print(i," ")
    }
 }
 func d2()  {
  for i := 3; i>0; i-- {
    defer func ()  { // 应用在一个匿名函数，而且匿名函数没有参数
      fmt.Print(i," ") // i 在外层函数执行完，就是循环执行完就已经不存在了。所以当defer在外层函数返回时执行，得到的i值就是0。 特别注意，这是一个程序的bug
    }()
  }
  fmt.Println()
 }

 func d3(){
   for i := 3; i >0 ; i-- {
      defer func (n int)  { // 匿名函数，有一个参数.每次定义的时候，都会使用当前的i值。而不是循环后去获取的i值
         fmt.Print(n," ")
      }(i)
   }
 }

 func main() {
   d1()
   d2()
   fmt.Println()
   d3()
   fmt.Println()
 }
 ```

 ```bash
 $ go run defer.go
1 2 3
0 0 0
1 2 3
 ```
## panic和recover
-  panic() ： 终止当前的Go程序工作流，开始paniking
- recover() ： 一个goroutine已经调用了panic(),这个函数可以获取回来这个goroutine的控制权。

```go
// panicRecover.go
package main
import (
  "fmt"
)
func a() {
  fmt.Println("Inside a()")

  // defer定义了一个匿名函数，如果有一个panic()被调用，这个匿名函数会被调用。
  defer func ()  {
    if c := recover(); c != nil {
      fmt.Println("recover inside a()!")
    }()
    fmt.Println("About to call b()")
    b()
    fmt.Println("b() exited!")
    fmt.Println("Exiting a()")
  }
}

func b(){
  fmt.Println("Inside b()")
  panic("Panic in b()!")
  fmt.Println("Exiting b()")
}

func main()  {
  a()
  fmt.Println("main() ended!")
}

```

```bash
go run panicRecover.go
Inside a()
About to call b()
Inside b()
Recover inside a()!
main() ended!

```
- 这是挺神奇的，a()函数没有执行完，
- 没有抛出paniking的异常，正常结束了，这是我们想要的结果
- 这是个因为外层函数a()里面定义了defer延迟执行的匿名函数。处理了b()函数抛出来的panic()

### 自己使用一些这个panic函数
- 可以使用panic函数，不试图去recover
- 这个panic会抛出异常，终止程序。所以使用panic和recover配合使用比较专业一点。
- panic()函数的输入和log package的panic()函数输入看起来很像，但是很重要的是，这里的panic()是不会发送任何日志到Unix机器的。

- os.Args[0]: 第一个参数固定是文件路径：/var/folders/rt/xbcdvxgn43s2p0n9wn_z6bz80000gn/T/go-build563527000/b001/exe/justPanic
- os.Args[1]: 第二个参数开始是用户输入
- 所以：os.Args[]至少长度为1
```go
// justPanic.go
package main
import (
  "fmt"
  "os"
)
func main()  {
  if len(os.Args) == 1 {
    panic("Not enough arguments!") // 这个panic会抛出异常，终止程序
  }

  fmt.Println("Thanks for the arguments!")
}
```
```bash
$ go run justPanic.go
panic: Not enough arguments!

goroutine 1 [running]:
main.main()
      /Users/mtsouk/ch2/code/justPanic.go:10 +0x9e
exit status 2
```
## 两个比较常用的Unix基础服务
- 排查Unix程序的运行效率或者异常终止问题的工具
- 这个工具可以查看C 系统调用，把这个执行文件调用的。
- strace(1) / dtrace(1)
- 要理解：所有的程序，只要在Unix上面运行，都是使用C语言的系统调用，去和Unix内核通信的，和执行几乎所有的任务。
- 虽然两个工具都可以在go run命令下使用。go run 创建很多临时文件，这个两个工具会试图像是这些临时文件，这不是你想要的。

### strace tool
- 只能在linux系统使用这个工具，macos是不能使用的
- 显示函数输入参数，返回结果。
- 自己终止，并利用grep命令获取你要结果
- 加-c格式化输出结果：系统调用名称，系统调用次数，产生的错误次数，占用总执行时间的百分比，一个系统调用总的时间

```bash
$ strace ls
execve("/bin/ls", ["ls"], [/* 15 vars */]) = 0
brk(0)                                  = 0x186c000
fstat(3, {st_mode=S_IFREG|0644, st_size=35288, ...}) = 0

$ strace find /usr 2>&1 | grep ioctl
ioctl(0, SNDCTL_TMR_TIMEBASE or SNDRV_TIMER_IOCTL_NEXT_DEVICE or TCGETS, 0x7ffe3bc59c50) = -1 ENOTTY (Inappropriate ioctl for device)
ioctl(1, SNDCTL_TMR_TIMEBASE or SNDRV_TIMER_IOCTL_NEXT_DEVICE or TCGETS, 0x7ffe3bc59be0) = -1 ENOTTY

(Inappropriate ioctl for device)

$ strace -c find /usr 1>/dev/null
% time     seconds  usecs/call     calls    errors syscall
------ ----------- ----------- --------- --------- -------------
 82.88    0.063223           2     39228           getdents
 16.60    0.012664           1     19587           newfstatat
  0.16    0.000119           0     19618        13 open

```

### the dtrace tool
- strace(1),truss(1)可以跟踪进程的系统调用，但是他们很慢
- dtrace(1)可以监控运行程序，或者服务进程，开销没有那么大
- dtrace(1)有可以在Linux运行的版本，但是最好的运行环境是macos,和其他的freeBSD实体
- dtrace(1),dtruss(1)需要root运行
- dtrace是远远比strace强大的
- 必须熟练一种这样的工具
```bash

$ sudo dtruss godoc
ioctl(0x3, 0x80086804, 0x7FFEEFBFEC20)           = 0 0
close(0x3)         = 0 0
access("/AppleInternal/XBS/.isChrooted\0", 0x0, 0x0)   = -1 Err#2
thread_selfid(0x0, 0x0, 0x0)         = 1895378 0
geteuid(0x0, 0x0, 0x0)         = 0 0
getegid(0x0, 0x0, 0x0)         = 0 0


$ sudo dtruss -c go run unsafe.go 2>&1

CALL                                        COUNT
access                                          1
bsdthread_register                              1
getuid                                          1
ioctl                                           1
issetugid                                       1
kqueue                                          1
write                                           1
read                                          244
kevent                                        474
fcntl                                         479
lstat64                                       553

```

## 认识你的go开发运行环境
- 使用runtime pacakge 去获取go的环境信息
- 这个package包含函数和属性，获取这些信息
```go
//goEnv.go
package main
import (
  "fmt"
  "runtime"
)

func main()  {
  fmt.Print("you are using: " runtime.Compiler, " ")
  fmt.Println("on a", runtime.GOARCH, "machine")
  fmt.Println("go version : ",runtime.Version())
  fmt.Println("Number of CPUs : ", runtime.NumCPU())
  fmt.Println("Number of Goroutines: ", runtime.NumGoroutine())
}
````
```bash
katedeMacBook-Pro:master-go kate$  go run goEnv.go
you are using: gc on a amd64 machine
go version :  go1.11.1
Number of CPUs :  4
Number of Goroutines:  1
```

- 例子：判断go版本，并提示
```go
// requiredVersion.go
package main
import (
  "fmt"
  "runtime"
  "strconv"
  "strings"
)

func main()  {
  myVersion := runtime.Version()
	//取 go1.11.1中的go1, 再取go1中的1。
  major := strings.Split(myVersion,".")[0][2]
  fmt.Println("major: ",strings.Split(myVersion,".")[0][2])
	// 取 go1.11.1中的11
  minor := strings.Split(myVersion,".")[1]
  fmt.Println("minor",minor)
  m1,_ := strconv.Atoi(string(major))
  m2,_ := strconv.Atoi(minor)
  fmt.Println("m1",m1)
  fmt.Println("m2",m2)
  if m1 ==1 && m2 <8 {
    fmt.Println("need to use 1.8 or higher!")
    return
  }
  fmt.Println("you are using go version 1.8 or higher!")
}
```

## go汇编
- 汇编语言
- go 汇编器：一个go的工具允许你去到go编译器使用的汇编语言
- GOOS： 操作系统
- GOARCH：编译架构
-
```bash

GOOS=darwin GOARCH=amd64 go tool compile -S goEnv.go

GOOS=darwin GOARCH=amd64 go build -gcflags -S goEnv.go
```
# 节点树
- go编译器生成的节点树

# 学习更多关于go build
```bash
go build -x defer.go
```
## 一般go编码建议
## 资源附件
## 练习
## 总结

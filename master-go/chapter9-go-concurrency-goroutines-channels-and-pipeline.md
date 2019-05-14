chapter 9 go 并发送-goroutines,channels,pipelines


- 创建channels
- 从一个channel读取和接收数据
- 写入或者发送数据去一个channel
- 创建pipeline
- 等待giroutines返回
- 等待你的goroutines结束


### 进程，线程，goroutine的不同
### go定时器
	- Unix定时器负责执行线程程序
	- go运行环境有自己的定时器，负责执行goroutine使用一种m:n定时
		- m个goroutine执行的时候，使用n个系统线程，利用复用性。
		- 这个定时器，负责goroutine执行的方式和顺序
		- go定时器只是处理一个go程序的goroutines,这样比内核定时器的更快，简单，开销小。
		- 下一章更详细讲述这个go定时器

- 并发和并行
	- 并行是同事发生
	- 并发是独立执行
	- 只有组件可以独立执行，你可以安全地并行。
	- 更好的并行，来自如更好的并发实现。
	- 程序设计负责实现良好的并发性，能利用系统的并行优势。
	- 程序设计不用考虑并行，但是要实现把事情设计为独立，小的实现，然后组合起来运行。
	- 即使不在一个并行的Unix机器，你并发的设计仍然提高了程序的可维护性。
	- 并发是比并行要好的。
## goroutines
### 创建goroutines
	- 执行顺序决定于：操作系统的定时器，go定时器，操作系统开销
	- go + 函数名或者一个匿名函数
	- 这个go关键字的函数，立即返回，后台执行一个goroutine,主线程继续执行。
	- 使用正常函数创建goroutine
	- 是使用匿名函数创建goroutine

```go
// simple.go
package main
import (
	"fmt"
	"time"
)
func function()  {
	for i:=0; i<10; i++ {
		fmt.Print(i)
	}
	fmt.Println()
}

func main()  {
	//使用一般函数，启动一个goroutine
	go function()

	// 使用匿名函数启动一个goroutine，这种方式更实际如果实现小函数
	go func(){
		for i := 10; i<20; i++ {
			fmt.Print(i," ")
		}
	}()

//等待gorouitne结束，打印到输出。再结束main函数
// 如果没有这一行，结果不会出现在控制台
	time.Sleep(1 * time.Second)
}
```
```bash
# 两次执行结果显示，两个goroutine的执行顺序是随机的。
bogon:master-go kate$ go run simple.go
0 1 2 3 4 5 6 7 8 9
10 11 12 13 14 15 16 17 18 19

bogon:master-go kate$ go run simple.go
10 11 12 13 14 15 16 17 18 19 0 1 2 3 4 5 6 7 8 9

```
- 下一章会有写另外的代码控制goroutine的执行顺序
- 也会控制打印一个goroutine的结果，然后再打印下一个。

### 创建多个goroutines
- 使用flag package 读取传入参数
```go
// create.go
// go run create.go -n 100
package main
import (
	"flag"
	"fmt"
	"time"
)
func main()  {
	n := flag.Int("n", 10, "Number of goroutines")
	// 使用flag package 读取传入参数
	flag.Parse()

	count := *n
	fmt.Printf("Going to create %d goroutines.\n" , count)

	for i := 0; i<count; i++ {
		go func (x int)  {
			fmt.Printf("%d", x)
		}(i)
	}

	//等待goroutine完成后，结束main函数
	//还有其他更加优雅的方式等待goroutine结束，在main()函数之前。
	// 这种sleep方式是有风险的，很有可能看不到gorountine的输出。
	time.Sleep(time.Second)
	fmt.Println("\nExiting...")

}
````
```bash
bogon:master-go kate$ go run create.go -n 100
Going to create 100 goroutines.
1 4 2 3 7 9 8 13 10 11 12 18 14 15 16 17 5 20 19 24 21 22 23 31 25 26 27 28 29 30 35 32 33 6 45 36 37 34 38 39 40 41 44 46 42 58 59 64 43 49 48 47 60 61 62 50 51 52 53 54 55 56 57 81 65 66 67 68 69 70 71 72 73 74 75 76 77 78 79 80 63 87 95 83 84 89 86 90 82 88 85 91 92 93 0 94 97 98 99 96
Exiting...

```


## 等待goroutines结束
- 防止main()在goroutine之前结束。
- 使用sync package
- 高级章：学习两种技术应对，超时的goroutines
```go
// syncGo.go
package main
import(
	"flag"
	"fmt"
	"sync"
)
func main()  {
	n := flag.Int("n", 20, "Number of gorouitnes")
	flag.Parse()
	count := *n
	fmt.Printf("Going to create %d gorouitnes. \n", count)
	var waitGroup sync.WaitGroup

	fmt.Printf("%#v\n", waitGroup)

	for i := 0; i<count; i++ {
		//每个gorouitne启动之前增加计数器的值
		waitGroup.Add(1)

		go func (x int)  {
			//每个goroutine结束之后减少计数器的值
			defer waitGroup.Done()
			fmt.Printf("%d ",x)
		}(i)
	}

	fmt.Printf("%#v\n", waitGroup)
	waitGroup.Wait()
	fmt.Println("\nExiting....")

}
````
- state1在sync.WaitGroup保存计数器，这个计数器用来增加或者减少，goroutine数量。根据sync.Add()，sync.Done()

```bash
bogon:master-go kate$ go run syncGo.go -n 10
Going to create 10 gorouitnes.
sync.WaitGroup{noCopy:sync.noCopy{}, state1:[3]uint32{0x0, 0x0, 0x0}}
sync.WaitGroup{noCopy:sync.noCopy{}, state1:[3]uint32{0x0, 0xa, 0x0}}
9 5 6 7 8 0 3 4 1 2
Exiting....
bogon:master-go kate$ go run syncGo.go -n 20
Going to create 20 gorouitnes.
sync.WaitGroup{noCopy:sync.noCopy{}, state1:[3]uint32{0x0, 0x0, 0x0}}
1 0 3 8 5 6 7 10 9 11 15 19 sync.WaitGroup{noCopy:sync.noCopy{}, state1:[3]uint32{0x0, 0xa, 0x0}}
16 17 18 4 13 14 12 2
Exiting....
bogon:master-go kate$ go run syncGo.go -n 30
Going to create 30 gorouitnes.
sync.WaitGroup{noCopy:sync.noCopy{}, state1:[3]uint32{0x0, 0x0, 0x0}}
1 6 2 3 4 5 9 7 8 11 10 sync.WaitGroup{noCopy:sync.noCopy{}, state1:[3]uint32{0x0, 0x14, 0x0}}
29 13 0 22 21 23 25 19 18 24 26 27 20 28 14 16 15 12 17
Exiting....
```
### 如果Add()和Done()发起这两个请求的数量不一致，会发生什么？（创建goroutine的时候调用Add()，执行完成没有调用Done()）
- 情况1 ： 在调用第一次打印之前多执行一次waitGroup.Add(1),
	- 结果就会死锁，因为程序在等待n+1个goroutine结束，但是只有n个结束了。所以一直等待fatal error: all goroutines are asleep - deadlock!
s
- 情况2 ：在循环后面加多一行waitGroup.Done()。
	- 结果：panic: sync: negative WaitGroup counter.

## channels
- 1 一个channel只允许一种类型的数据
- 2 操作一个channel需要一个接收信息方
- 3 定义channel： chan关键字
- 4 关闭channel: close()函数
- 5 channel作为函数参数，一定要声明方向，是写入还是读取。这样省了很多错误。

### 写入一个channel
- 值：x和channel:c类型相同
- 写入：c <- x
```go
package main
import (
	"fmt"
	"time"
)
//chan声明c是一个channel, 后面是channel类型
func writeToChannel(c chan int, x int)  {
	fmt.Println(x)
	c <- x
	close(c) // 关闭channel，
	fmt.Println(x)
}

func main()  {
	// 创建channel，指定类型
	c := make(chan int)
	go writeToChannel(c,10)
	time.Sleep(1 * time.Second)
}

```
````bash
$ go run writeCh.go
10
````
- 当c <- x执行的时候，writeToChannel()函数剩下的部分就被阻塞了。当写入一个channel的时候，这个channel没有接收者，当main()函数time.Sleep(1 * time.Second)结束后，也没有等待writeToChannel()执行完成。

### 读取一个channel
```go
// readCh.go
package main
import (
	"fmt"
	"time"
)
func writeToChannel(c chan int, x int)  {
	fmt.Println("1", x)
	c <- x
	close(c)
	fmt.Println("2",x)
}
func main()  {
	c := make(chan int)
	go writeToChannel(c, 10)
	time.Sleep(1 * time.Second)
	fmt.Println("Read:", <-c)
	time.Sleep(1 * time.Second)

	_, ok := <-c
	if ok {
		fmt.Println("Channel is open!")
	} else {
		fmt.Println("Channel is closed!")
	}
}
```
- writeToChannel()函数里的channel有接受者，所以不会阻塞。能继续执行第二个打印。

```bash
$ go run readCh.go
1 10
Read: 10
2 10
Channel is closed!
$ go run readCh.go
1 10
2 10
Read: 10
Channel is closed!
```

### channel 作为一个函数参数
- 作为参数是可以声明channel方向的。
- channel默认是默认是双向的。
```go
func f1(c chan int, x int)  {
	fmt.Println(x)
	c <- x
}

//chan<- 声明这是个只能写入的channel
// 读取这个channel就会报错
func f2(c chan<- int, x int) {
	fmt.Println(x)
	c <- x
}

```
## pipelines
- 是一个虚拟的概念，就是利用channel来传递数据，在不同的goroutines之间。
```go
package main
import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)
var CLOSEA = false

var DATA = make(map[int]bool)

func random(min, max int) int {
	return rand.Intn(max-min) + min
}
// 产生随机数，并写入channelA
func first(min, max int, out chan<- int)  {
	for {
		if CLOSEA {
			close(out) // 只有写入channel可以关闭
			return
		}
		out <- random(min,max)
	}
}

//从channelA读取数据并写入DATA
//并且把数据写入channelB
func second(out chan<- int, in <-chan int)  {
	for x := range in {
		fmt.Print(x," ")
		// 如果发现一个随机数已经存在DATA map里面，就设置CLOSEA=true,停止first()函数再写入任何数据到channelA.
		_,ok := DATA[x]
		if ok {
			CLOSEA = true
		} else {
			DATA[x] = true // DATA没有这个参数，就设置这个参数为key,value为true
				out <- x
		}
	}

	fmt.Println()
	close(out) // 数据写完可以关闭channelB

}

//从channelB读取数据，并求和
func third(in <-chan int) {
	var sum int
	sum = 0
	// 当second()函数关闭channelB,就停止读取的循环。
	for x2 := range in {
		sum = sum + x2
	}

	fmt.Printf("The sum of the random numbers is %d\n", sum)
}


func main(){
	if len(os.Args) != 3 {
		fmt.Println("Need two integer parameters!")
		os.Exit(1)
	}

	n1, _ := strconv.Atoi(os.Args[1])
	n2, _ := strconv.Atoi(os.Args[2])

	if n1 > n2 {
		fmt.Printf("%d should be smaller than %d\n",n1, n2)
		return
	}

	rand.Seed(time.Now().UnixNano())
	A := make(chan int)
	B := make(chan int)

	go first(n1, n2, A)
	go second(B, A)
	third(B)

}

```
## addtional resource
## 练习
## 总结

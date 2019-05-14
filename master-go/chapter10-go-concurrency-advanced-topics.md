chapter 10 并发-高级话题
- select关键字
- go定时器怎么工作的
- 两个技术让运行过久的goroutine超时
- signals channels
- buffered channels
- nil channels
- 监控 gorouitnes
- channels of channels
- 共享内存和互斥
- sync.Mutex 和 sync.RWMutex类型
- context package和函数性（泛函数）
- worker pools
- 检查竞争态

## 重温一下go的定时器-scheduler
- 有效率的方式分发工作量到可用资源上面
- fork-join 并发模型
- 公平分配策略：所有处理器分配平均一样的开销
- goroutine是一个task.
- go定时器工作的三类实体：
	- 操作系统线程
	- gorouitne
	- 逻辑处理器：个数由GOMAXPROCS环境变量定义。
- 两种队列：
	- gorouitne全局队列：
	- gorouitne本地队列：依附于每个逻辑处理器。
- 全局队列的gorouitne分配给逻辑处理器的本地队列。
- 定时器需要检查全局队列，为了避免执行的goroutine只存在本地队列。
- 每个逻辑处理器有多个线程，所以两个本地队列之前就会发生偷窃.
- go定时器是可以创建操作系统线程的。

### GOMAXPROCS环境变量
- 可以限制系统线程数，用来执行go
```go
package main
import (
	"fmt"
	"runtime"
)
// 获取环境变量：GOMAXPROCS的值
func getGOMAXPROCS() int {
	return runtime.GOMAXPROCS(0)
}
func main()  {
	fmt.Printf("GOMAXPROCS: %d\n", getGOMAXPROCS())
}
```
```bash
# 运行
go run maxprocs.go
GOMAXPROCS:8

# 改变环境变量的值
go version

export GOMAXPROCS=800; go run maxprocs.go

export GOMAXPROCS=4; go run maxprocs.go

```

## select 关键字
- 允许goroutine去等待其他操作
- select是同步的，如果没有一个在select下面声明的channel是准备好的，那么select声明会阻塞，直到有一个channel是准备好的。如果有多个select的channel是准备好的，那么会随机选择一个装备好的channel.
- select声明最大的优势：可以连接，编排，管理多个channel.
- select声明可以连接channel, channel可以连接goroutine,所以select声明是go并发模型最重要部分。

```go
// select.go
package main
import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func gen(min, max int, createNumber chan int, end chan bool)  {
	for {
		select {
		case createNumber <- rand.Intn(max-min) + min:
		case <-end:
			close(end)
			return
		case <-time.After(4 * time.Second):// 这个函数4秒后自动返回。所以会取消select的阻塞。
			fmt.Println("\ntime.After()!")
		}
	}
}

func main()  {
	rand.Seed(time.Now().Unix())
	createNumber := make(chan int)
	end := make(chan bool)

	if len(os.Args) != 2 {
		fmt.Println("Please give me an integer!")
		return
	}

	n,_ := strconv.Atoi(os.Args[1])
	fmt.Printf("Going to create %d random numbers.\n" ,n)

	go gen(0, 2*n, createNumber, end)

	for i := 0; i<n; i++ {
		fmt.Printf("%d",<-createNumber)
	}

	time.Sleep(5 * time.Second)// 给足够的时间case <-time.After(4 * time.Second):可以返回，这个声明可以生效。
	fmt.Println("Exiting...")
	end <- true // 激活case <-end:的select声明。

}

```
```go
$ go run select.go 10
Going to create 10 random numbers.
13 17 8 14 19 9 2 0 19 5
time.After()!
Exiting...
```
## goroutine超时
- 两个最重要的方式，让goroutine超时。
- 可以控制，多久让一个goroutine停止。
- 两种方式都是使用select关键字和time.After()函数结合。
- 方式一：让一个goroutine超时
```go
// timeOut1.go
package main
import(
	"fmt"
	"time"
)

func main()  {
	c1 := make(chan string)

	// 需要3秒才能执行完这个goroutine
	go func() {
		time.Sleep(time.Second * 3)
		c1 <- "c1 ok"
	}()

	// 使用select阻塞main()函数
	select {
	case res := <-c1
		fmt.Println(res)
	case <-time.After(time.Second * 1): // 1 秒就返回，结束select的阻塞。继续执行。main()执行完，就不等待goroutine的结果。
		fmt.Println("timeout c1")
	}

	c2 := make(chan string)
	go func ()  {
		 time.Sleep(3 * time.Second)
		 c2 <- "c2 ok"
	}()

	select {
	case res := <-c2:
		fmt.Println(res)
	case <-time.After(4 * time.Second): // select阻塞后，有足够的时间给第一个case解除阻塞往下执行，所以这个case不会执行。
		fmt.Println("timeout c2")
	}
}
```
```bash
go run timeOut1.go
timeout c1
c2 ok
```

- 方式二：让一个goroutine超时

```go
// timeOut2.go
package main
import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

func timeout(w *sync.WaitGroup, t time.Duration) bool  {
	temp := make(chan int)
	go func () {
		time.Sleep(5 * time.Second)
		defer close(temp)
		w.Wait() // 阻塞并永远等待下去，等待sync.Done()请求，w.Wait()返回后select的第一个分支会被执行。
	}()

	select {
	case <-temp: //w.Wait()返回后select的第一个分支会被执行
		return false
	case <-time.After(t):
		return true
	}
}

	func mian() {
		arguments := os.Args
		if len(arguments) != 2 {
			fmt.Println("Need a time duration!")
			return
		}

		var w sync.WaitGroup
		w.Add(1)

		t, err := strconv.Atoi(arguments[1])
		if err != nil {
			fmt.Println(err)
			return
		}

		duration := time.Duration(int32(t)) * time.Millisecond
		fmt.Printf("Timeout period is %s\n", duration)

		if timeout(&w, duration) {
			fmt.Println("time out!")
		} else {
			fmt.Println("ok!")
		}

	w.Done()
	if timeout(&w, duration){
			fmt.Println("Timed out!")
	} else {
		fmt.Println("OK!")
	}
}

```
```bash

# 超时为10s, 第一个timeout()函数，w.Wait()阻塞了，<-time.After(t):先返回，所以超时。第二个timeout()函数，	w.Done()了，所以这个goroutine里面不需要等待任何东西，所以不会超时。case <-temp:先返回，可以执行。
$ go run timeOut2.go 10000
Timeout period is 10s
Timed out!
OK!

# 超时很短case <-time.After(t):很快返回，所以两个timeout()函数都没有足够时间执行完成。
$ go run timeOut2.go 100
Timeout period is 100ms
Timed out!
Timed out!

```
## go channel重温
- 
## 共享内存和共享变量
## 捕获竞态条件--catching race conditions
## context package
## 资源
## 练习
## 总结

chapter3 go的基本数据类型
- go loop
- go arrays
- go slices
- go maps
- go constants
- go pointers
- 和times工作
- 操作日期


## go loop

## go arrays
```go
// 大小，类型，元素
anArray := [4]int{1,2,4,-4}
len(anArray)
// 多维数组
twoD := [4][4]int{{1,2,3,4},{1,2,3,4},{1,24,45,5},{1,2,3,4}}
treeD := [2][2][2]int{{{1,0},{1,2}},{{3,4},{4,5}}}
```
## go slices
- 底层也是使用array实现的。
- slice传递的是引用，传递内存地址。所以函数里面修改了slice，再函数退出后也是影响原来的slice的。
- 比起array，传递一个大的slice不会触发批量的复制，效率更高。
```go
//和array定义一样的，只是没有声明大小。如果声明大小，返回的是一个array
asliceLiteral := []int{1,2,3,4,5}

//创建定长的slice，make创建还可以自动初始化，但是这个初始化的动作，可能未来版本会被改变。
integer := make([]int,20)

for i := 0; i < len(integer); i++ {
	fmt.Println(integer[i])
}

//
asliceLiteral = nil

//添加元素
integer = append(integer,-5000)

//包含头不包含尾：index=1,index=2的两个元素
integer[1:3]

// 创建一个新的slice, re-slicing，会引起问题的
// 指向同一个内存地址，改变原来的或者reslice的都会互相影响
//这个操作不是复制操作
// 第二个问题：小的reslice只要存在，都不会回收原来的slice
s2 := integer[1:3]

package main
import "fmt"
func main()  {
	s1 := make([]int,5)
	reSlice := s1[1:3]
	fmt.Println(s1)
	fmt.Println(reSlice)

	reSlice[0] = -100
	reSlice[1] = 123456
	fmt.Println(s1)
	fmt.Println(reSlice)
}
```
### slice的自动扩展
- 容量：cap(),当前分配给slice的空间。
- 长度：len(),很数组的长度相同意思。
- 如果一个容量和长度相同的slice，添加一个元素，长度增加1，容量增加一倍。所以，添加一个元素到一个很大的slice，可能会消耗很大的内存。
```go
//lenCap.go
package main
import (
	"fmt"
)
func printSlice(x []int)  {
	for _,number := range x {
		fmt.Print(number," ")
	}
	fmt.Println()
}

func main()  {
	alice := []int{-1,0,4}
	fmt.Printf("aSlice: ")
	printSlice(aSlice)

	fmt.Printf("Cap: %d, Lengths: %d\n", cap(aSlice),len(aSlice))
	aSlice = append(aSlice,-100)
	fmt.Printf("aSlice: ")
	printSlice(aSlice)
	fmt.Printf("Cap: %d, Lengths: %d\n", cap(aSlice),len(aSlice))

	aSlice = append(aSlice,-2)
	aSlice = append(aSlice,-3)
	aSlice = append(aSlice,-4)
	fmt.Printf("aSlice: ")
	printSlice(aSlice)
	fmt.Printf("Cap: %d, Lengths: %d\n", cap(aSlice),len(aSlice))
}
```
```bash
bogon:master-go kate$ go run lenCap.go
aSlice: -1 0 4
Cap: 3, Lengths: 3
aSlice: -1 0 4 -100
Cap: 6, Lengths: 4
bogon:master-go kate$ go run lenCap.go
aSlice: -1 0 4
Cap: 3, Lengths: 3
aSlice: -1 0 4 -100
Cap: 6, Lengths: 4
aSlice: -1 0 4 -100 -2 -3 -4
Cap: 12, Lengths: 7
```
### byte slices
- 类型是byte的slice
- 使用在文件的输入输出

```go
s := make([]byte,5)

```

### copy()函数
- 用途：
	- 只接受slice作为参数
	- 从一个数组创建一个slice
	- 从一个slice创建另外一个slice
	- 有需要注意的问题
	- copy(dst,src): 选择len(dst),len(src)中最小的个数复制
```go
package main
import (
	"fmt"
)
func main() {
	a6 := []int{1,2,3,4,5,6}
	a4 := []int{-1,-2,-3,-4}
	fmt.Println("a6",a6)
	fmt.Println("a4",a4)

	copy(a6,a4)//copy(dst,src)
	fmt.Println("a6",a6)
	fmt.Println("a4",a4)
	fmt.Println()

	array4 := [4]int{-1,-2,-3,-4}
	s6 := []int{1,2,3,4,5,6}
	//array4[0:]是把array转化为了slice
	copy(s6,array4[0:]) //copy(dst,src)
	fmt.Println("array4",array4[0:])
	fmt.Println("a6",a6)
	fmt.Println()

	array5 := [5]int{-1,-2,-3,-4,-5}
	s7 := []int{1,2,3,4,5,6,7}
	copy(array5[0:],s7) //copy(dst,src)
	fmt.Println("array5",array5[0:])
	fmt.Println("s7",s7)
	fmt.Println()
}

```

```bash
bogon:master-go kate$ go run copySlice.go
a6 [1 2 3 4 5 6]
a4 [-1 -2 -3 -4]
a6 [-1 -2 -3 -4 5 6]
a4 [-1 -2 -3 -4]

array4 [-1 -2 -3 -4]
a6 [-1 -2 -3 -4 5 6]

array5 [1 2 3 4 5]
s7 [1 2 3 4 5 6 7]

```

### 高纬度slice

### 另外一个slice的例子

### slices排序使用sort.slice()

## go maps
### 存储一个nil的map
### 什么时候使用

## go constants
### 常量生成器：iota

## go pointers
## 处理时间和日期

- 转换不同的日期，和时间
- 打印不同的日期时间格式
- 字符串转日期时间
```go
package main
import(
	"fmt"
	"time"
)
func main()  {
	// 获取时间从1970年1月00：00：00 UTC开始的时间长度，单位为秒
	fmt.Println("Epoch time:", time.Now().Unix())
	t := time.Now()
	fmt.Println(t,t.Format(time.RFC3339))

	time.Sleep(time.Second)//停止一秒，
	t1 := time.Now()
	// 计算时间差
	fmt.Println("Time difference:", t1.Sub(t))

	formatT := t.Format("01 January 2006")
	fmt.Println(formatT)
	loc,_ := time.LoadLocation("Europe/Paris")
	londonTime := t.In(loc)
	fmt.Println("Paris: ", londonTime)
}
```
### 处理times日期
- time package
- time.Parse() - 解析字符串里面的时间
- 15： 解析小时
- 04： 解析分钟
- 05： 解析秒
- 11： 解析月份
- Jan： 解析月份， January: 也是解析月份
- 2006： 解析年份
- 02: 解析日期
- Monday：解析星期，Mon： 也是解析星期
```go
package main
import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main()  {
	var myDate string
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s string\n",filepath.Base(os.Args[0]))
		return
	}
	myDate = os.Args[1]
	// 定义字符串的日期格式，并传入真正的字符串
	d,err := time.Parse("2006-January-02", myDate)
	if err == nil {
		fmt.Println("Full:", d)
		fmt.Println("Date: ",d.Day(),d.Month(),d.Year())
	} else {
		fmt.Println(err)
	}
}

````
### 处理时间日期格式
```go

```
## 资源附件
## 练习
## 总结

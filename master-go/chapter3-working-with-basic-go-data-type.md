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
- 如果你的代码里面用到了很多这种高纬度的slice那么久重新考虑设计
```go
s1 := make([][]int,4)
```
### 另外一个slice的例子
// todo



### slices排序使用sort.slice()
// todo

## go maps
- 可以使用make()函数创建
- map是一个hash table的引用。
- map的key的数据类型必须是可以比较的，这样编译器才能区分不同的key
- 不要使用floating-point,不同机器存在精确度不一样。
```go
// key：string value:int
iMap = make(map[string]int)

anotherMap := map[string]int {
	"k1": 12
	"k2": 13
}

//
anotherMap["k1"]

delete(anotherMap,"k1")
// usingMaps.go
package main

import (
	"fmt"
)
func main()  {
	iMap := make(map[string]int)
	iMap["k1"] = 12
	iMap["k1"] = 13
	fmt.Println("iMap",iMap)

	anotherMap := map[string]int {
		"k1":12,
		"k2":13,
	}

	fmt.Println("anotherMap",anotherMap)
	delete(anotherMap,"k1")
	delete(anotherMap,"k1")
	delete(anotherMap,"k1")
	fmt.Println("anotherMap",anotherMap)

	value,ok := iMap["doesItExist"]
	fmt.Println("value: " ,value) // 如果这个key不存在，返回的value也是0，就很难判断了。
	fmt.Println("ok: " ,ok)// 使用这个参数判断这个key是否存在，存在：ok == true
	if ok {
		fmt.Println("Exists!")
	} else {
		fmt.Println("Does NOT exist!")
	}

	for key,value := range iMap{ // map的元素返回顺序是随机的
		fmt.Println(key,value)
	}

}

```
```bash
bogon:master-go kate$ go run usingMap.go
iMap map[k1:13]
anotherMap map[k1:12 k2:13]
anotherMap map[k2:13]
a:  0
ok:  false
Does NOT exist!
k1 13

```
### 存储一个key到nil的map
- 但是查询，删除，获取长度，range，应用到一个nil的map是不会报错的
```go
// 这样是没有问题的
aMap := map[string]int{}
aMap["test"] = 1

// 这样是有问题的, nil的map是不能赋值的
aMap := map[string]int{}
aMap = nil
fmt.Println(aMap)
aMap["test"] = 1

```
```bash
$ go run failMap.go
map[]
panic: assignment to entry in nil map
````

### 什么时候使用
- map的比slice和array灵活，处理上开销相对大，但是内置的go结果还是很快的。所以如果觉得有需要还是推荐使用

## go constants
- 程序执行过程中，不能改变他们的值。是在编译过程中定义的。不是运行过程中。
- const 关键字
- 最好把所有的常量写到一个package里面
````go
const HEIGHT = 200

const (
		C1 = "C1C1C1"
		C2 = "C2C2C3"
		C3 = "C3C3C3"
)

// 一样的意思
s1 := "My String"
var s2 = "My String"
var s3 string = "My String"

// 有类型的使用的时候要注意类型限制
const s1 = "My String"
const s2 string = "My String"

// 例子
const s1 = 123
const s2 float64 = 123
var v1 float32 = s1 * 12 // 没有问题
var v2 float32 = s2 * 12 // 类型不同，报错
````


### 常量生成器：iota

```go
// constants.go
package main
import (
	"fmt"
)

// 定义两个新的类型
type Digit int
type Power2 int

const PI = 3.1415926

const (
	C1 = "C1C1C1"
	C2 = "C2C2C3"
	C3 = "C3C3C3"
)

func main(){
	const s1 = 123
	var v1 float32 = s1 * 12
	fmt.Println(v1)
	fmt.Println(PI)


	const(
		Zero Digit = iota
		One
		Two
		Three
		Four
	)
// 上面的常量定义方式等价于
// const(
// 	Zero = 0
// 	One = 1
// 	Two = 2
// 	Three = 3
// 	Four = 4
// )


	fmt.Println(One)
	fmt.Println(Two)

// iota永远都是增长的，_是忽略掉不使用的值6
	const(
		p2_0 Power2 = 1 << iota // iota = 0
		_
		p2_2 // iota = 2； 1 << 2 ; 00000100
		_
		p2_4 // iota = 4； 1 << 4 ; 00010000
		_
		p2_6 // iota = 6； 1 << 4 ; 01000000
		_
	)

	fmt.Println("2^0: ",p2_0)
	fmt.Println("2^2: ",p2_2)
	fmt.Println("2^4: ",p2_4)
	fmt.Println("2^6: ",p2_6)

}

```

```bash
bogon:master-go kate$ go run constants.go
1476
3.1415926
1
2
2^0:  1
2^2:  4
2^4:  16
2^6:  64
```

## go pointers
- 只有必要时使用
- 为了兼容一下原来用ponter实现的语言
- * ：获取指针的值，叫做指针解析
- & ： 获取一个不是指针的变量的内存地址

```go
//指针作为参数的函数
func getPointer(n *int)  {

}

// 指针作为返回值的函数
func returnPointer(n int) *int {

}

package main
import (
	"fmt"
)

// 指针参数
func getPointer(n *int)  {
 *n = *n * *n
}

//指针返回值
func returnPointer(n int) *int  {
	v := n * n
	return &v
}

func main(){
	i := -10
	j := 25

	pI := &i
	pJ := &j

	fmt.Println("pI memory:", pI)
	fmt.Println("pJ memory:", pJ)
	fmt.Println("pI value:", *pI)
	fmt.Println("pJ value:", *pJ)

	*pI = 123456
	*pI--
	fmt.Println("i: ",i)

	getPointer(pJ)
	fmt.Println("j:",j)
	k := returnPointer(12)
	fmt.Println(*k)
	fmt.Println(k)

}
````
```bash
bogon:master-go kate$ go run pointer.go
pI memory: 0xc00006e008
pJ memory: 0xc00006e010
pI value: -10
pJ value: 25
i:  123455
j: 625
144
0xc00006e048

```
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
package main
import (
	"fmt"
	"regexp"
	"time"
)
func main(){
	logs := []string{"127.0.0.1 - - [16/Nov/2017:10:49:46 +0200]325504",
		"127.0.0.1 - - [16/Nov/2017:10:16:41 +0200] GET CVEN HTTP/1.1 200 12531 - Mozilla/5.0 AppleWebKit/537.36",
    "127.0.0.1 200 9412 - - [12/Nov/2017:06:26:05 +0200] GET http://www.mtsoukalos.eu/taxonomy/term/47 1507",
    "[12/Nov/2017:16:27:21 +0300]",
    "[12/Nov/2017:20:88:21 +0200]",
    "[12/Nov/2017:20:21 +0200]",
    }

		for _, logEntry := range logs {
      r := regexp.MustCompile(".*\[(\d\d\/\w+/\d\d\d\d:\d\d:\d\d:\d\d.*)\].*")
			if r.MatchString(logEntry) {
          match := r.FindStringSubmatch(logEntry)
					dt, err := time.Parse("02/Jan/2006:15:04:05 -0700",
					 match[1])

					 if err == nil {
                newFormat := dt.Format(time.RFC850)
                fmt.Println(newFormat)
            } else {
                fmt.Println("Not a valid date time format!")
            }

				} else {
            fmt.Println("Not a match!")
        }
	}

```

## 资源附件
## 练习
## 总结

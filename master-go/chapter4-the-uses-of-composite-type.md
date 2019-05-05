# chapter4 使用composite类型
- 之前学了：arrays, slices，maps，points, canstants，for loop, range关键字，times和dates的使用。
- 这章学习高级特性：tuples（元组）， strings，switch命令，struct关键字定义结构，定义正则表达式和实现模式匹配在go里面，
- 实现一个简单的key-value存储
- go strings, runes，byte slices, string literals
- 正则表达式
- 模式匹配
- switch
- strings包函数

## 关于composite类型
- 元组：允许方程返回多个值，而不用把他们包含在一个结构中

## 结构
- 能存储不通类型的数据的结构
```go

type aStructure struct {
  person string
  height int
  weight int
}

// 创建一个该结构的变量
var s1 aStructure

// 获取一个属性
s1.person

// 赋值, 不需要给每个属性初始化
p1 := aStructure{weight:12, height:-2}

```
- 一般structures定义在main()函数外面，这样可以全局使用。
- 除非你想这个类型只被当前域使用，其他地方不能使用。
- 结构的属性的顺序是很重要的，如果两个结果属性的顺序不一致，也是两个不同结构。
```go
// 创建代码文件：structures.go
package main
import(
  "fmt"
)
func main(){
  // 在这个函数内部定义也不会报错，但是你必须有一个好的理由
  type XYZ struct {
    X int
    Y int
    Z int
  }
  var s1 XYZ
  fmt.Println(s1.Y,s1.Z)
  p1 := XYZ{23,12,-2}
  P2 := XYZ{Z:12,Y:13}
  fmt.Println(p1)
  fmt.Println(p2)

  // 创建一个该结构的数组
  pSlice := [4]XYZ{}
  pSlice[2] = p1
  pSlice[0] = p2
  fmt.Println(pSlice)
  // 结构已经复制到数组，所以后面改变原来的结构不会影响数据里面的结构对象
  p2 = XYZ{1,2,3}
  fmt.Println(pSlice)
}
```

```bash
go run structures.goa
# 0 0
# {23,12,-2}
# {0,13,12}
# [{0,13,12}{0,0,0}{23,12,-2}{0,0,0}]
# [{0,13,12}{0,0,0}{23,12,-2}{0,0,0}]
```
### 结构指针
```go
// 实例代码pointerStruct.go
package main
import(
  "fmt"
)
type myStructure struct {
  Name string
  Surname string
  Height int32
}
// 一个好的习惯，在一个函数里面初始化一个结构。习惯命名：NewStruct()
// c/c++ 更习惯于从函数返回一个变量的地址，
func createStruct(n,s string, h int32) *myStructure {
  if h > 300 {
    h = 0
  }
  return &myStructure{n,s,h}
}
// 不适用指针，或者习惯命名：NewStructurePointer , NewStructure
func retStructure(n,s string, h int32) myStructure{
  if h > 300 {
    h = 0
  }
  return myStructure{n,s,h}
}

func main() {
  s1 := createStruct("Mihalis","Tsoukalos",123)
  s2 := retStructure("Mihalis","Tsoukalos",123)
  fmt.Println((*s1).Name) //从指针中获取变量
  fmt.Println(s2.Name)
  fmt.Println(s1)
  fmt.Println(s2)
}
```

```bash
go run pointerStruct.go
# Mihalis
# Mihalis
# &{Mihalis Tsoukalos 123}
# {Mihalis Tsoukalos 123}
```
### 使用new关键字
- new返回一个指针，就是一个内存地址
- make：返回的对象是已经初始化的，而且不是一个指针，只能用于maps,channels,slices

```go
// 获取了这个结构的分配的内存地址，没有初始化
pS := new(aStructure)

// 使用new创建一个slice，并指向nil
sP := new([]aStructure)

```

## tuples数据结构(元组)，一个函数返回多个变量
```go
min,_ := strconv.ParseFloat(arguaments[1],64)
```

```go
// 示例：tuples.go
package main

import(
  "fmt"
)
// 第六章，会学到分配一个名字给函数的返回值
// 返回多个参数，
func retThree(x int)(int, int, int)  {
  return 2*x, x*x, -x
}

func main()  {
  // 元祖赋值，可以使用_忽略你不使用的参数
  fmt.Println(retThree(10))
  n1,n2,n3 := retThree(20)
  fmt.Println(n1,n2,n3)

  // 其他元祖的使用
  // 交换位置
  n1,n2 = n2,n1
  fmt.Println(n1,n2,n3)

  //表达式计算
  x1, x2, x3 := n1*2, n1*n1, -n1
  fmt.Println(x1,x2,x3)

}
// 执行：go run tuples.go
```
## 正则表达式和模式匹配
- regexp package

### 理论
### 简单例子
- 从一行文本中获取一列
- 一行一行读取一个文本
- 两个参数：列序号，文本文件路径
```go
// seletctColumn.go
package main
import (
  "bufio"
  "fmt"
  "io"
  "os"
  "strconv"
  "strings"
)
func main()  {
  arguaments := os.Args
  if len(arguaments) < 2 {
    fmt.Printf("usage: selectColumn column <file> | <file2> [...fileN]\n")
    os.Exit(1)
  }
  temp, err := strconv.Atoi(arguaments[1])
  if err != nil {
    fmt.Println("Column value is not an integer:", temp)
    return
  }
  column := temp
  if column < 0 {
    fmt.Println("Invalid Column number!")
    os.Exit(1)
  }

  for _,filename := range arguaments[2:] {
    fmt.Println("\t\t",filename)
    f,err := os.Open(filename)
    if err != nil {
      fmt.Printf("error opening file %s\n", err)
      continue
    }
    defer f.Close()

    r := bufio.NewReader()
    for {
      // 按照行读取，返回一个byte slice
      line,err := r.ReadString('\n')

      if err == io.EOF {
        beak
      } else if err != nil {
        fmt.Printf("error reading file %s", err)
      }
      // strings.Fields()函数分割字符串使用一个空格
      data := strings.Fields(line)
      if len(data) >= column {
        fmt.Println(data[column-1])
      }
    }
  }

}


```
```bash
go run selectColumn.go 15 /tmp/swtag.log /tmp/adobegc.log
```
### 高级例子
- 匹配一个日期和时间字符，在日志里面
- 改变日志文件的日期和时间格式
- 也要一行行读取文本文件
- 使用了两个正则表达式
```go
package main
import (
  "fmt"
  "os"
  "bufio"
  "io"
  "regexp"
  "strings"
  "time"
)
  func main()  {
    arguaments := os.Args
    if len(arguaments) == 1 {
      fmt.Println("Please provide one text file to process!")
      os.Exit(1)
    }
    filename := arguaments[1]
    f,err := os.Open(filename)
    if err != nil {
      fmt.Printf("error opening file %s", err)
      os.Exit(1)
    }
    defer f.Close()

    notAMatch := 0
    r := bufio.NewReader(f)
    for {
      line,err := r.ReadString('\n')
      if err == io.EOF {
        break
      } else if err != nil {
        fmt.Printf("error reading file %s", err)
      }

      // 匹配第一种日期格式
      r1 := regexp.MustCompile('.*\[(\d\d\/w+\d\d\d\d:\d\d:\d\d:\d\d.*)\] .*')
      if r1.MatchString(line) {
        match := r1.FindStringSubmatch(line)
        d1, err := time.Parse("02/Jan/2006:15:04:05 -0700",match[1])
        if err == nil {
          newFormat := d1.Format(time.Stamp)
          fmt.Print(strings.Replace(line,match[1],newFormat,1))
        } else {
          notAMatch++
        }
        continue
      }
      // 匹配第二种日期格式
      r2 := regexp.MustCompile('.*\[(\w+\-\d\d-\d\d:\d\d:\d\d:\d\d.*)\].*')
      if r2.MarchString(line) {
        match := r2.FindStringSubmatch(line)
        d1,err := time.Parse("Jan-02-06:15:04:05 -0700", match[1])
        if err == nil {
          newFormat := d1.Format(time.Stamp)
          fmt.Print(strings.Replace(line,match[1],newFormat,1))
        } else {
          notAMatch++
        }
        continue
      }    
    }
    fmt.Println(notAMatch,"lines did not match!")
  }

// go run changDT.go logEntries.txt

```
### 匹配IPv4地址
- 8位二进制数字
- 每位0-255
```go
// findIPv4.go
package main
import(
  "bufio"
  "fmt"
  "io"
  "net"
  "os"
  "path/filepath"
  "regexp"
)
// 可以修改这个方法，获取你需要监控的IP
func findIP(input string) string  {
  //匹配一个部分，
  // 25开头后面一位可以是：0，1，2，3，4，5
  // 或 2开头0，1，2，3，4，最后一位：0，1，2，3，4，5，6，7，8，9
  // 或 1开头：后面[0-9][0-9]
  // 或 第一位：[1-9],第二位[0-9]
  partIP := "25([0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])"
  //
  grammar := partIP + "\\." + partIP + "\\." + partIP + "\\." + partIP
  matchMe := regexp.MustCompile(grammar)
  return matchMe.FindString(input)

}
func main()  {
  arguments := os.Args
  if len(arguments) < 2 {
    fmt.Printf("usage: %s logFile\n", filepath.Base(arguments[0]))
    os.Exit(1)
  }

  for _,filename := range arguments[1:] {
    f, err := os.Open(filename)
    if err != nil {
      fmt.Printf("error opening file %s\n", err)
      os.Exit(-1)
    }
    defer f.Close()

    r := bufio.NewReader(f)
    for {
      line,err := r.ReadString("\n")
      if err == io.EOF {
        break
      } else if err != nil{
         fmt.Printf("error reading file %s", err)
         break
      }
      // 每一行都调用一次IP正则匹配
      ip := findIP(line)
      // 二次验证是否是ipv4的地址，如果是就打印出来
      trial := net.ParseIP(ip)
      if trial.To4() == nil {
        continue
      } else {
        fmt.Println(ip)
      }
    }
  }
}

// go run findIPv4.go /tmp/auth.log

// ipv4 地址出现次数最多的会出现在前面
// go run findIPv4.go /tmp/auth.log.1 /tmp/auth.log | sort -rn | uniq -c | sort -rn

```
## Strings
### 什么是rune
- rune
- character
- string
- 字符串常量（string literal）

### unicode package
### strings package

## switch命令
## 计算π使用高的精确度
## 开发一个key/value存储使用Go
## 另外的资源
## 练习
## 总结

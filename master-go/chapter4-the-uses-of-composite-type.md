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

- rune 特殊符号
- character 字符
- string 字符串
- 字符串常量（string literal）
```go
//字符串常量
const sLiteral = "\x99\x42"
s2 := "k"
// 获取字符长度len()

// strings.go
package main
import (
  "fmt"
)
func main()  {
  const sLiteral = "\x99\x42\x32\x55\x50\x35\x23\x50\x29\x9c"
  fmt.Println(sLiteral) //
  fmt.Println("x: %x\n",sLiteral) // 返回\xAB的AB部分
  fmt.Printf("sLiteral length: %d\n", len(sLiteral))

  for i := i<len(sLiteral); i++ {
    fmt.Printf("%x",sLiteral[i]) //打印字符的Unicode的编码
  }
  fmt.Println()
  fmt.Printf("q: %q\n",sLiteral)
  fmt.Printf("+q: %+q\n",sLiteral)
  fmt.Printf("x: % x\n",sLiteral) //打印字符的Unicode的编码，每个字符之间用空格隔开

  // 打印字符串
  fmt.Printf("s: As a string: %s\n",sLiteral)

  s2 := "€£³"
  // 遍历Unicode编码。
  for x,y := range s2 {
    // %#U 会打印字符串使用U+0058格式
    fmt.Printf("%#U starts at byte position %d\n",y,x)
  }
  // 长度大于字符数量
  fmt.Printf("s2 length: %d\n",len(s2))

  const s3 = "ab12AB"
  fmt.Println("s3:",s3) //打印字符串
  fmt.Printf("x: % x\n", s3) //打印Unicode编码

  fmt.Printf("s3 length: %d\n", len(s3))

  for i := 0; i<len(s3); i++ {
    fmt.Printf("%x ",s3[i]) //打印Unicode编码
  }
  fmt.Println()

}
```

```bash
$ go run strings.go
�B2UP5#P)�
x: 9942325550352350299c
sLiteral length: 10
99 42 32 55 50 35 23 50 29 9c
q: "\x99B2UP5#P)\x9c"
+q: "\x99B2UP5#P)\x9c"
 x: 99 42 32 55 50 35 23 50 29 9c
s: As a string: �B2UP5#P)�
U+20AC '€' starts at byte position 0
U+00A3 '£' starts at byte position 3
U+00B3 '³' starts at byte position 5
s2 length: 7
s3: ab12AB
x: 61 62 31 32 41 42
s3 length: 6
61 62 31 32 41 42

```
- 一般常用fmt.Println(), fmt.Printf()
- 如果不是欧洲或者美国，这里解释的字符编码就有点作用了
### 什么是rune
- 特殊字符及时一个int32值，是一个go类型用来表示一个Unicode编码点，Unicode编码点或者编码位置，是一个数值，用来代表一个Unicode字符。也有一个歧义，例如提供格式信息。
- 可以把string当作一个runes集合
- 一个rune字面量就是一个字符在单引号内。
- 背后，一个rune字面量关联着一个Unicode code point
- 例子
```go
//runes.go
package main
import (
  "fmt"
)
func main()  {
  const r1 = '€' // 欧洲符合，不属于ASCII字符表，
  fmt.Println("(int32) r1:", r1) // 打印32位int值
  fmt.Printf("(HEX) r1: %x\n",r1) // 打印16进制值
  fmt.Printf("(as a string) r1: %s\n", r1) // 打印字符串
  fmt.Printf("(as a character) r1: %c\n",r1) //打印字符

  fmt.Println("A string is a collection of runes:", []byte("Mihalis")) //打印Unicode编码
  aString := []byte("Mihalis")
  for x,y := range aString {
    fmt.Println(x,y)//打印序号，和Unicode编码
    fmt.Printf("Char: %c\n", aString[x]) //打印字符
  }
  fmt.Printf("%s\n",aString) //打印字符
}
```

```bash
$ go run runes.go
(int32) r1: 8364
(HEX) r1: 20ac
(as a String) r1: %!s(int32=8364)
(as a character) r1: €
A string is a collection of runes: [77 105 104 97 108 105 115]
0 77
Char: M
1 105
Char: i
2 104
Char: h
3 97
Char: a
4 108
Char: l
5 105
Char: i
6 115
Char: s
Mihalis

```
- 获得一个illegal rune literal错误信息例子
```bash
$ cat a.go
package main
import (
    'fmt'
)
func main(){
}

$ go run a.go
package main:a.go:4:2: illegal rune literal

```
### unicode package
- 包含许多有用的方法
- unicode.IsPrint(): 确认一个部分的string是否可以打印使用runes
```go
// unicode.go
package main
import(
  "fmt"
  "unicode"
)
func main()  {
  const sL = "\x99\x00ab\x50\x00\x23\x50\x29\29c"
for i := 0; i < len(sL); i++ {
  if unicode.IsPrint(rune(sL[i])) {
    fmt.Printf("%c\n",sL[i]) //打印字符
  } else {
    fmt.Println("Not printable!")
  }
}

}

```
```bash
$ go run unicode.go
Not printable!
Not printable!
a
b
P
Not printable!
#
 P
)
Not printable!

```

### strings package
- 处理UTF-8字符串
- 有很多有用的函数
- 例子中使用的文件输入输出将在第8章展示
```go
// useString.go
package main
import (
  "fmt"
  s "string" //包使用别名， s.FunctionName(), string.FunctionName()不能再使用
  "unicode"
)

var f = fmt.Printf() // 方法使用别名，一般不使用
func main() {
  upper := s.ToUpper("Hello there!")
  f("To Upper: %s\n", upper)
  f("To Lower: %s\n", s.ToLower("Hello THERE"))
  f("%s\n", s.Title("tHis will be a title!"))

 // 觉得两个字符串是否相等，即使实际是不相等的。
  f("EqualFold: %v\n", s.EqualFold("Mihalis","MIHALIS"))
  f("EqualFold: %v\n", s.EqualFold("Mihalis","MIHAli"))

  f("Prefix: %v\n",s.HasPrefix("Mihalis","Mi"))
  f("Prefix: %v\n",s.HasPrefix("Mihalis","mi"))
  f("Prefix: %v\n",s.HasSuffix("Mihalis","is"))
  f("Prefix: %v\n",s.HasSuffix("Mihalis","IS"))

  f("Index: %v\n",s.Index("Mihalis","ha"))
  f("Index: %v\n",s.Index("Mihalis","Ha"))
  f("Count: %v\n",s.Count("Mihalis","i"))
  f("Count: %v\n",s.Count("Mihalis","I"))
  f("Repeat: %s\n",s.Repeat("ab",5))

  f("TrimSpace: %s\n", s.TrimSpace(" \tThis is a line. \n"))
  f("TrimLeft: %s", s.TrimLeft(" \tThis is a\t line. \n", "\n\t "))
  f("TrimRight: %s\n", s.TrimRight(" \tThis is a\t line. \n", "\n\t "))

// 如果相同返回0，否则返回-1，+1
  f("Compare: %v\n", s.Compare("Mihalis", "MIHALIS"))
  f("Compare: %v\n", s.Compare("Mihalis", "Mihalis"))
  f("Compare: %v\n", s.Compare("MIHALIS", "MIHalis"))

//使用空格分割字符串
  f("Fields: %v\n", s.Fields("This is a string!"))
  f("Fields: %v\n", s.Fields("Thisis\na\tstring!"))

// 返回一个string slice， 使用"",按照字符分割
  f("%s\n", s.Split("abcd efg", ""))

//Replace参数：1-处理的字符穿，2-被替换字符串 3-替换字符串 4-最大的替换数量如果为负数，没有限制
  f("%s\n", s.Replace("abcd efg", "", "_", -1))
  f("%s\n", s.Replace("abcd efg", "", "_", 4))
  f("%s\n", s.Replace("abcd efg", "", "_", 2))

  lines := []string{"Line 1", "Line 2", "Line 3"}
  f("Join: %s\n", s.Join(lines, "+++"))

  f("SplitAfter: %s\n", s.SplitAfter("123++432++", "++"))

// 自定义一个函数去过滤出自己想要的字符串
  trimFunction := func(c rune) bool {
      return !unicode.IsLetter(c)
  }
  f("TrimFunc: %s\n", s.TrimFunc("123 abc ABC \t .", trimFunction))

}

```
```bash
$ go run useStrings.go
To Upper: HELLO THERE!
To Lower: hello there
THis WiLL Be A Title!
EqualFold: true
EqualFold: false
Prefix: true
Prefix: false
Suffix: true
Suffix: false
Index: 2
Index: -1
Count: 2
Count: 0
Repeat: ababababab
TrimSpace: This is a line.
TrimLeft: This is a      line.
TrimRight:        This is a    line.
Compare: 1
Compare: 0
Compare: -1
Fields: [This is a string!]
Fields: [Thisis a string!]
[a b c d   e f g]
_a_b_c_d_ _e_f_g
__a_b_c_d efg
_a_bcd efg
Join: Line 1+++Line 2+++Line 3
SplitAfter: [123++ 432++ ]
TrimFunc: abc ABC

```
## switch命令
- 这个命令可以使用正则表达式
```go
switch asString {
case "1":
  fmt.Println("one!")
case "0":
  fmt.Println("zero!")
default:
  fmt.Println("do not care!")
}

switch {
case number < 0:
    fmt.Println("Less than zero!")
case number > 0:
    fmt.Println("Bigger than zero!")
default:
    fmt.Println("Zero!")
}
//switch.go
package main
import(
  "fmt"
  "os"
  "regexp"
  "strconv"
)
func main(){
  arguments := os.Args
  if len(arguments) < 2 {
    fmt.Println("usage: switch number")
    os.Exit(1)
  }

  number, err != strconv.Atoi(arguments[1])
  if err != nil {
    fmt.Println("This value is not an integer: ", number)
  } else {
    switch {
    case number < 0:
        fmt.Println("Less than zero!")
    case number > 0:
        fmt.Println("Bigger than zero!")
    default:
        fmt.Println("Zero!")
    }
  }
  asString := arguments[1]
  switch asString {
  case "5":
    fmt.Println("five!")
  case "0":
    fmt.Println("zero!")
  default:
    fmt.Println("do not care!")
  }
}
```

## 计算π使用高的精确度
## 开发一个key/value存储使用Go

## 另外的资源
## 练习
## 总结

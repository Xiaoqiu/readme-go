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
func retThree(x int)(int, int, int)  {
  return 2*x, x*x, -x
}

```
## 正则表达式和模式匹配
## Strings
## switch命令
## 计算π使用高的精确度
## 开发一个key/value存储使用Go
## 另外的资源
## 练习
## 总结

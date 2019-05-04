# flag package

# the io.reader and io.writer interfaces
- 需要实现Read()方法
- 需要实现Write()方法
## buffered 和 unbuffered 文件输入 ，输出
- 一般比较重要的数据就不用缓冲，不然容易丢失

# the bufio package
- 这个包是缓冲输入输出的
- 这个包还是使用io.Reader 和 io.Writer 对象
- 封装为bufio.Reader 和 bufio.Writer 对象
- 这个包非常常用来读取文件

# 读取文本文件
- 文本文件是Unix系统最多的文件
- 一行行读取，一个个词读取，一个个字符读取
- 一行行读取时最容易的，一个个词读取时最难的
## 一行行读取文本文件
```go
// 文件byLine.go
package main

import(
  "bufio"
  "flag"
  "fmt"
  "io"
  "os"
)

func lineByLine(file string) error {
  var err error

  f,err := os.Open(file)
  if err != nil {
    return err
  }

  defer f.Close()

  r := bufio.NewReader(f)

  for {
    line.err := r.ReadString("\n")

    if err == io.EOF {
      break
    }else if err != nil {
      fmt.Printf("error reading file %s", err)
      break
    }

    fmt.Print(line)
  }
  returen nil
}

func main()  {
  flag.Parse()
  if len(flag.Args() == 0) {
    fmt.Printf("usage: byline <file1> [<file2> ...] \n")
    return
  }

  for _,file := range flag.Args() {
    err := lineByLine(file)
    if err != nil {
      fmt.Println(err)
    }
  }
}

```

- 执行文件
```bash
# 执行命令
go run byLine.go /tmp/swatag.log /tmp/adobegc.log | wc

# 输出 4761 88521 568402

# 验证上面命令输出是否正确
wc /tmp/swatag.log /tmp/adobegc.log

```
## 一个个词读取文本文件
```go
// byWord.go
package main

import(
  "bufio"
  "flag"
  "fmt"
  "io"
  "os"
  "regexp"
)

func wordByWord(file string)error  {
  var err error
  // 打开文件
  f,err := os.Open(file)
  if err != nil {
    return err
  }
  defer f.Close()
  // 创建reader
  r := bufio.NewReader(f)

  for {
    // 用reader读取文件
    line, err := r.ReadString("\n")
    if err == io.EOF {
      break
    } else if err != nil {
      fmt.Printf("error reading file %s",err)
      return err
    }
    // 使用正则表达式在每行里面分割词； 使用空格分割词
    r := regexp.MustCompile("[^\\s]+")
    words := r.FindAllString(line,-1)
    for i:= 0; i<len(words); i++ {
      fmt.Println(words[i])
    }

  }

  return nil
}

func main() {
  flag.Parse()
  if len(flag.Args() == 0) {
    fmt.Printf("usage:byWord <file1> [<file2>...]\n")
    return
  }

  for _,file := range flag.Args() {
    err := wordByWord(filr)
    if err != nil {
      fmt.Println(err)
    }
  }
}

```

```bash
# 运行命令
go run byWord.go /tmp/adobegc.log

# 输出结果
01/08/18
...

# 代码执行锁获取的行和词数量
go run byWord.go /tmp/adobegc.log | wc

# 统计文件的词数量
wc /tmp/adobegc.log
```

## 一个个字符读取文本文件
```go
package main
import(
  "bufio"
  "flag"
  "fmt"
  "io"
  "os"
)

func charByChar(file string) error  {
  var err error
  f,err := os.Open(file)
  if err != nil {
    return err
  }
  defer f.Close()

  r := bufio.NewReader(f)
  for {
    line,err := r.ReadString("\n")
    if err == io.EOF {
      break
    } else if err != nil {
      fmt.Printf("error reading file %s", err)
      return err
    }
    // range返回两个值：一序列号，一个值。但是这个值不是string， 所以要装换一下。
    for _,x := range line {
        fmt.Println(string(x))
      }    
  }
  return nil

  func main()  {
    flag.Parse()
    if len(flag.Args()) == 0 {
      fmt.Printf("usage:byWord <file1> [<file2>...]\n")
      return
    }

    for _,file := reange flag.Args() {
      err := charByChar(file)
      if err != nil {
        fmt.Println(err)
      }
    }

  }
}
```

```bash
# 这个方法可以当做wc命令统计字符数量
go run byCharacter.go /tmp/adobegc.log
0
1
...

```
## 从/dev/random读取
- /dev/random的作用是生成随机数
```go
// devRandom.go
package main
import(
  "encoding/binary" //读取二进制数，需要解码
  "fmt"
  "os"
)
func main(){
  f,err := os.Open("/dev/random")
  defer f.Close()

  if err != nil {
    fmt.Println(err)
    return
  }

  var seed int64
// binary.LittleEndian 计算机的小端
  binary.Read(f,binary.LittleEndian, &seed)
  fmt.Println("Seed:", seed)
}
```
```bash
go run devRandom.go
# 输出一串数字
```

# 从文件读取一定数量的数据
- 特别是读取二进制的数据，你需要读取固定的长度，然后解码
- 使用byte slice数据结果存储
- 输入为一个读取数量(byte slice的大小)，和一个文件类型*os.File
- 输出为你想要读的数据
- 这个代码的用途，可以帮你使用buffer size的大小读取你想要的文件
```go
// readSize.go
package main
import(
  "fmt"
  "io"
  "os"
  "strconv"
)
func readSize(f *os.File, size int)[]byte  {
  buffer := make([]byte, size)
  // io.Reader.Read()方法返回：读取的字节数，和err
  // buffer 是一个slice数据结构
  n,err := f.Read(buffer)   

  if err == io.EOF {
     return nil
  }

  if err != nil {
    fmt.Print(err)
    return nil
  }
  return buffer[0:n]
}

func main()  {
  arguments := os.Args
  if len(arguments) != 3 {
    fmt.Println("<buffer size> <filename>")
    return
  }

//
  bufferSize, err := strconv.Atoi(os.Args[1])

  if err != nil {
    fmt.Println(err)
    return
  }

  file := os.Args[2]
  f,err := os.Open(file)

  f,err != nil {
    fmt.Println(err)
    return
  }

  defer f.close()
// 重复读取文件一定数量数据，直到返回error或者nil
  for {
    readData := readSize(f,bufferSize)
    if readData != nil {
      fmt.Print(string(readData))
    } else {
      break
    }
  }
}
```

```bash
go run readSize.go 1000 /bin/ls |  wc
#输出 80 1032 38688

wc /bin/ls
#输出 80 1032 38688 /bin/ls

```
# 为什么使用二进制格式
- readSize.go是最好的例子用来读取二进制文件
- 例如存储数组20作为字符串存储，你需要两个字节存储ASII字符，一个存2，一个存0
- 使用二进制存储只要一个字节，因为20可以用二进制表示00010100，或者0x14用16进制

# 读取cvs文件
- 有格式的文本文件
- 学习读取一个文本文件包含一个飞机的坐标，每一行包含一对坐标。
- 使用go library Glot, 创建一个plot来存取读取的点。
- Glot使用的是Gunplot, 所以你要安装Gnuplot在Unix系统上
```go
// CSVplot.go
package main
import(
  "encoding/csv"
  "fmt"
  "github.com/Arafatk/glot"
  "os"
  "strconv"
)

func main(){
  if len(os.Args) != 2 {
    fmt.Println("Need a data file!")
    return
  }

  file := os.Args[1]
  //判断文件是否存在
  _,err := os.Stat(file)
  if err != nil {
    fmt.Println("Cannot stat", file)
    return
  }

  f,err := os.Open(file)
  if err != nil {
    fmt.Println("Cannot open", file)
    fmt.Println(err)
    return
  }

  defer f.Close()

  reader := csv.NewReader(f)
  reader.FieldsPerRecord = -1
  allRecords, err := reader.ReaderAll()
  if err != nil {
    fmt.Println(err)
    return
  }

  xP := []float64{}
  yP := []float64{}

// 转化读取的string为numbers放入二维的slice变量points中
  for _,rec := range allRecords {
    x,_ := strconv.ParseFloat(rec[0],64)
    y,_ := strconv.ParseFloat(rec[1],64)
    xP = append(xP,x)
    yP = append(yP,y)
  }

  points := [][]float64{}
  points = append(points,xP)
  points = append(points,yP)

  fmt.Println(points)

  //打印坐标
  dimensions := 2
  persist := true
  debug := false
  plot,_ := glot.NewPlot(dimensions,persist,debug)

  plot.setTitle("Using Glot with CSV data")
  plot.SetXLabel("X-Axis")
  plot.SetYLabel("Y-Axis")
  style := "circle"
  plot.AddPointGroup("Circle",style,points)
  plot.SavePlot("output.png")

}

```
```bash
# 下载依赖库Glot
go get github.com/Arafatk/glot

# 下载成功后，可见，
ls -l ~/go/pkg/darwin_amd64/github.com/Arafatk/glot/
# 输出 total 120 。。。。
# 。。。
#

# 源数据文件格式
cat /tmp/dataFile
# 1,2
# 2,3
# 3,3
# 4,4
# 5,8
# 6,5
# -1,12
# -2,10
# -3,10
#-4,10
# 转为points [[1,2,3,4,5,6,-1,-2,-3,-4][2,3,3,4,8,5,12,10,10,10]]
go run CSVplot.go /tmp/dataFile

```
# 写入文件
- 使用io.Writer接口可以写数据到硬盘文件，
- save.go 使用5种方式写入文件
```go
package main
import (
  "bufio"
  "fmt"
  "io"
  "io/ioutil"
  "os"
)

func main()  {
// slice数据结构
  s := []byte("Data to write\n")
  f1,err := os.Create("f1.txt")
  if err != nil {
    fmt.Println("Cannot create file", err)
    return
    }
  defer f1.Close()
  //使用喜欢的格式写入数据到文件
  fmt.Fprintf(f1,string(s))

  f2,err := os.Create("f2.txt")
  if err != nil {
    fmt.Println("Cannot create file", err)
    return
    }
  defer f2.Close()

  //方式2：
  n,err := f2.WriteString(string(s))
  fmt.Printf("wrote %d bytes\n",n)  

  f3,err := os.Create("f3.txt")
  if err != nil {
    fmt.Println(err)
    return
  }
  w := bufio,NewWriter(f3)
  // 方式3
  n,err = w.WriteString(string(s))
  fmt.Printf("wrote %d bytes\n",n)
  w.Flush()

  // 方式4
  f4 := "f4.txt"
  err = ioutil.WriteFile(f4,s,0644)
  if err != nil {
    fmt.Println(err)
    return
  }

  // 方式5
  f5,err := os.Create("f5.txt")
  if err != nil {
    fmt.Println(err)
    return
  }
  n,err = io.WriteString(f5,string(s))
  if err != nil{
    fmt.Println(err)
    return
  }
  fmt.Printf("wrote %d bytes\n",n)
}

```
```bash
go run save.go

ls -l f?.txt

cat f?.txt

```

# 加载和保存数据到硬盘disk
- 第四章，使用Composite types，示例代码：keyValue.go
- package encoding/gob 序列化和反序列化
- 使用gob format 存储数据，也叫做流编码
- 其他序列化和反序列化包：encoding/xml
- package ncoding/json
- 使用diff 命令现实keyValue.go kvSaveLoad.go不同之处，不显示完整代码。
# string package

# bytes package

# 文件权限

# 处理unix signals

# 编程 unix pipes in go

# 遍历路径树

# 使用eBPF from go

# 关于syscall.PtraceRegs

# User ID 和 group ID

# Additional resources

# 联系

# 总结

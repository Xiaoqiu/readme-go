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

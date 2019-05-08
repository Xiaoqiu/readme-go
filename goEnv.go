package main
import (
  "fmt"
  "runtime"
)

func main()  {
  fmt.Print("you are using: ", runtime.Compiler, " ")
  fmt.Println("on a", runtime.GOARCH, "machine")
  fmt.Println("go version : ",runtime.Version())
  fmt.Println("Number of CPUs : ", runtime.NumCPU())
  fmt.Println("Number of Goroutines: ", runtime.NumGoroutine())
}

package main
import (
  "fmt"
)

func main()  {
  // 这样是没有问题的
  aMap := map[string]int{}
  aMap["test"] = 1

  // 这样是有问题的, nil的map是不能赋值的
  aMap := map[string]int{}
  aMap = nil
  fmt.Println(aMap)
  aMap["test"] = 1

}

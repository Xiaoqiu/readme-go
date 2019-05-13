// create.go
// go run create.go -n 100
package main
import (
	"flag"
	"fmt"
	"time"
)
func main()  {
	// 设置一个默认值
	n := flag.Int("n", 10, "Number of goroutines")
	// 使用flag package 读取传入参数
	flag.Parse()

	count := *n
	fmt.Printf("Going to create %d goroutines.\n" , count)

	for i := 0; i<count; i++ {
		go func (x int)  {
			fmt.Printf("%d ", x)
		}(i)
	}

	//等待goroutine完成后，结束main函数
	//还有其他更加优雅的方式等待goroutine结束，在main()函数之前。
	time.Sleep(time.Second)
	fmt.Println("\nExiting...")

}

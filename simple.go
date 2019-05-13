package main
import (
	"fmt"
	"time"
)
func function()  {
	for i:=0; i<10; i++ {
		fmt.Print(i," ")
	}
	fmt.Println()
}

func main()  {
	//使用一般函数，启动一个goroutine
go function()

	// 使用匿名函数启动一个goroutine，这种方式更实际如果实现小函数
	go func(){
		for i := 10; i<20; i++ {
			fmt.Print(i," ")
		}
	}()

//等待gorouitne结束，打印到输出。再结束main函数
	time.Sleep(1 * time.Second)

}

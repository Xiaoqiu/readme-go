package main
import(
	"flag"
	"fmt"
	"sync"
)
func main()  {
	n := flag.Int("n", 20, "Number of gorouitnes")
	flag.Parse()
	count := *n
	// 情况1： 多个一个 waitGroup.Add(1)
	fmt.Printf("Going to create %d gorouitnes. \n", count)
	var waitGroup sync.WaitGroup

	fmt.Printf("%#v\n", waitGroup)

	for i := 0; i<count; i++ {
		//每个gorouitne启动之前增加计数器的值
		waitGroup.Add(1)

		go func (x int)  {
			//每个goroutine结束之后减少计数器的值
			defer waitGroup.Done()
			fmt.Printf("%d ",x)
		}(i)
	}
	// 情况2： 多个一个 waitGroup.Done()
	
	fmt.Printf("%#v\n", waitGroup)
	waitGroup.Wait()
	fmt.Println("\nExiting....")

}

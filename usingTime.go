package main
import(
	"fmt"
	"time"
)
func main()  {
	// 获取时间从1970年1月00：00：00 UTC开始的时间长度，单位为秒
	fmt.Println("Epoch time:", time.Now().Unix())
	t := time.Now()
	fmt.Println(t,t.Format(time.RFC3339))

	time.Sleep(time.Second)//停止一秒，
	t1 := time.Now()
	// 计算时间差
	fmt.Println("Time difference:", t1.Sub(t))

	formatT := t.Format("2006 Jan 02")
	fmt.Println(formatT)
	loc,_ := time.LoadLocation("Europe/Paris")
	londonTime := t.In(loc)
	fmt.Println("Paris: ", londonTime)
	fmt.Println("Paris Hour: ", londonTime.Hour())
}

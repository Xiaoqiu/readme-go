package main
import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main()  {
	var myDate string
	if len(os.Args) != 2 {
		fmt.Printf("usage: %s string\n",filepath.Base(os.Args[0]))
		return
	}
	myDate = os.Args[1]
	// 定义字符串的日期格式，并传入真正的字符串
	d,err := time.Parse("2006-January-02", myDate)
	if err == nil {
		fmt.Println("Full:", d)
		fmt.Println("Date: ",d.Day(),d.Month(),d.Year())
	} else {
		fmt.Println(err)
	}
}

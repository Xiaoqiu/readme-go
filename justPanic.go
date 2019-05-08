package main
import(
	"os"
	"fmt"
)
func main()  {
	fmt.Println("arg[0] : ", os.Args[0])
	fmt.Println("args len:", len(os.Args))
	if len(os.Args) == 1 { //
		panic("not enough arguments")
	}
	fmt.Println("gooooo!")
}

package main
import "fmt"
func main()  {
	s1 := make([]int,5)
	s1[0] = 0
	s1[1] = 1
	s1[2] = 2
	s1[3] = 3
	s1[4] = 4
	reSlice := s1[1:3]
	fmt.Println(s1)
	fmt.Println(reSlice)

	reSlice[0] = -100
	reSlice[1] = 123456
	fmt.Println(s1)
	fmt.Println(reSlice)
}

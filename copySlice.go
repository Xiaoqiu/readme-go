package main
import (
	"fmt"
)
func main() {
	a6 := []int{1,2,3,4,5,6}
	a4 := []int{-1,-2,-3,-4}
	fmt.Println("a6",a6)
	fmt.Println("a4",a4)

	copy(a6,a4)//copy(dst,src)
	fmt.Println("a6",a6)
	fmt.Println("a4",a4)
	fmt.Println()

	array4 := [4]int{-1,-2,-3,-4}
	s6 := []int{1,2,3,4,5,6}
	//array4[0:]是把array转化为了slice
	copy(s6,array4[0:]) //copy(dst,src)
	fmt.Println("array4",array4[0:])
	fmt.Println("a6",a6)
	fmt.Println()

	array5 := [5]int{-1,-2,-3,-4,-5}
	s7 := []int{1,2,3,4,5,6,7}
	copy(array5[0:],s7) //copy(dst,src)
	fmt.Println("array5",array5[0:])
	fmt.Println("s7",s7)
	fmt.Println()
}

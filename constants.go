package main
import (
	"fmt"
)

// 定义两个新的类型
type Digit int
type Power2 int

const PI = 3.1415926

const (
	C1 = "C1C1C1"
	C2 = "C2C2C3"
	C3 = "C3C3C3"
)

func main(){
	const s1 = 123
	var v1 float32 = s1 * 12
	fmt.Println(v1)
	fmt.Println(PI)


	const(
		Zero Digit = iota
		One
		Two
		Three
		Four
	)
// 上面的常量定义方式等价于
// const(
// 	Zero = 0
// 	One = 1
// 	Two = 2
// 	Three = 3
// 	Four = 4
// )


	fmt.Println(One)
	fmt.Println(Two)

// iota永远都是增长的，_是忽略掉不使用的值6
	const(
		p2_0 Power2 = 1 << iota // iota = 0
		_
		p2_2 // iota = 2； 1 << 2 ; 00000100
		_
		p2_4 // iota = 4； 1 << 4 ; 00010000
		_
		p2_6 // iota = 6； 1 << 4 ; 01000000
		_
	)

	fmt.Println("2^0: ",p2_0)
	fmt.Println("2^2: ",p2_2)
	fmt.Println("2^4: ",p2_4)
	fmt.Println("2^6: ",p2_6)

}

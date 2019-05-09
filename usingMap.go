package main

import (
	"fmt"
)
func main()  {
	iMap := make(map[string]int)
	iMap["k1"] = 12
	iMap["k1"] = 13
	fmt.Println("iMap",iMap)

	anotherMap := map[string]int {
		"k1":12,
		"k2":13,
	}

	fmt.Println("anotherMap",anotherMap)
	delete(anotherMap,"k1")
	delete(anotherMap,"k1")
	delete(anotherMap,"k1")
	fmt.Println("anotherMap",anotherMap)

	value,ok := iMap["doesItExist"]
	fmt.Println("value: " ,value) // 如果这个key不存在，返回的value也是0，就很难判断了。
	fmt.Println("ok: " ,ok)// 使用这个参数判断这个key是否存在，存在：ok == true
	if ok {
		fmt.Println("Exists!")
	} else {
		fmt.Println("Does NOT exist!")
	}

	for key,value := range iMap {
		fmt.Println(key,value)
	}

}

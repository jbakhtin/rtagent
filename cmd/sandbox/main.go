package main

import "fmt"

func main() {
	testMap := make(map[string]string, 2)

	testMap["test1"] = "Test1"
	testMap["test2"] = "Test2 "

	copiedMap := make(map[string]string, len(testMap))

	for k, v := range testMap {
		copiedMap[k] = v
	}

	copiedMap["test2"] = "Test3"

	fmt.Println(testMap)
	fmt.Println(copiedMap)
}

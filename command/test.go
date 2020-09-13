package main

import (
	"fmt"
	"github.com/tongsq/gin-example/component/logger"
	"os"
)

func main() {
	logger.Error("test")
	dir, _ := os.Getwd()
	println(dir)
	var arr []string
	arr = append(arr, "a")
	arr = append(arr, "b")
	arr = append(arr, "c")
	var arr2 []string = arr[1 : len(arr)-1]
	fmt.Println(arr2)
}

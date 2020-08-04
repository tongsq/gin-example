package main

import (
	"github.com/tongsq/gin-example/component/logger"
	"os"
)

func main() {
	logger.Error("test")
	dir, _ := os.Getwd()
	println(dir)

}

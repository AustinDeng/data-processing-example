package main

import (
	"log"
	"os"
	// 调用 matchers 的 init 函数
	_ "data-processing-example/matchers"
	"data-processing-example/search"
)

func init() {
	// 设置日志输出到标准输出
	log.SetOutput(os.Stdout)
}

func main() {
	// 使用特定的项进行搜索
	search.Run("president")
}

package search

import (
	"fmt"
	"log"
)

// 保存搜索结果的结构
type Result struct {
	Field   string
	Content string
}

// 定义了要实现的新搜索类型的行为
type Matcher interface {
	Search(feed *Feed, searchTerm string) ([]*Result, error)
}

// 为每个数据源单独启动 goroutine 来执行这个函数并发地执行搜素
func Match(matcher Matcher, feed *Feed, searchTerm string, results chan<- *Result) {
	// 对特定的匹配器执行搜素
	searchResults, err := matcher.Search(feed, searchTerm)
	if err != nil {
		log.Println(err)
		return
	}

	// 将结果写入通道
	for _, result := range searchResults {
		results <- result
	}
}

// 从每个单独的 goroutine 中接受到结果后,在终端中输出
func Display(results chan *Result) {
	// 通道会一直被阻塞,知道有结果写入
	// 一旦通道被关闭, for 循环就会终止
	for result := range results {
		fmt.Printf("%s:\n%s\n\n", result.Field, result.Content)
	}
}

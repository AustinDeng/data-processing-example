package search

import (
	"log"
	"sync"
)

// 注册用于搜索的匹配器类型
var matchers = make(map[string]Matcher)

// 搜索逻辑
func Run(searchTerm string) {
	// 获取需要搜索的数据源列表
	feeds, err := RetrieveFeeds()
	if err != nil {
		log.Fatal(err)
	}

	// 创建一个无缓冲的通道,接受匹配后的结果
	results := make(chan *Result)

	// 构造一个 waitGroup, 以便处理所有的数据源
	var waitGroup sync.WaitGroup

	// 设置需要等待处理每个数据源的 goroutine 数量
	waitGroup.Add(len(feeds))

	// 为每个数据源启动一个 goroutine 来查找结果
	for _, feed := range feeds {
		// 获取一个匹配器用于查找
		matcher, exists := matchers[feed.Type]
		if !exists {
			matcher = matchers["default"]
		}

		// 启动一个 goroutine 来执行搜索
		go func(matcher Matcher, feed *Feed) {
			Match(matcher, feed, searchTerm, results)
			waitGroup.Done()
		}(matcher, feed)
	}

	// 启动一个 goroutine 来监控是否所有任务都已经完成
	go func() {
		// 等待所有任务完成，这里会导致 goroutine 阻塞，直到 WaitGroup 内部计数为 0
		waitGroup.Wait()

		// 关闭通道，通知 Display 函数可以退出程序了
		close(results)
	}()

	// 启动函数，显示返回结果，并且在最后一个结果显示完成后返回
	Display(results)
}

// 注册匹配器,提供给后面的程序使用
func Register(feedType string, matcher Matcher) {
	if _, exists := matchers[feedType]; exists {
		log.Fatalln(feedType, "Matcher already registered")
	}

	log.Println("Register", feedType, "matcher")
	matchers[feedType] = matcher
}

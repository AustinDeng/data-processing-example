# data-processing-example 

这个程序从不同的数据源拉取数据，将数据内容与一组搜索项作对比,然后将匹配的内容显示在终端窗口。

首先程序会读取文本文件，进行网络调用，解码 XML 和 JSON 成为结构化类型数据，并且利用 Go 语言的并发机制保证这些操作的速度。

## 程序流程图
![程序流程架构图](https://i.loli.net/2018/10/27/5bd3e35e053d0.png)

## 安装

    // 获取项目源码
    $ cd $GOPATH/src
    $ git clone git@github.com:AustinDeng/data-processing-example.git
    

## 运行

    // 进入目录
    $ cd $GOPATH/src/data-processing-example
    
    // 执行命令
    $ go run main.go

## 接口实现

        // 定义了要实现的新搜索类型的行为
        type Matcher interface {
        	Search(feed *Feed, searchTerm string) ([]*Result, error)
        }
        
        // 实现默认匹配器的行为
        func (m defaultMatcher) Search(feed *Feed, searchTerm string) ([]*Result, error) {
        	return nil, nil
        }
        
        // 实现 rss 匹配器的行为
        func (m rssMatcher) Search(feed *search.Feed, searchTerm string) ([]*search.Result, error) {
    	var results []*search.Result
    
    	log.Printf("Search Feed Type[%s] Site[%s] For Uri[%s]\n", feed.Type, feed.Name, feed.URI)
    
    	// 获取要搜索的数据
    	document, err := m.retrieve(feed)
    	if err != nil {
    		return nil, err
    	}
    
    	for _, channelItem := range document.Channel.Item {
    		// 检查标题部分是否包含搜索项
    		matched, err := regexp.MatchString(searchTerm, channelItem.Title)
    		if err != nil {
    			return nil, err
    		}
    
    		// 如果找到匹配项,将其作为结果保存下来
    		if matched {
    			results = append(results, &search.Result{
    				Field:   "Title",
    				Content: channelItem.Title,
    			})
    		}
    
    		// 检查描述部分是否包含搜索项
    		matched, err = regexp.MatchString(searchTerm, channelItem.Description)
    		if err != nil {
    			return nil, err
    		}
    
    		// 如果找到匹配项,将其作为结果保存下来
    		if matched {
    			results = append(results, &search.Result{
    				Field:   "Description",
    				Content: channelItem.Description,
    			})
    		}
    	}
    
    	    return results, nil
        }


## 搜索逻辑

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

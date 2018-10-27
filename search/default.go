package search

// 默认匹配器
type defaultMatcher struct{}

// 将默认匹配器注册到程序中
func init() {
	var matcher defaultMatcher
	Register("default", matcher)
}

// 实现默认匹配器的行为
func (m defaultMatcher) Search(feed *Feed, searchTerm string) ([]*Result, error) {
	return nil, nil
}

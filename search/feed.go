package search

import (
	"encoding/json"
	"os"
)

// Feed 包含所需处理的数据源信息
type Feed struct {
	Name string `json:"site"`
	URI  string `json:"link"`
	Type string `json:"type"`
}

const dataFile = "data/data.json"

// 读取并反序列化源数据文件
func RetrieveFeeds() ([]*Feed, error) {
	// 打开文件
	file, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}
	// 当函数返回时,关闭文件
	defer file.Close()

	// 将文件解码到一个切片里, 里面的每一个元素都是一个指向一个 Feed 类型值的指针
	var feeds []*Feed
	err = json.NewDecoder(file).Decode(&feeds)

	return feeds, err
}

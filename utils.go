// description:
// @author renshiwei
// Date: 2024/7/12

package main

import (
	"encoding/json"
	"os"
	"time"
)

// SaveToJSONFile 将数据保存到 JSON 文件中
func SaveToJSONFile(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// LoadFromJSONFile 从 JSON 文件中读取数据
func LoadFromJSONFile(filename string) ([]*TwitterRes, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []*TwitterRes
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	return data, err
}

// RemoveDuplicates 去重函数
func RemoveDuplicates(tweets []*TwitterRes) []*TwitterRes {
	seen := make(map[string]bool)
	result := make([]*TwitterRes, 0)

	for _, tweet := range tweets {
		if !seen[tweet.Date+tweet.Twitter] {
			seen[tweet.Date+tweet.Twitter] = true
			result = append(result, tweet)
		}
	}
	return result
}

func TimestampToStr(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	layout := "2006-01-02 15:04:05"
	return t.Format(layout)
}

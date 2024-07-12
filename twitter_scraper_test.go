// description:
// @author renshiwei
// Date: 2024/7/11

package main

import (
	"encoding/json"
	"fmt"
	twitterscraper "github.com/n0madic/twitter-scraper"
	"log"
	"net/http"
	"os"
	"testing"
)

// 1. 登录并保存cookie文件
func TestTwitterLoginAndSetCookies(t *testing.T) {
	scraper := twitterscraper.New()
	// 根据需要设置代理
	err := scraper.SetProxy("http://localhost:6152")
	if err != nil {
		panic(err)
	}

	err = scraper.Login("<userId>", "<password>", "<email>")
	if err != nil {
		panic(err)
	}

	cookies := scraper.GetCookies()
	// serialize to JSON
	js, _ := json.Marshal(cookies)
	// save to file
	f, _ := os.Create("cookies.json")
	_, err = f.Write(js)
	if err != nil {
		panic(err)
	}

}

// 2. 根据cookie文件登录并下载推文
func TestTwitterScraper(t *testing.T) {
	scraper := twitterscraper.New()
	err := scraper.SetProxy("http://localhost:6152")
	if err != nil {
		panic(err)
	}

	f, _ := os.Open("cookies.json")
	// deserialize from JSON
	var cookies []*http.Cookie
	err = json.NewDecoder(f).Decode(&cookies)
	if err != nil {
		panic(err)
	}
	// load cookies
	scraper.SetCookies(cookies)
	// check login status
	isLogin := scraper.IsLoggedIn()
	if !isLogin {
		panic("not login")
	}

	//err = scraper.Login("<userId>", "<password>", "<email>")
	//if err != nil {
	//	panic(err)
	//}

	getMaxSize := 100000
	totalFetched := 0
	execCount := 0
	twitterResList := make([]*TwitterRes, 0)
	cursor := ""

	for totalFetched < getMaxSize {
		tweets, nextCursor, err := scraper.FetchTweets("okx", 50, cursor)
		if err != nil {
			log.Fatal("Failed to fetch tweets:", err)
			break
		}

		if len(tweets) == 0 {
			break // 没有更多的推文
		}

		source := fmt.Sprintf("Sourced from %s's Twitter.", "okx")

		for _, tweet := range tweets {
			tw := &TwitterRes{
				Date:    TimestampToStr(tweet.Timestamp),
				Source:  source,
				Twitter: tweet.Text,
			}
			twitterResList = append(twitterResList, tw)
			totalFetched++

		}
		execCount++
		cursor = nextCursor
		fmt.Println("当前抓取次数：", execCount)
		fmt.Println("当前抓取总数量：", totalFetched)
		fmt.Println("nextCursor:", nextCursor)
	}

	fmt.Println("Total tweets fetched:", totalFetched)
	fmt.Println("Last Cursor:", cursor)

	saveTweetsFileName := "okx.json"

	// 读取现有的 tweets.json 文件
	existingTweets, err := LoadFromJSONFile(saveTweetsFileName)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File does not exist, creating new file.")
			existingTweets = []*TwitterRes{} // 如果文件不存在，初始化为空切片
		} else {
			log.Fatal("Failed to load existing tweets from JSON file:", err)
		}
	}

	// 将新抓取的推文数据追加到已存在的数据中
	allTweets := append(existingTweets, twitterResList...)

	// 对整个 JSON 数据进行去重
	uniqueTweets := RemoveDuplicates(allTweets)

	// 将去重后的数据保存回 tweets.json 文件中
	err = SaveToJSONFile(saveTweetsFileName, uniqueTweets)
	if err != nil {
		log.Fatal("Failed to save tweets to JSON file:", err)
	}

	fmt.Println("tweets json save success")
	fmt.Println("allTweets Count:", len(uniqueTweets))
}

func TestJsonAddSource(t *testing.T) {
	saveTweetsFileName := "ssv_network.json"
	existingTweets, err := LoadFromJSONFile(saveTweetsFileName)
	if err != nil {
		panic(err)
	}

	for _, tweet := range existingTweets {
		if tweet.Source == "" {
			tweet.Source = "Sourced from ssv_network's Twitter."
		}
	}

	// 将去重后的数据保存回 tweets.json 文件中
	err = SaveToJSONFile(saveTweetsFileName, existingTweets)
	if err != nil {
		log.Fatal("Failed to save tweets to JSON file:", err)
	}

	fmt.Println("tweets json save success")
	fmt.Println("allTweets Count:", len(existingTweets))
}

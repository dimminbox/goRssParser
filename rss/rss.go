package rss

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
	"goRssParser/core"
	"strings"
	"time"
)

func ParseRss(url string, keywordList []string, keywordExList []string) []core.Result {

	var result []core.Result

	h := md5.New()
	p := bluemonday.UGCPolicy()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURLWithContext(url, ctx)

	if feed == nil {
		return nil
	}

	for _, item := range feed.Items {
		for _, keyword := range keywordList {
			if !strings.Contains(item.Title, keyword) && !strings.Contains(p.Sanitize(item.Description), keyword) {
				continue
			}
			notException := true
			for _, keywordEx := range keywordExList {
				if strings.Contains(item.Title, keywordEx) || strings.Contains(p.Sanitize(item.Description), keywordEx) {
					notException = false
					break
				}
			}

			if notException {
				hash := fmt.Sprintf("%x\n", h.Sum([]byte(fmt.Sprintf("%s%s", item.Link))))
				result = append(result, core.Result{DateCreated: time.Now(), Hash: hash, Url: item.Link,
					Title: item.Title, Keyword: keyword, DatePublicated: *item.PublishedParsed})
				break
			}
		}
	}

	return result
}

package main

import (
	"github.com/cute-angelia/go-utils/syntax/ijson"
	"github.com/cute-angelia/weibospider"
	"log"
)

func main() {
	wb := weibospider.NewWeiboSpider()
	if posts, err := wb.GetUserPosts(2464203544, 2); err == nil {
		log.Println(ijson.Pretty(posts))
	} else {
		log.Println(err)
	}
}

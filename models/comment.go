package models

import "time"

// Comment 评论信息
type Comment struct {
	UID       int32
	Content   string
	WeiboURL  string
	LikeNum   int32
	CreateAt  time.Time
	CrawlTime time.Time
}

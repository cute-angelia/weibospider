package weibospider

import "weibospider/models"

type PostsResponse struct {
	OK   int32 `json:"ok"`
	Data struct {
		Cards []Card `json:"cards"`
	} `json:"data"`
}

type Card struct {
	Type  int32       `json:"card_type"`
	URL   string      `json:"scheme"`
	MBlog models.Post `json:"mblog"`
}

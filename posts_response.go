package weibospider

import "github.com/cute-angelia/weibospider/models"

type PostsResponse struct {
	OK   interface{} `json:"ok"`
	Data struct {
		List  []models.Post `json:"list"`
		Cards []Card        `json:"cards"`
	} `json:"data"`
}

type Card struct {
	Type  int32       `json:"card_type"`
	URL   string      `json:"scheme"`
	MBlog models.Post `json:"mblog"`
}

func (r PostsResponse) Posts() []models.Post {
	if len(r.Data.List) > 0 {
		return r.Data.List
	}

	posts := make([]models.Post, 0, len(r.Data.Cards))
	for _, card := range r.Data.Cards {
		post := card.MBlog
		post.URL = card.URL
		posts = append(posts, post)
	}
	return posts
}

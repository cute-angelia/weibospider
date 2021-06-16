package weibospider

type LongTestResponse struct {
	OK   int32 `json:"ok"`
	Data struct {
		OK             int32  `json:"ok"`
		Content        string `json:"longTextContent"`
		RepostsCount   int32  `json:"reposts_count"`
		CommentsCount  int32  `json:"comments_count"`
		AttitudesCount int32  `json:"attitudes_count"`
	} `json:"data"`
}

package models

type Like struct {
	Data LikeData `json:"data"`
	Ok   int      `json:"ok"`
}

type LikeData struct {
	CardlistInfo struct {
		Containerid string `json:"containerid"`
		Page        int    `json:"page"`
		Title       string `json:"title"`
	} `json:"cardlist_info"`
	Cards     []LikeCard `json:"cards"`
	Count     int        `json:"count"`
	TitleTop  string     `json:"title_top"`
	WeiboNeed string     `json:"weibo_need"`
}

type LikeCard struct {
	CardType int       `json:"card_type"`
	Cols     int       `json:"cols"`
	Pics     []LikePic `json:"pics"`
}

type LikePic struct {
	Pic       string `json:"pic"`
	PicId     string `json:"pic_id"`
	PicMiddle string `json:"pic_middle"`
	PicOri    string `json:"pic_ori"`
	PicSmall  string `json:"pic_small"`
}

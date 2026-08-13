package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

// Post 微博的具体信息
type Post struct {
	URL            string `json:"url"`
	UID            uint64 `json:"uid"`
	PostCreatedAt  string `json:"created_at"`
	ID             string `json:"id" gorm:"primaryKey"`
	MID            string `json:"mid"`
	MBlogID        string `json:"mblogid"`
	Text           string `json:"text"`
	RepostsCount   int32  `json:"reposts_count"`
	CommentsCount  int32  `json:"comments_count"`
	AttitudesCount int32  `json:"attitudes_count"`
	IsLongText     bool   `json:"isLongText"`
	PicNum         int32  `json:"pic_num"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	PageInfo       PageInfo `json:"page_info"`

	PicIDs   []string           `json:"pic_ids"`
	PicInfos map[string]PicInfo `json:"pic_infos"`
	Pics     []Pic              `json:"pics"`
}

func (p *Post) UnmarshalJSON(data []byte) error {
	type postAlias Post
	aux := struct {
		*postAlias
		ID      flexibleString `json:"id"`
		MID     flexibleString `json:"mid"`
		RawText string         `json:"text_raw"`
		User    User           `json:"user"`
	}{postAlias: (*postAlias)(p)}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	p.ID = string(aux.ID)
	p.MID = string(aux.MID)
	if p.Text == "" {
		p.Text = aux.RawText
	}
	if p.ID == "" {
		p.ID = p.MID
	}
	if p.MID == "" {
		p.MID = p.ID
	}
	if p.UID == 0 {
		p.UID = aux.User.ID
	}
	if len(p.Pics) == 0 && len(p.PicInfos) > 0 {
		p.Pics = picsFromPicInfos(p.PicIDs, p.PicInfos)
	}
	return nil
}

type flexibleString string

func (s *flexibleString) UnmarshalJSON(data []byte) error {
	var text string
	if err := json.Unmarshal(data, &text); err == nil {
		*s = flexibleString(text)
		return nil
	}

	var number json.Number
	if err := json.Unmarshal(data, &number); err == nil {
		*s = flexibleString(number.String())
		return nil
	}

	var f float64
	if err := json.Unmarshal(data, &f); err == nil {
		*s = flexibleString(strconv.FormatFloat(f, 'f', -1, 64))
		return nil
	}

	return nil
}

func (p Post) LongTextID() string {
	if p.MBlogID != "" {
		return p.MBlogID
	}
	if p.MID != "" {
		return p.MID
	}
	return p.ID
}

func (p Post) WeiboURL() string {
	if p.UID == 0 || p.MBlogID == "" {
		return p.URL
	}
	return fmt.Sprintf("https://weibo.com/%d/%s", p.UID, p.MBlogID)
}

type Pic struct {
	Pid   string `json:"pid"`
	Geo   PicGeo `json:"geo"`
	Large struct {
		Size string `json:"size"`
		Url  string `json:"url"`
		Geo  PicGeo `json:"geo"`
	} `json:"large"`
}

type PicInfo struct {
	PicID    string       `json:"pic_id"`
	Large    PicInfoImage `json:"large"`
	Largest  PicInfoImage `json:"largest"`
	Original PicInfoImage `json:"original"`
	Bmiddle  PicInfoImage `json:"bmiddle"`
}

type PicInfoImage struct {
	Url    string      `json:"url"`
	Width  interface{} `json:"width"`
	Height interface{} `json:"height"`
}

type PicGeo struct {
	Width  interface{} `json:"width"`
	Height interface{} `json:"height"`
	Croped bool        `json:"croped"`
}

func picsFromPicInfos(picIDs []string, picInfos map[string]PicInfo) []Pic {
	if len(picIDs) == 0 {
		for picID := range picInfos {
			picIDs = append(picIDs, picID)
		}
	}

	pics := make([]Pic, 0, len(picIDs))
	for _, picID := range picIDs {
		info, ok := picInfos[picID]
		if !ok {
			continue
		}

		image := bestPicInfoImage(info)
		if image.Url == "" {
			continue
		}

		pid := info.PicID
		if pid == "" {
			pid = picID
		}
		pic := Pic{Pid: pid}
		pic.Large.Url = image.Url
		pic.Large.Geo.Width = image.Width
		pic.Large.Geo.Height = image.Height
		pics = append(pics, pic)
	}
	return pics
}

func bestPicInfoImage(info PicInfo) PicInfoImage {
	for _, image := range []PicInfoImage{info.Largest, info.Original, info.Large, info.Bmiddle} {
		if image.Url != "" {
			return image
		}
	}
	return PicInfoImage{}
}

type PageInfo struct {
	Type       string `json:"type"`
	ObjectType string `json:"object_type"`
	Urls       struct {
		Mp4720pMp4 string `json:"mp4_720p_mp4"`
		Mp4HdMp4   string `json:"mp4_hd_mp4"`
		Mp4LdMp4   string `json:"mp4_ld_mp4"`
	} `json:"urls"`
	MediaInfo PageMediaInfo `json:"media_info"`
}

type PageMediaInfo struct {
	StreamURL   string `json:"stream_url"`
	StreamURLHD string `json:"stream_url_hd"`
	Mp4720pMp4  string `json:"mp4_720p_mp4"`
	Mp4HdURL    string `json:"mp4_hd_url"`
	Mp4SDURL    string `json:"mp4_sd_url"`
}

func (p *PageInfo) UnmarshalJSON(data []byte) error {
	type pageInfoAlias PageInfo
	aux := struct {
		*pageInfoAlias
		Type flexibleString `json:"type"`
	}{pageInfoAlias: (*pageInfoAlias)(p)}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	p.Type = string(aux.Type)
	if p.Urls.Mp4720pMp4 == "" {
		p.Urls.Mp4720pMp4 = p.MediaInfo.Mp4720pMp4
	}
	if p.Urls.Mp4HdMp4 == "" {
		p.Urls.Mp4HdMp4 = firstNonEmpty(p.MediaInfo.Mp4HdURL, p.MediaInfo.StreamURLHD, p.MediaInfo.StreamURL)
	}
	if p.Urls.Mp4LdMp4 == "" {
		p.Urls.Mp4LdMp4 = firstNonEmpty(p.MediaInfo.Mp4SDURL, p.MediaInfo.StreamURL)
	}
	return nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

func (p Post) Save() error {
	if err := CreateDirIfNotExist("./output"); err != nil {
		log.WithFields(log.Fields{"dir": "./output", "err": err.Error()}).Error("failed to create directory")
		return err
	}
	if err := CreateDirIfNotExist(fmt.Sprintf("./output/%v", p.UID)); err != nil {
		log.WithFields(log.Fields{"dir": fmt.Sprintf("./output/%v", p.UID), "err": err.Error()}).Error("failed to create directory")
		return err
	}

	filename := fmt.Sprintf("./output/%v/%v.json", p.UID, p.ID)
	if !FileExist(filename) {
		jstr, _ := json.MarshalIndent(p, "", "")
		_ = ioutil.WriteFile(filename, jstr, 0644)
	}

	return nil
}

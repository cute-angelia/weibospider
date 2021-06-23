package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	Text           string `json:"text"`
	RepostsCount   int32  `json:"reposts_count"`
	CommentsCount  int32  `json:"comments_count"`
	AttitudesCount int32  `json:"attitudes_count"`
	IsLongText     bool   `json:"isLongText"`
	PicNum         int32  `json:"pic_num"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
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

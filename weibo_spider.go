package weibospider

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cute-angelia/weibospider/models"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	userInfoUrlFmt  = "https://m.weibo.cn/api/container/getIndex?type=__uid&value=%d&containerid=100505%d"
	userPostsUrlFmt = "https://m.weibo.cn/api/container/getIndex?type=__uid&value=%d&containerid=107603%d&page=%d"
	longTextUrlFmt  = "https://m.weibo.cn/statuses/extend?id=%s"
)

func init() {
	log.SetReportCaller(true)
}

type weiboSpider struct {
	delay    time.Duration
	wg       sync.WaitGroup
	longtext bool
}

func NewWeiboSpider(options ...Option) *weiboSpider {
	// default
	c := &weiboSpider{
		delay:    5 * time.Second,
		wg:       sync.WaitGroup{},
		longtext: false,
	}

	for _, option := range options {
		option(c)
	}

	return c
}

// getUserInfoUrl 生成用户信息 URL
func getUserInfoUrl(uid uint64) string {
	return fmt.Sprintf(userInfoUrlFmt, uid, uid)
}

// getUserPostsUrl 生成微博列表爬取 URL
func getUserPostsUrl(uid uint64, page uint32) string {
	return fmt.Sprintf(userPostsUrlFmt, uid, uid, page)
}

// 生成长微博爬取 URL
func getLongTextUrl(id string) string {
	return fmt.Sprintf(longTextUrlFmt, id)
}

// randomUserAgent 随机生成 UserAgent
func randomUserAgent() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// RandomSleep 随机睡眠
func RandomSleep(minSeconds, maxSeconds int32) {
	dura := minSeconds*1000 + rand.Int31n((maxSeconds-minSeconds)*1000)
	log.WithField("duration", dura).Debug("sleep(ms)")
	time.Sleep(time.Duration(dura) * time.Millisecond)
}

// GetUserInfo 爬取用户信息
func (wb *weiboSpider) GetUserInfo(uid uint64) (models.User, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", getUserInfoUrl(uid), nil)
	if err != nil {
		log.WithFields(log.Fields{"req": req.URL, "err": err.Error()}).Error()
		return models.User{}, err
	}
	req.Header.Set("User-Agent", randomUserAgent())

	resp, err := client.Do(req)
	if err != nil {
		log.WithField("err", err.Error()).Error("request failed")
		return models.User{}, err
	}
	defer resp.Body.Close()

	log.WithFields(log.Fields{"status": resp.StatusCode, "url": resp.Request.URL.String()}).Info("request success")

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.WithField("err", err.Error()).Error("response read failed")
		return models.User{}, err
	}
	uinfores := UserInfoResponse{}
	if err := json.Unmarshal(body, &uinfores); err != nil {
		log.WithField("err", err.Error()).Error("response unmarshal failed")
		return models.User{}, err
	}

	if uinfores.OK != 1 {
		log.WithField("response", string(body)).Error("response not ok")
		return models.User{}, errors.New("response not ok")
	}

	return uinfores.Data.UserInfo, nil
}

//  爬取长微博文本
func getLongText(id string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", getLongTextUrl(id), nil)
	if err != nil {
		log.WithFields(log.Fields{"req": req.URL, "err": err.Error()}).Error()
		return "", err
	}
	req.Header.Set("User-Agent", randomUserAgent())

	resp, err := client.Do(req)
	if err != nil {
		log.WithField("err", err.Error()).Error("request failed")
		return "", err
	}
	defer resp.Body.Close()

	log.WithFields(log.Fields{"status": resp.StatusCode, "url": resp.Request.URL.String()}).Info("request success")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithField("err", err.Error()).Error("response read failed")
		return "", err
	}

	ltresp := LongTestResponse{}
	if err := json.Unmarshal(body, &ltresp); err != nil {
		log.WithFields(log.Fields{"err": err.Error(), "response": string(body)}).Error("response unmarshal failed")
		return "", err
	}

	if ltresp.OK != 1 || ltresp.Data.OK != 1 {
		log.WithField("response", string(body)).Error("response not ok")
		return "", errors.New("response not ok")
	}
	return ltresp.Data.Content, nil
}

// GetUserPosts 爬取用户微博
func (wb *weiboSpider) GetUserPosts(uid uint64, page uint32) ([]models.Post, error) {
	// delay
	wb.wg.Wait()
	if wb.delay > 0 {
		defer func() {
			wb.wg.Add(1)
			go func() {
				time.Sleep(wb.delay)
				wb.wg.Done()
			}()
		}()
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", getUserPostsUrl(uid, page), nil)
	if err != nil {
		log.WithFields(log.Fields{"req": req.URL, "err": err.Error()}).Error()
		return []models.Post{}, err
	}
	req.Header.Set("User-Agent", randomUserAgent())

	resp, err := client.Do(req)
	if err != nil {
		log.WithField("err", err.Error()).Error("request failed")
		return []models.Post{}, err
	}
	defer resp.Body.Close()

	log.WithFields(log.Fields{"status": resp.StatusCode, "url": resp.Request.URL.String()}).Info("request success")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithField("err", err.Error()).Error("response read failed")
		return []models.Post{}, err
	}

	uposts := PostsResponse{}
	if err := json.Unmarshal(body, &uposts); err != nil {
		log.WithField("err", err.Error()).Error("response unmarshal failed")
		return []models.Post{}, err
	}

	posts := []models.Post{}
	for _, card := range uposts.Data.Cards {
		if card.MBlog.IsLongText && wb.longtext {
			RandomSleep(2, 5)
			content, err := getLongText(card.MBlog.ID)
			if err != nil {
				log.WithFields(log.Fields{"err": err.Error()}).Error("failed to get longtext")
				continue
			}
			card.MBlog.Text = content
		}
		post := card.MBlog
		post.URL = card.URL
		post.UID = uid
		posts = append(posts, post)
	}

	if uposts.OK != 1 {
		log.WithField("response", string(body)).Error("response not ok")
		return []models.Post{}, errors.New("response not ok")
	}
	return posts, nil
}

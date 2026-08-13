package weibospider

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cute-angelia/weibospider/models"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	defaultBaseURL = "https://weibo.com"
	defaultUA      = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0 Safari/537.36"

	userInfoUrlFmt  = "%s/ajax/profile/info?uid=%d"
	userPostsUrlFmt = "%s/ajax/statuses/searchProfile?uid=%d&page=%d&hasori=1&hastext=1&haspic=1&hasvideo=1&hasmusic=1&hasret=1"
	longTextUrlFmt  = "%s/ajax/statuses/longtext?id=%s"
)

func init() {
	log.SetReportCaller(true)
}

type weiboSpider struct {
	delay      time.Duration
	wg         sync.WaitGroup
	longtext   bool
	cookie     string
	baseURL    string
	userAgent  string
	httpClient *http.Client
}

func NewWeiboSpider(options ...Option) *weiboSpider {
	// default
	c := &weiboSpider{
		delay:     5 * time.Second,
		wg:        sync.WaitGroup{},
		longtext:  false,
		cookie:    os.Getenv("WEIBO_COOKIE"),
		baseURL:   defaultBaseURL,
		userAgent: defaultUA,
		httpClient: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	}

	for _, option := range options {
		option(c)
	}

	return c
}

func GetUserInfo(uid uint64) (models.User, error) {
	return NewWeiboSpider().GetUserInfo(uid)
}

func GetUserPosts(uid uint64, page uint32) ([]models.Post, error) {
	return NewWeiboSpider().GetUserPosts(uid, page)
}

// getUserInfoUrl 生成用户信息 URL
func (wb *weiboSpider) getUserInfoUrl(uid uint64) string {
	return fmt.Sprintf(userInfoUrlFmt, wb.baseURL, uid)
}

// getUserPostsUrl 生成微博列表爬取 URL
func (wb *weiboSpider) getUserPostsUrl(uid uint64, page uint32) string {
	return fmt.Sprintf(userPostsUrlFmt, wb.baseURL, uid, page)
}

// 生成长微博爬取 URL
func (wb *weiboSpider) getLongTextUrl(id string) string {
	return fmt.Sprintf(longTextUrlFmt, wb.baseURL, id)
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
	if maxSeconds <= minSeconds {
		if minSeconds > 0 {
			time.Sleep(time.Duration(minSeconds) * time.Second)
		}
		return
	}
	dura := minSeconds*1000 + rand.Int31n((maxSeconds-minSeconds)*1000)
	log.WithField("duration", dura).Debug("sleep(ms)")
	time.Sleep(time.Duration(dura) * time.Millisecond)
}

func (wb *weiboSpider) getJSON(rawURL string, target interface{}) error {
	if strings.TrimSpace(wb.cookie) == "" {
		return errors.New("weibo cookie is required; pass WithCookie(...) or set WEIBO_COOKIE")
	}

	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", wb.userAgent)
	req.Header.Set("Cookie", wb.cookie)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Referer", wb.baseURL+"/")

	resp, err := wb.httpClient.Do(req)
	if err != nil {
		log.WithField("err", err.Error()).Error("request failed")
		return err
	}
	defer resp.Body.Close()

	log.WithFields(log.Fields{"status": resp.StatusCode, "url": resp.Request.URL.String()}).Info("request success")

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("weibo request failed: status %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		log.WithField("err", err.Error()).Error("response unmarshal failed")
		return err
	}
	return nil
}

func responseOK(ok interface{}) bool {
	switch v := ok.(type) {
	case bool:
		return v
	case float64:
		return v == 1
	case int:
		return v == 1
	case int32:
		return v == 1
	default:
		return false
	}
}

// GetUserInfo 爬取用户信息
func (wb *weiboSpider) GetUserInfo(uid uint64) (models.User, error) {
	uinfores := UserInfoResponse{}
	if err := wb.getJSON(wb.getUserInfoUrl(uid), &uinfores); err != nil {
		return models.User{}, err
	}

	if !responseOK(uinfores.OK) {
		log.WithField("response", uinfores).Error("response not ok")
		return models.User{}, errors.New("response not ok")
	}

	if uinfores.Data.User.ID != 0 {
		return uinfores.Data.User, nil
	}
	return uinfores.Data.UserInfo, nil
}

// 爬取长微博文本
func (wb *weiboSpider) getLongText(id string) (string, error) {
	ltresp := LongTestResponse{}
	if err := wb.getJSON(wb.getLongTextUrl(id), &ltresp); err != nil {
		return "", err
	}

	if !responseOK(ltresp.OK) {
		log.WithField("response", ltresp).Error("response not ok")
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

	uposts := PostsResponse{}
	if err := wb.getJSON(wb.getUserPostsUrl(uid, page), &uposts); err != nil {
		return []models.Post{}, err
	}

	posts := []models.Post{}
	for _, post := range uposts.Posts() {
		if post.IsLongText && wb.longtext {
			if wb.delay > 0 {
				RandomSleep(2, 5)
			}
			content, err := wb.getLongText(post.LongTextID())
			if err != nil {
				log.WithFields(log.Fields{"err": err.Error()}).Error("failed to get longtext")
				continue
			}
			post.Text = content
		}
		post.UID = uid
		if post.URL == "" {
			post.URL = post.WeiboURL()
		}
		posts = append(posts, post)
	}

	if !responseOK(uposts.OK) {
		log.WithField("response", uposts).Error("response not ok")
		return []models.Post{}, errors.New("response not ok")
	}
	return posts, nil
}

// 用户喜欢
func (wb *weiboSpider) GetUserLikes() {
	// https://m.weibo.cn/api/container/getSecond?containerid=1078031878498994_-_photolike&page=2&count=24&title=%E8%B5%9E%E8%BF%87%E7%9A%84%E5%9B%BE%E7%89%87&luicode=10000011&lfid=1078031878498994&type=like
}

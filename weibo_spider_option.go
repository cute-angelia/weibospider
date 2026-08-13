package weibospider

import (
	"net/http"
	"strings"
	"time"
)

type Option func(c *weiboSpider)

func WithDelay(delay time.Duration) Option {
	return func(c *weiboSpider) {
		c.delay = delay
	}
}

func WithLongText(longtext bool) Option {
	return func(c *weiboSpider) {
		c.longtext = longtext
	}
}

func WithCookie(cookie string) Option {
	return func(c *weiboSpider) {
		c.cookie = strings.TrimSpace(cookie)
	}
}

func WithUserAgent(userAgent string) Option {
	return func(c *weiboSpider) {
		if strings.TrimSpace(userAgent) != "" {
			c.userAgent = userAgent
		}
	}
}

func WithHTTPClient(client *http.Client) Option {
	return func(c *weiboSpider) {
		if client != nil {
			c.httpClient = client
		}
	}
}

func withBaseURL(baseURL string) Option {
	return func(c *weiboSpider) {
		c.baseURL = strings.TrimRight(baseURL, "/")
	}
}

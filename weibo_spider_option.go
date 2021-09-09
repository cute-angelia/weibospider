package weibospider

import (
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

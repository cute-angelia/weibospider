# WeiboSpider

基于微博手机网页版接口实现的 weibo 爬虫，不需要账户登陆

## 使用介绍

### 爬取用户信息

```golang
package main

import (
    "github/wyxpku/weibospider/spider"
)

func main() {
    spider.GetUserInfo()
}

```

# weibospider

基于微博网页版接口实现的 weibo 爬虫。当前微博接口需要登录态 cookie，请先从浏览器复制 `weibo.com` 请求里的 cookie。

## Cookie

推荐使用环境变量：

```bash
export WEIBO_COOKIE='SUB=...; SUBP=...; XSRF-TOKEN=...'
```

也可以在代码里传入：

```golang
wb := weibospider.NewWeiboSpider(weibospider.WithCookie("SUB=...; SUBP=..."))
```

## 使用介绍

### 爬取用户信息

```golang
package main

import (
	"fmt"

	"github.com/wyxpku/weibospider"
)

func main() {
	wb := weibospider.NewWeiboSpider()
	uinfo, _ := wb.GetUserInfo(2993720115)
	fmt.Printf("%#v\n", uinfo)
}
```


### 爬取用户微博

```golang
package main

import (
	"fmt"

	"github.com/wyxpku/weibospider"
)

func main() {
	wb := weibospider.NewWeiboSpider(weibospider.WithLongText(true))
	posts, _ := wb.GetUserPosts(2993720115, 1)
	fmt.Printf("%#v\n", posts)
}
```

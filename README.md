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

### 通过环境变量使用

```golang
package main

import (
	"fmt"
	"log"

	"github.com/cute-angelia/weibospider"
)

func main() {
	wb := weibospider.NewWeiboSpider()

	user, err := wb.GetUserInfo(3261134763)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("用户：%s，粉丝：%d\n", user.Name, user.FollowersCount)

	posts, err := wb.GetUserPosts(3261134763, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("第 1 页微博数：%d\n", len(posts))
	if len(posts) > 0 {
		fmt.Printf("第一条：%s\n%s\n", posts[0].URL, posts[0].Text)
	}
}
```

运行：

```bash
export WEIBO_COOKIE='SUB=...; SUBP=...; XSRF-TOKEN=...'
go run ./cmd
```

### 在代码里传入 cookie

```golang
package main

import (
	"fmt"
	"log"

	"github.com/cute-angelia/weibospider"
)

func main() {
	wb := weibospider.NewWeiboSpider(
		weibospider.WithCookie("SUB=...; SUBP=...; XSRF-TOKEN=..."),
		weibospider.WithLongText(true),
	)

	posts, err := wb.GetUserPosts(3261134763, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("第 1 页微博数：%d\n", len(posts))
}
```

# weibospider

基于微博手机网页版接口实现的 weibo 爬虫，不需要账户登陆

## 使用介绍

### 爬取用户信息

```golang
package main

import (
	"fmt"

	"github.com/wyxpku/weibospider"
)

func main() {
	uinfo, _ := weibospider.GetUserInfo(2993720115)
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
	posts, _ := weibospider.GetUserPosts(2993720115, 1)
	fmt.Printf("%#v\n", posts)
}
```

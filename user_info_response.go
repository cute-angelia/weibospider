package weibospider

import "github.com/cute-angelia/weibospider/models"

type UserInfoResponse struct {
	OK   interface{} `json:"ok"`
	Data struct {
		User     models.User `json:"user"`
		UserInfo models.User `json:"userInfo"`
	} `json:"data"`
}

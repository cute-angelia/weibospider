package weibospider

import "github.com/cute-angelia/weibospider/models"

type UserInfoResponse struct {
	OK   int32 `json:"ok"`
	Data struct {
		UserInfo models.User `json:"userInfo"`
	} `json:"data"`
}

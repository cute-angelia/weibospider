package models

// User 微博用户信息
type User struct {
	ID              uint64 `json:"id" gorm:"primaryKey"`
	Name            string `json:"screen_name"`
	Verified        bool   `json:"verified"`
	VerifiedType    int32  `json:"verified_type"`
	VerifiedTypeExt int32  `json:"verified_type_ext"`
	VerifiedReasone int32  `json:"verified_reason"`
	Description     string `json:"description"`
	Gender          string `json:"gender"`
	FollowersCount  int32  `json:"followers_count"`
	FollowCount     int32  `json:"follow_count"`
	AvatorURL       string `json:"avatar_hd"`
}

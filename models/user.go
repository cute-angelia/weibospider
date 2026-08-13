package models

import "encoding/json"

// User 微博用户信息
type User struct {
	ID              uint64 `json:"id" gorm:"primaryKey"`
	Name            string `json:"screen_name"`
	Verified        bool   `json:"verified"`
	VerifiedType    int32  `json:"verified_type"`
	VerifiedTypeExt int32  `json:"verified_type_ext"`
	VerifiedReasone string `json:"verified_reason"`
	Description     string `json:"description"`
	Gender          string `json:"gender"`
	FollowersCount  int32  `json:"followers_count"`
	FollowCount     int32  `json:"follow_count"`
	AvatorURL       string `json:"avatar_hd"`
}

func (u *User) UnmarshalJSON(data []byte) error {
	type userAlias User
	aux := struct {
		*userAlias
		FriendsCount    int32  `json:"friends_count"`
		ProfileImageURL string `json:"profile_image_url"`
	}{userAlias: (*userAlias)(u)}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if u.FollowCount == 0 {
		u.FollowCount = aux.FriendsCount
	}
	if u.AvatorURL == "" {
		u.AvatorURL = aux.ProfileImageURL
	}
	return nil
}

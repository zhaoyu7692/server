package model

import (
	"time"
)

type UserType int

const (
	UserTypeGeneral UserType = 0
	UserTypeAdmin   UserType = 1
)

type User struct {
	Uid           int64      `json:"uid,omitempty" db:"UID"`
	Username      string     `json:"username,omitempty" db:"USERNAME"` // 用户名
	Password      string     `json:"password,omitempty" db:"PASSWORD"` // SHA256编码后的密码
	Token         string     `json:"token,omitempty" db:"TOKEN"`       // 令牌
	UserType      UserType   `json:"user_type" db:"USER_TYPE"`
	LastAuthority *time.Time `json:"last_modified,omitempty" db:"LAST_AUTHORITY"`
	//Timestamp
}

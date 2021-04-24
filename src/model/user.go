package model

import "time"

type User struct {
	Uid          int64     `json:"uid,omitempty" db:"UID"`
	Username     string    `json:"username,omitempty" db:"USERNAME"` // 用户名
	Password     string    `json:"password,omitempty" db:"PASSWORD"` // SHA256编码后的密码
	Token        string    `json:"token,omitempty" db:"TOKEN"`       // 令牌
	LastModified *time.Time `json:"last_modified,omitempty" db:"LAST_MODIFIED"`
	//Timestamp
}

package model

import (
	"encoding/json"
	"time"
)

type Contest struct {
	Cid           int64      `json:"cid" db:"CID"`
	Title         string     `json:"title" db:"TITLE"`
	Indexes       []int64    `json:"indexes"`
	BeginTime     *time.Time `json:"begin_time" db:"BEGIN_TIME"`
	Duration      int64      `json:"duration" db:"DURATION"`
	RegisterCount int64      `json:"register_count" db:"REGISTER_COUNT"`
}

func (c Contest) MarshalJSON() ([]byte, error) {
	type Alias Contest
	var beginTime string
	if c.BeginTime != nil {
		beginTime = c.BeginTime.Format("2006-01-02 15:04:05")
	}
	return json.Marshal(&struct {
		Alias
		BeginTime string `json:"begin_time"`
	}{
		Alias:     Alias(c),
		BeginTime: beginTime,
	})
}

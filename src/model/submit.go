package model

import (
	"encoding/json"
	"time"
)

type Submit struct {
	Rid                int64       `json:"rid" db:"RID"`
	Cid                int64       `json:"cid" db:"CID"`
	Index              int64       `json:"index" db:"INDEX"`
	Code               string      `json:"code" db:"CODE"`
	Status             JudgeStatus `json:"status" db:"STATUS"`
	Language           int64       `json:"language" db:"LANGUAGE"`
	TimeCost           *int64      `json:"time_cost" db:"RUN_TIME"`
	CurrentCase        int64       `json:"current_case" db:"CURRENT_CASE"`
	MemoryCost         *int64      `json:"memory_cost" db:"RUN_MEMORY"`
	CompilationMessage *string     `json:"compilation_message" db:"COMPILATION_MESSAGE"`
	SubmitTime         *time.Time  `json:"submit_time,omitempty" db:"SUBMIT_TIME"`
	User
	//Timestamp
}

func (s Submit) MarshalJSON() ([]byte, error) {
	type Alias Submit
	var submitTime string
	if s.SubmitTime != nil {
		submitTime = s.SubmitTime.Format("2006-01-02 15:04:05")
	}
	return json.Marshal(&struct {
		Alias
		SubmitTime string `json:"submit_time,omitempty"`
	}{
		Alias:      Alias(s),
		SubmitTime: submitTime,
	})
}

type JudgeSubmitModel struct {
	Rid         int64  `json:"rid" db:"RID"`
	Cid         int64  `json:"cid" db:"CID"`
	Index       int64  `json:"index" db:"INDEX"`
	Code        string `json:"code" db:"CODE"`
	Status      int64  `json:"status" db:"STATUS"`
	Language    int64  `json:"language" db:"LANGUAGE"`
	Pid         int64  `json:"pid" db:"PID"`
	TimeLimit   int64  `json:"time_limit,omitempty" db:"TIME_LIMIT"`
	MemoryLimit int64  `json:"memory_limit,omitempty" db:"MEMORY_LIMIT"`
}

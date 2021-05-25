package dao

import (
	"encoding/json"
	"main/mysql"
	"time"
)

type ContestTableModel struct {
	CID           int64      `json:"cid" db:"CID"`
	Title         string     `json:"title" db:"TITLE"`
	BeginTime     *time.Time `json:"begin_time" db:"BEGIN_TIME"`
	Duration      int64      `json:"duration" db:"DURATION"`
	RegisterCount int64      `json:"register_count" db:"REGISTER_COUNT"`
}

func (c ContestTableModel) MarshalJSON() ([]byte, error) {
	type Alias ContestTableModel
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

type ContestProblemMappingTableModel struct {
	CID    int64 `json:"cid" db:"CID"`
	Index  int64 `json:"index" db:"INDEX"`
	PID    int64 `json:"pid" db:"PID"`
	Accept int64 `json:"accept" db:"ACCEPT"`
	Total  int64 `json:"total" db:"TOTAL"`
}

type ProblemTableModel struct {
	Pid         int64  `json:"pid" db:"PID"`
	Title       string `json:"title" db:"TITLE"`
	Description string `json:"description" db:"DESCRIPTION"`
	Difficulty  int64  `json:"difficulty" db:"DIFF"`
	Input       string `json:"input" db:"INPUT"`
	Output      string `json:"output" db:"OUTPUT"`
	TimeLimit   int64  `json:"time_limit" db:"TIME_LIMIT"`
	MemoryLimit int64  `json:"memory_limit" db:"MEMORY_LIMIT"`
	CaseCount   int64  `json:"case_count" db:"CASE_COUNT"`
	Source      string `json:"source" db:"SOURCE"`
}

type SampleTableModel struct {
	PID    int64  `json:"pid" db:"PID"`
	SID    int64  `json:"sid" db:"SID"`
	Input  string `json:"input" db:"INPUT"`
	Output string `json:"output" db:"OUTPUT"`
}

func GetContestWithCid(cid int64) *ContestTableModel {
	var contest ContestTableModel
	if err := mysql.DBConn.Get(&contest, "SELECT CID, TITLE, BEGIN_TIME, DURATION, REGISTER_COUNT FROM contest WHERE CID = ?", cid); err != nil {
		return nil
	}
	return &contest
}

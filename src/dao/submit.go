package dao

import (
	"main/model"
	"main/mysql"
	"time"
)

type SubmitTableModel struct {
	Rid                int64             `json:"rid" db:"RID"`
	Cid                int64             `json:"cid" db:"CID"`
	Index              int64             `json:"index" db:"INDEX"`
	Uid                int64             `json:"uid" db:"UID"`
	Code               string            `json:"code" db:"CODE"`
	Status             model.JudgeStatus `json:"status" db:"STATUS"`
	Language           int64             `json:"language" db:"LANGUAGE"`
	TimeCost           *int64            `json:"time_cost" db:"RUN_TIME"`
	MemoryCost         *int64            `json:"memory_cost" db:"RUN_MEMORY"`
	CurrentCase        int64             `json:"current_case" db:"CURRENT_CASE"`
	CompilationMessage *string           `json:"compilation_message" db:"COMPILATION_MESSAGE"`
	SubmitTime         *time.Time        `json:"submit_time,omitempty" db:"SUBMIT_TIME"`
}

func GetSubmitWithRid(rid int64) *SubmitTableModel {
	submitModel := SubmitTableModel{}
	if err := mysql.DBConn.Get(&submitModel, "SELECT RID, CID, `INDEX`, UID, CODE, STATUS, LANGUAGE, RUN_TIME, RUN_MEMORY, CURRENT_CASE, COMPILATION_MESSAGE, SUBMIT_TIME FROM submit WHERE RID = ?", rid); err != nil {
		return nil
	}
	return &submitModel
}

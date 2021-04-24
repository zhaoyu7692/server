package model

import "time"

//type ResponseTimestamp struct {
//	TimestampCreated  string `json:"timestamp_created,omitempty"`
//	TimestampModified string `json:"timestamp_modified,omitempty"`
//}

type Timestamp struct {
	GmtCreated  *time.Time `json:"gmt_created,omitempty" db:"GMT_CREATED"`
	GmtModified *time.Time `json:"gmt_modified,omitempty" db:"GMT_MODIFIED"`
	//ResponseTimestamp
}

type RequestPaginationModel struct {
	Page int64 `json:"page"`
	Size int64 `json:"size"`
}

type ResponsePaginationModel struct {
	Size  int64 `json:"size"`
	Total int64 `json:"total"`
}

type ResponseCode int

const (
	PublicFail ResponseCode = iota
	Success
	JumpLogin
)

type ResponseBaseModel struct {
	Code    ResponseCode `json:"code"`
	Message string       `json:"message,omitempty"`
}

type JudgeStatus int

const (
	JudgeStatusSystemError JudgeStatus = -1
	JudgeStatusWaiting                 = iota
	JudgeStatusCompiling
	JudgeStatusCompilationError
	JudgeStatusCompilationTimeLimitExceeded
	JudgeStatusRunning
	JudgeStatusTimeLimitExceeded
	JudgeStatusMemoryLimitExceeded
	JudgeStatusOutputLimitExceeded
	JudgeStatusRuntimeError
	JudgeStatusPresentationError
	JudgeStatusWrongAnswer
	JudgeStatusAccept
)

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
	JudgeStatusSystemError                  JudgeStatus = -1
	JudgeStatusWaiting                                  = 0
	JudgeStatusCompiling                                = 1
	JudgeStatusCompilationError                         = 2
	JudgeStatusCompilationTimeLimitExceeded             = 3
	JudgeStatusRunning                                  = 4
	JudgeStatusTimeLimitExceeded                        = 5
	JudgeStatusMemoryLimitExceeded                      = 6
	JudgeStatusOutputLimitExceeded                      = 7
	JudgeStatusRuntimeError                             = 8
	JudgeStatusPresentationError                        = 9
	JudgeStatusWrongAnswer                              = 10
	JudgeStatusAccept                                   = 11
)

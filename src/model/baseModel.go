package model

type Timestamp struct {
	GmtCreated  string `json:"gmt_created,omitempty" db:"GMT_CREATED"`
	GmtModified string `json:"gmt_modified,omitempty" db:"GMT_MODIFIED"`
}

type RequestPaginationModel struct {
	Offset int64 `json:"offset"`
	Size int64 `json:"size"`
}

type ResponsePaginationModel struct {
	Size int64 `json:"size"`
	Total int64 `json:"total"`
}
package model

type BaseSubmit struct {
	Rid                int64   `json:"rid" db:"RID"`
	Pid                int64   `json:"pid" db:"PID"`
	Code               string  `json:"code" db:"CODE"`
	Language           int64   `json:"language" db:"LANGUAGE"`
	Status             int64   `json:"status" db:"STATUS"`
	User
	Timestamp
}

type Submit struct {
	BaseSubmit
	TimeCost   *int64 `json:"time_cost" db:"RUN_TIME"`
	MemoryCost *int64 `json:"memory_cost" db:"RUN_MEMORY"`
	CompilationMessage *string `json:"compilation_message" db:"COMPILATION_MESSAGE"`
	User
}

//type SubmitStatus struct {
//	BaseSubmit
//
//}

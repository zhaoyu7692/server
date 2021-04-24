package model

type Example struct {
	Pid    int64  `json:"pid,omitempty"`
	Sid    int64  `json:"sid,omitempty"`
	Input  string `json:"input" db:"INPUT"`
	Output string `json:"output" db:"OUTPUT"`
	//Timestamp
}

type Problem struct {
	Pid               int64     `json:"pid" db:"PID"`
	Title             string    `json:"title,omitempty" db:"TITLE"`
	Description       string    `json:"description,omitempty" db:"DESCRIPTION"`
	Difficulty        int64     `json:"difficulty,omitempty" db:"DIFF"`
	InputDescription  string    `json:"input_description,omitempty" db:"INPUT"`
	OutputDescription string    `json:"output_description,omitempty" db:"OUTPUT"`
	TimeLimit         int64     `json:"time_limit,omitempty" db:"TIME_LIMIT"`
	MemoryLimit       int64     `json:"memory_limit,omitempty" db:"MEMORY_LIMIT"`
	CaseCount         int64     `json:"case_count,omitempty" db:"CASE_COUNT"`
	Accept            int64     `json:"accept" db:"ACCEPT"`
	Total             int64     `json:"total" db:"TOTAL"`
	Sample            []Example `json:"samples,omitempty"`
	Source            string    `json:"source,omitempty" db:"SOURCE"`
	//Timestamp
}

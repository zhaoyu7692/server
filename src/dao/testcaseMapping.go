package dao

import "main/mysql"

type TestcaseMappingTableModel struct {
	Pid      string `json:"pid" db:"PID"`
	Filename string `json:"filename" db:"FILENAME"`
	Key      string `json:"key" db:"KEY"`
	Path     string `json:"path" db:"PATH"`
}

func GetTestCasesWithPid(pid int64) *[]TestcaseMappingTableModel {
	var testcases []TestcaseMappingTableModel
	if err := mysql.DBConn.Select(&testcases, "SELECT PID, FILENAME, `KEY`, PATH FROM testcase_mapping WHERE PID = ?", pid); err != nil {
		return nil
	}
	return &testcases
}

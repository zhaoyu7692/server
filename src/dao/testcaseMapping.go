package dao

import "main/mysql"

type TestcaseMappingTableModel struct {
	Pid      string `db:"PID"`
	Filename string `db:"FILENAME"`
	Key      string `db:"KEY"`
	Path     string `db:"PATH"`
}

func GetTestCasesWithPid(pid int64) *[]TestcaseMappingTableModel {
	var testcases []TestcaseMappingTableModel
	if err := mysql.DBConn.Select(&testcases, "SELECT PID, FILENAME, `KEY`, PATH FROM testcase_mapping WHERE PID = ?", pid); err != nil {
		return nil
	}
	return &testcases
}

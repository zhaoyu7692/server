package dao

import "main/mysql"

func GetProblemWithPid(pid int64) *ProblemTableModel {
	var problem ProblemTableModel
	if err := mysql.DBConn.Get(&problem, "SELECT PID, TITLE, DESCRIPTION, DIFF, INPUT, OUTPUT, SOURCE, TIME_LIMIT, MEMORY_LIMIT FROM problem WHERE PID = ?", pid); err != nil {
		return nil
	}
	return &problem
}
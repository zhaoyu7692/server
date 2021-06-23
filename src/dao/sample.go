package dao

import "main/mysql"

func GetSamplesWithPid(pid int64) *[]SampleTableModel{
	var samples []SampleTableModel
	if err := mysql.DBConn.Select(&samples, "SELECT PID, SID, INPUT, OUTPUT FROM sample WHERE PID = ? ORDER BY SID", pid); err != nil {
		return nil
	}
	return &samples
}

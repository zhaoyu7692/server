package dao

import "main/mysql"

type ProblemDetailModel struct {
	ProblemTableModel
	Accept int64
	Total  int64
	Index  int64
}

func GetProblemsWithCid(cid int64) *[]ProblemDetailModel {
	var mappings []ContestProblemMappingTableModel
	if err := mysql.DBConn.Select(&mappings, "SELECT CID, `INDEX`, PID, ACCEPT, TOTAL FROM contest_problem_mapping WHERE CID = ?", cid); err != nil {
		return nil
	}
	var problems []ProblemDetailModel
	for _, value := range mappings {
		problem := GetProblemWithPid(value.PID)
		problems = append(problems, ProblemDetailModel{
			ProblemTableModel: *problem,
			Accept:            value.Accept,
			Total:             value.Total,
			Index:             value.Index,
		})
	}
	return &problems
}


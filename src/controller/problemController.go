package controller

import (
	"encoding/json"
	"main/model"
	"main/mysql"
	"main/utils"
	"math"
	"net/http"
)

type problemResponseModel struct {
	Problem model.Problem `json:"problem"`
	model.ResponseBaseModel
}

type ProblemController struct {
}

func (c *ProblemController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	var responseModel problemResponseModel
	defer func() {
		if stream, err := json.Marshal(responseModel); err == nil {
			_, _ = w.Write(stream)
		}
	}()
	responseModel.Code = model.PublicFail
	query := r.URL.Query()
	cid := utils.StringConstraint(query.Get("cid"), 0, math.MaxInt64, math.MaxInt64)
	index := utils.StringConstraint(query.Get("index"), 0, math.MaxInt64, math.MaxInt64)
	sql := "SELECT p.PID, TITLE, DESCRIPTION, DIFF, INPUT, OUTPUT, SOURCE, TIME_LIMIT, MEMORY_LIMIT FROM problem as p, contest_problem_mapping as cp WHERE cp.CID = ? AND cp.`INDEX` = ? AND cp.PID = p.PID"
	if err := mysql.DBConn.Get(&responseModel.Problem, sql, cid, index); err != nil {
		return
	}
	sql = "SELECT INPUT, OUTPUT FROM sample as s, contest_problem_mapping as cp WHERE s.PID = cp.PID AND cp.CID = ? AND cp.`INDEX` = ?"
	if err := mysql.DBConn.Select(&responseModel.Problem.Sample, sql, cid, index); err != nil {
		return
	}
	responseModel.Code = model.Success
}

func init() {
	RegisterController("/problem/", new(ProblemController))
}

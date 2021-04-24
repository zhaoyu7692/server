package controller

import (
	"encoding/json"
	"main/model"
	"main/mysql"
	"main/utils"
	"math"
	"net/http"
)

type ContestResponseModel struct {
	Problems []model.Problem `json:"problems"`
	model.ResponseBaseModel
}

type ContestController struct {
}

func (c *ContestController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	response := ContestResponseModel{}
	defer func() {
		if stream, err := json.Marshal(response); err == nil {
			_, _ = w.Write(stream)
		}
	}()
	response.Code = model.PublicFail
	query := r.URL.Query()
	cid := utils.StringConstraint(query.Get("cid"), 1, math.MaxInt64, math.MaxInt64)
	sql := "SELECT p.PID, TITLE, DIFF, ACCEPT, TOTAL FROM contest_problem_mapping as cp, problem as p WHERE CID = ? AND cp.PID = p.PID ORDER BY cp.`INDEX`"
	if err := mysql.DBConn.Select(&response.Problems, sql, cid); err != nil {
		return
	}
	response.Code = model.Success
}

func init() {
	RegisterController("/contest/", new(ContestController))
}

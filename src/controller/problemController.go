package controller

import (
	"encoding/json"
	"main/model"
	"main/mysql"
	"net/http"
	"strconv"
)

type problemResponseModel struct {
	model.Problem
	Samples []model.Example `json:"samples"`
}

type ProblemController struct {
}

func (c *ProblemController) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	if r.Method == "POST" {
		return
	}

	query := r.URL.Query()
	problemID, err := strconv.Atoi(query.Get("problemID"))
	if err != nil {
		return
	}

	var responseModel problemResponseModel
	if err := mysql.DBConn.Get(&responseModel.Problem, "SELECT PID, TITLE, DESCRIPTION, DIFF, INPUT, OUTPUT, SOURCE, TIME_LIMIT, MEMORY_LIMIT, ACCEPT, TOTAL FROM PROBLEM WHERE PID = ?", problemID); err != nil {
		return
	}
	if err := mysql.DBConn.Select(&responseModel.Samples, "SELECT INPUT, OUTPUT FROM SAMPLE WHERE PID = ?", problemID); err != nil {
		return
	}

	stream, err := json.Marshal(responseModel)
	if err != nil {
		return
	}
	_, _ = w.Write(stream)
}

func init() {
	RegisterController("/problem/", new(ProblemController))
}

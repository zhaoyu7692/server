package controller

import (
	"encoding/json"
	"main/model"
	"main/mysql"
	"main/utils"
	"math"
	"net/http"
)

type GetProblemInfoResponseModel struct {
	Problem model.Problem `json:"problem"`
	model.ResponseBaseModel
}

type GetProblemInfoController struct {

}

func (c *GetProblemInfoController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	response := GetProblemInfoResponseModel{}
	response.Code = model.PublicFail
	defer func() {
		if stream, err := json.Marshal(response); err == nil {
			_, _ = w.Write(stream)
		}
	}()
	query := r.URL.Query()
	pid := utils.StringConstraint(query.Get("pid"), -1, math.MaxInt64, -1)
	if pid == -1 {
		return
	}
	if err := mysql.DBConn.Get(&response.Problem, "SELECT PID, TITLE FROM problem WHERE PID = ?", pid); err != nil {
		response.Message = "题目不存在"
		return
	}
	response.Code = model.Success
}

func init() {
	RegisterController("/getProblemInfo/", new(GetProblemInfoController))
}

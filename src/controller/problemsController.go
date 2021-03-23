package controller

import (
	"encoding/json"
	"io/ioutil"
	"main/model"
	"main/mysql"
	"net/http"
)

type problemsResponseModel struct {
	ItemList []model.Problem `json:"item_list"`
	model.ResponsePaginationModel
}

type ProblemsController struct {

}

func (c *ProblemsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	var page model.RequestPaginationModel
	if err := json.Unmarshal(body, &page); err != nil {
		return
	}

	var problems []model.Problem
	if err := mysql.DBConn.Select(&problems, "SELECT PID, TITLE, DIFF, ACCEPT, TOTAL FROM PROBLEM"); err != nil {
		return
	}
	var responseModel problemsResponseModel
	if err := mysql.DBConn.Get(&responseModel.Total, "SELECT COUNT(PID) FROM PROBLEM"); err != nil {
		return
	}
	responseModel.ItemList = problems
	responseModel.Size = page.Size

	stream, err := json.Marshal(responseModel)
	if err != nil {
		return
	}
	_, _ = w.Write(stream)
}

func init() {
	RegisterController("/problems/", new(ProblemsController))
	//runtime.ReadMemStats()
}

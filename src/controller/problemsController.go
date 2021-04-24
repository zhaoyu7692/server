package controller

import (
	"encoding/json"
	"io/ioutil"
	"main/model"
	"main/mysql"
	"main/utils"
	"net/http"
)

type problemsResponseModel struct {
	Problems struct {
		ItemList []model.Problem `json:"item_list"`
		model.ResponsePaginationModel
	} `json:"problems"`
	model.ResponseBaseModel
}

type ProblemsController struct {
}

func (c *ProblemsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	responseModel := problemsResponseModel{}
	defer func() {
		if stream, err := json.Marshal(responseModel); err == nil {
			_, _ = w.Write(stream)
		}
	}()
	responseModel.Code = model.PublicFail
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	var page model.RequestPaginationModel
	if err := json.Unmarshal(body, &page); err != nil {
		return
	}

	offset := utils.Max((page.Page-1)*page.Size, 0)
	sql := "SELECT p.PID, TITLE, DIFF, ACCEPT, TOTAL FROM problem as p, contest_problem_mapping as cp WHERE CID = 0 AND p.PID = cp.PID ORDER BY `INDEX` LIMIT ?, ?"
	if err := mysql.DBConn.Select(&responseModel.Problems.ItemList, sql, offset, page.Size); err != nil {
		return
	}

	if err := mysql.DBConn.Get(&responseModel.Problems.Total, "SELECT COUNT(PID) FROM contest_problem_mapping WHERE CID = 0"); err != nil {
		return
	}
	responseModel.Problems.Size = page.Size
	responseModel.Code = model.Success
}

func init() {
	RegisterController("/problems/", new(ProblemsController))
}

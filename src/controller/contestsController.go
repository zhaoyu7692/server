package controller

import (
	"encoding/json"
	"main/model"
	"main/mysql"
	"main/utils"
	"math"
	"net/http"
)

type contestRequestModel struct {
	model.RequestPaginationModel
}

type contestResponseModel struct {
	Contests struct {
		ItemList []model.Contest `json:"item_list"`
		model.ResponsePaginationModel
	} `json:"contests"`
	model.ResponseBaseModel
}

type ContestsController struct {
}

func (c *ContestsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	response := contestResponseModel{}
	defer func() {
		if stream, err := json.Marshal(response); err == nil {
			_, _ = w.Write(stream)
		}
	}()
	query := r.URL.Query()
	page := utils.StringConstraint(query.Get("page"), 0, math.MaxInt64, 20)
	size := utils.StringConstraint(query.Get("size"), 20, 20, 20)
	offset := utils.Max((page-1)*size, 0)
	if err := mysql.DBConn.Select(&response.Contests.ItemList, "SELECT CID, TITLE, BEGIN_TIME, DURATION, REGISTER_COUNT FROM contest ORDER BY CID DESC LIMIT ?, ?", offset, size); err != nil {
		return
	}
	if err := mysql.DBConn.Get(&response.Contests.Total, "SELECT  COUNT(*) FROM contest"); err != nil {
		return
	}
	response.Contests.Size = size
	response.Code = model.Success
}

func init() {
	RegisterController("/contests/", new(ContestsController))
}

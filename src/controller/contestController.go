package controller

import (
	"encoding/json"
	"main/dao"
	"main/model"
	"main/mysql"
	"main/utils"
	"math"
	"net/http"
)

func init() {
	RegisterHandler("/contest/", contest)
	RegisterHandler("/contests/", contests)

}

type ContestResponseModel struct {
	Problems []struct {
		model.Problem
		Index int64 `json:"index" db:"INDEX"`
	} `json:"problems"`
	Contest dao.ContestTableModel `json:"contest"`
	model.ResponseBaseModel
}

func contest(w http.ResponseWriter, r *http.Request) {
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
	sql := "SELECT `INDEX`, TITLE, DIFF, ACCEPT, TOTAL FROM contest_problem_mapping as cp, problem as p WHERE CID = ? AND cp.PID = p.PID ORDER BY `INDEX`"
	if err := mysql.DBConn.Select(&response.Problems, sql, cid); err != nil {
		return
	}
	response.Contest = *dao.GetContestWithCid(cid)
	response.Code = model.Success
}

type contestsResponseModel struct {
	Contests struct {
		ItemList []model.Contest `json:"item_list"`
		model.ResponsePaginationModel
	} `json:"contests"`
	model.ResponseBaseModel
}

func contests(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	response := contestsResponseModel{}
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

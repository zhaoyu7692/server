package controller

import (
	"encoding/json"
	"main/model"
	"main/mysql"
	"main/utils"
	"math"
	"net/http"
	"time"
)

type contestRankResponseModel struct {
	ContestRank []struct {
		model.User
		Problem []struct {
			Index      int64      `json:"index" db:"INDEX"`
			TryCount   int64      `json:"try_count" db:"TRY_COUNT"`
			AcceptTime *time.Time `json:"accept_time,omitempty" db:"ACCEPT_TIME"`
		} `json:"problem"`
	} `json:"contest_rank"`
	Contest model.Contest `json:"contest"`
	model.ResponseBaseModel
}

type ContestRankController struct {
}

func (c *ContestRankController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	response := contestRankResponseModel{}
	defer func() {
		if stream, err := json.Marshal(response); err == nil {
			_, _ = w.Write(stream)
		}
	}()
	response.Code = model.PublicFail
	query := r.URL.Query()
	cid := utils.StringConstraint(query.Get("cid"), 1, math.MaxInt64, math.MaxInt64)
	if err := mysql.DBConn.Get(&response.Contest, "SELECT CID, TITLE, BEGIN_TIME, DURATION FROM contest WHERE  CID = ?", cid); err != nil {
		return
	}
	if err := mysql.DBConn.Select(&response.Contest.Indexes, "SELECT `INDEX` FROM contest_problem_mapping WHERE CID = ? ORDER BY `INDEX`", cid); err != nil {
		return
	}
	if err := mysql.DBConn.Select(&response.ContestRank, "SELECT UID FROM contest_rank WHERE CID = ?", cid); err != nil {
		return
	}
	sql := "SELECT `INDEX`, TRY_COUNT, ACCEPT_TIME FROM contest_rank WHERE CID = ? AND UID = ? ORDER BY `INDEX`"
	for i := 0; i < len(response.ContestRank); i++ {
		_ = mysql.DBConn.Select(&response.ContestRank[i].Problem, sql, cid, response.ContestRank[i].Uid)
	}
	response.Code = model.Success
}

func init() {
	RegisterController("/contestRank/", new(ContestRankController))
}

package controller

import (
	"encoding/json"
	"fmt"
	"main/dao"
	"main/model"
	"main/mysql"
	"main/redispool"
	"main/utils"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func init() {
	RegisterHandler("/contest/", contest)
	RegisterHandler("/contests/", contests)
	RegisterHandler("/contestRank/", contestRank)
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

type ProblemRankModel struct {
	Index          int64 `json:"index"`
	TryCount       int64 `json:"try_count"`
	AcceptDuration int64 `json:"accept_time,omitempty"`
}

type contestRankResponseModel struct {
	ContestRank []struct {
		//Username    string `json:"username"`
		AcceptCount int64 `json:"accept_count"`
		Penalty     int64 `json:"penalty"`

		model.User
		Problem []ProblemRankModel `json:"problem"`
	} `json:"contest_rank"`
	Contest model.Contest `json:"contest"`
	model.ResponseBaseModel
}


func contestRank(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	response := contestRankResponseModel{}
	response.Code = model.PublicFail
	query := r.URL.Query()
	cid := utils.StringConstraint(query.Get("cid"), 1, math.MaxInt64, math.MaxInt64)
	if cid == math.MaxInt64 {
		return
	}
	// TODO: redis
	//if reply, err := redis.Bytes(redispool.Get().Do("GET", fmt.Sprintf("contest_rank_key_cid_%d", cid))); err == nil {
	//	_, _ = w.Write(reply)
	//	return
	//}
	defer func() {
		if stream, err := json.Marshal(response); err == nil {
			_, _ = redispool.Get().Do("SET", fmt.Sprintf("contest_rank_key_cid_%d", cid), stream)
			_, _ = w.Write(stream)
		}
	}()
	if err := mysql.DBConn.Get(&response.Contest, "SELECT CID, TITLE, BEGIN_TIME, DURATION FROM contest WHERE  CID = ?", cid); err != nil {
		return
	}
	if err := mysql.DBConn.Select(&response.Contest.Indexes, "SELECT `INDEX` FROM contest_problem_mapping WHERE CID = ? ORDER BY `INDEX`", cid); err != nil {
		return
	}
	if err := mysql.DBConn.Select(&response.ContestRank, "SELECT DISTINCT UID FROM submit WHERE CID = ?", cid); err != nil {
		return
	}

	duration, err := time.ParseDuration(strconv.FormatInt(response.Contest.Duration, 10) + "s")
	if err != nil {
		return
	}
	// 区分用户
	for i := 0; i < len(response.ContestRank); i++ {
		if err := mysql.DBConn.Get(&response.ContestRank[i].Username, "SELECT USERNAME FROM user WHERE UID = ?", response.ContestRank[i].Uid); err != nil {
			return
		}
		response.ContestRank[i].Problem = []ProblemRankModel{}
		// 区分题目
		for j := 0; j < len(response.Contest.Indexes); j++ {
			var submits []model.Submit
			if err := mysql.DBConn.Select(&submits, "SELECT STATUS, SUBMIT_TIME FROM submit WHERE UID = ? AND CID = ? AND `INDEX` = ? ORDER BY RID", response.ContestRank[i].Uid, cid, response.Contest.Indexes[j]); err != nil {
				return
			}
			problemRank := ProblemRankModel{
				Index:          response.Contest.Indexes[j],
				AcceptDuration: -1,
			}
			tryCount := int64(0)
			for k := 0; k < len(submits); k++ {
				submitTime := submits[k].SubmitTime
				if submitTime.Before(*response.Contest.BeginTime) || submitTime.After(response.Contest.BeginTime.Add(duration)) {
					continue
				}
				switch submits[k].Status {
				case model.JudgeStatusCompilationError,
					model.JudgeStatusCompilationTimeLimitExceeded,
					model.JudgeStatusTimeLimitExceeded,
					model.JudgeStatusMemoryLimitExceeded,
					model.JudgeStatusOutputLimitExceeded,
					model.JudgeStatusRuntimeError,
					model.JudgeStatusPresentationError,
					model.JudgeStatusWrongAnswer,
					model.JudgeStatusAccept:
					tryCount++
				}

				if submits[k].Status == model.JudgeStatusAccept {
					problemRank.AcceptDuration = int64(math.Floor(submits[k].SubmitTime.Sub(*response.Contest.BeginTime).Minutes()))
					break
				}
			}
			problemRank.TryCount = tryCount
			if problemRank.AcceptDuration > 0 {
				response.ContestRank[i].Penalty += problemRank.AcceptDuration + (tryCount-1)*20
				response.ContestRank[i].AcceptCount++
			}
			response.ContestRank[i].Problem = append(response.ContestRank[i].Problem, problemRank)
		}
	}
	sort.Slice(response.ContestRank, func(i, j int) bool {
		if response.ContestRank[i].AcceptCount < response.ContestRank[j].AcceptCount {
			return false
		} else if response.ContestRank[i].AcceptCount == response.ContestRank[j].AcceptCount && response.ContestRank[i].Penalty >= response.ContestRank[j].Penalty {
			return false
		}
		return true
	})
	response.Code = model.Success
}

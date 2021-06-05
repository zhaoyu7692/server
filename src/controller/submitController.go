package controller

import (
	"encoding/json"
	"io/ioutil"
	"main/dao"
	"main/model"
	"main/mysql"
	"main/service"
	"net/http"
	"strconv"
	"time"
)

type submitResponseModel struct {
	model.ResponseBaseModel
}

type SubmitController struct {
}

func (c *SubmitController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	var response submitResponseModel
	// 返回
	defer func() {
		if stream, err := json.Marshal(response); err == nil {
			_, _ = w.Write(stream)
		}
	}()
	response.Code = model.PublicFail
	// 获取请求数据
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	var submit model.Submit
	if err := json.Unmarshal(body, &submit); err != nil {
		return
	}
	// 鉴权
	authStatus := service.AuthCheck(submit.Uid, submit.Token)
	//if authStatus != service.Authority {
	//	return
	//}
	switch authStatus {
	case service.UnAuthority:
		{
			response.Message = "请登录后提交题目"
		}
	case service.SessionOverdue:
		{
			response.Message = "会话过期，请登陆后提交题目"
		}
	}
	if authStatus != service.Authority && authStatus != service.AuthorityAdmin {
		response.Code = model.JumpLogin
		return
	}

	var contest model.Contest
	// 比赛时间校验
	now := time.Now()
	if submit.Cid > 0 {
		if err := mysql.DBConn.Get(&contest, "SELECT BEGIN_TIME, DURATION FROM contest WHERE CID = ?", submit.Cid); err != nil {
			return
		}
		duration, err := time.ParseDuration(strconv.FormatInt(contest.Duration, 10) + "s")
		if err != nil {
			return
		}
		if now.Before(*contest.BeginTime) || now.After(contest.BeginTime.Add(duration)) {
			response.Message = "不在比赛时间内，无法提交"
			return
		}
	}
	var problem model.JudgeSubmitModel
	sql := "SELECT p.PID, p.TIME_LIMIT, p.MEMORY_LIMIT FROM problem as p, contest_problem_mapping as cp WHERE CID = ? AND `INDEX` = ? AND p.PID = cp.PID"
	if err := mysql.DBConn.Get(&problem, sql, submit.Cid, submit.Index); err != nil {
		return
	}
	// 提交
	sql = "INSERT INTO submit (CID, `INDEX`, UID, CODE, STATUS, LANGUAGE, SUBMIT_TIME) values (?,?,?,?,?,?,?)"
	if result, err := mysql.DBConn.Exec(sql, submit.Cid, submit.Index, submit.Uid, submit.Code, 0, submit.Language, now); err == nil {
		if rid, err := result.LastInsertId(); err == nil {
			// 写入待测题库
			problem.Rid = rid
			problem.Cid = submit.Cid
			problem.Index = submit.Index
			problem.Code = submit.Code
			problem.Status = 0
			problem.Language = submit.Language
			service.StashSubmit(&problem)
			response.Code = model.Success
		}
	}
}

type getSubmitDetailRequestModel struct {
	Uid int64 `json:"uid"`
	Rid int64 `json:"rid"`
}

type getSubmitDetailResponseModel struct {
	SourceCode         string `json:"source_code"`
	CompilationMessage string `json:"compilation_message"`
	model.ResponseBaseModel
}

func getSubmitDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	responseModel := getSubmitDetailResponseModel{}
	responseModel.Code = model.PublicFail
	defer func() {
		if stream, err := json.Marshal(responseModel); err == nil {
			_, _ = w.Write(stream)
		}
	}()
	requestModel := getSubmitDetailRequestModel{}
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal(body, &requestModel); err != nil {
			return
		}
	} else {
		return
	}
	if submit := dao.GetSubmitWithRid(requestModel.Rid); submit != nil {
		if submit.Uid != requestModel.Uid {
			return
		}
		responseModel.SourceCode = submit.Code
		responseModel.CompilationMessage = ""
		if submit.CompilationMessage != nil {
			responseModel.CompilationMessage = *submit.CompilationMessage
		}
		responseModel.Code = model.Success
	}
}

func init() {
	RegisterHandler("/getSubmitDetail/", getSubmitDetail)
	RegisterController("/submit/", new(SubmitController))
}

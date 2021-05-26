package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/dao"
	"main/model"
	"main/mysql"
	"main/service"
	"main/utils"
	"math"
	"net/http"
	"os"
	"sync"
	"time"
)

func init() {
	//RegisterHandler("/getProblem/", getProblem)
	RegisterHandler("/createProblem/", createProblem)
	//RegisterHandler("/updateProblem/", updateProblem)
	RegisterHandler("/deleteProblem/", deleteProblem)
	RegisterHandler("/getContest/", getContest)
	RegisterHandler("/createContest/", createContest)
	RegisterHandler("/updateContest/", updateContest)
	RegisterHandler("/deleteContest/", deleteContest)
}

//type getProblemResponseModel struct {
//	Problem   *dao.ProblemTableModel           `json:"problem"`
//	Samples   *[]dao.SampleTableModel          `json:"samples"`
//	Testcases *[]dao.TestcaseMappingTableModel `json:"testcases"`
//	model.ResponseBaseModel
//}

//func getProblem(w http.ResponseWriter, r *http.Request) {
//	if r.Method != http.MethodGet {
//		return
//	}
//	responseModel := getProblemResponseModel{}
//	responseModel.Code = model.PublicFail
//	defer func() {
//		if stream, err := json.Marshal(responseModel); err == nil {
//			_, _ = w.Write(stream)
//		}
//	}()
//	query := r.URL.Query()
//	pid := utils.StringConstraint(query.Get("pid"), 1, math.MaxInt64, math.MaxInt64)
//	if pid == math.MaxInt64 {
//		return
//	}
//	responseModel.Problem = dao.GetProblemWithPid(pid)
//	responseModel.Samples = dao.GetSamplesWithPid(pid)
//	responseModel.Testcases = dao.GetTestCasesWithPid(pid)
//	responseModel.Code = model.Success
//}

type createProblemRequestModel struct {
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	Difficulty   int64           `json:"difficulty"`
	Input        string          `json:"input"`
	Output       string          `json:"output"`
	Source       string          `json:"source"`
	TimeLimit    int64           `json:"time_limit"`
	MemoryLimit  int64           `json:"memory_limit"`
	FilenameList []string        `json:"filename_list"`
	Samples      []model.Example `json:"samples"`
	User         model.User      `json:"user"`
}

type createProblemResponseModel struct {
	Pid int64 `json:"pid,omitempty"`
	model.ResponseBaseModel
}

var problemLock sync.Mutex

func createProblem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	responseModel := createProblemResponseModel{}
	responseModel.Code = model.PublicFail
	defer func() {
		if data, err := json.Marshal(responseModel); err == nil {
			_, _ = w.Write(data)
		}
	}()

	requestModel := createProblemRequestModel{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	if err = json.Unmarshal(body, &requestModel); err != nil {
		return
	}
	if service.AuthCheck(requestModel.User.Uid, requestModel.User.Token) != service.AuthorityAdmin {
		responseModel.Message = "没有权限进行该操作"
		return
	}
	problemLock.Lock()
	defer problemLock.Unlock()
	tx, err := mysql.DBConn.Begin()
	if err != nil {
		return
	}
	defer func() {
		if responseModel.Code != model.Success {
			_ = tx.Rollback()
		}
	}()
	// create problem
	var pid int64
	if err := mysql.DBConn.Get(&pid, "SELECT MAX(PID) FROM problem"); err != nil {
		return
	}
	pid++
	sql := "INSERT INTO problem (PID, TITLE, DESCRIPTION, DIFF, INPUT, OUTPUT, SOURCE, TIME_LIMIT, MEMORY_LIMIT) VALUES (?,?,?,?,?,?,?,?,?)"
	if _, err = tx.Exec(sql, pid, requestModel.Title, requestModel.Description, requestModel.Difficulty, requestModel.Input, requestModel.Output, requestModel.Source, requestModel.TimeLimit, requestModel.MemoryLimit); err != nil {
		return
	}

	// examples
	for index, sample := range requestModel.Samples {
		if _, err = tx.Exec("INSERT INTO sample (PID, SID, INPUT, OUTPUT) VALUES (?,?,?,?)", pid, index+1, sample.Input, sample.Output); err != nil {
			return
		}
	}

	// add to contest_problem_mapping
	if _, err = tx.Exec("INSERT INTO contest_problem_mapping (CID, `INDEX`, PID) VALUES (?,?,?)", 0, pid, pid); err != nil {
		return
	}

	// mapping file
	_, err = os.Lstat(fmt.Sprintf("%s%d", utils.GlobalConfig.Path.Data, pid))
	if os.IsNotExist(err) {
		if err := os.MkdirAll(fmt.Sprintf("%s%d", utils.GlobalConfig.Path.Data, pid), os.ModePerm); err != nil {
			return
		}
	}
	for _, hash := range requestModel.FilenameList {
		resourceModel := model.ResourceMappingModel{}
		if err = mysql.DBConn.Get(&resourceModel, "SELECT * FROM resource_mapping WHERE SHA_KEY = ?", hash); err != nil {
			return
		}

		realPath := fmt.Sprintf("%s%d/%s", utils.GlobalConfig.Path.Data, pid, resourceModel.Filename)
		data, err := ioutil.ReadFile(resourceModel.Path)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		if err := ioutil.WriteFile(realPath, data, os.ModePerm); err != nil {
			w.WriteHeader(500)
			return
		}

		if _, err = tx.Exec("INSERT INTO testcase_mapping (PID, FILENAME, `KEY`, PATH) VALUES (?,?,?,?)", pid, resourceModel.Filename, hash, realPath); err != nil {
			return
		}
	}

	if err = tx.Commit(); err != nil {
		return
	}
	responseModel.Pid = pid
	responseModel.Code = model.Success
}

//func updateProblem(w http.ResponseWriter, r *http.Request) {
//
//}

type deleteProblemRequestModel struct {
	Pid  int64      `json:"pid"`
	User model.User `json:"user"`
}

type deleteProblemResponseModel struct {
	model.ResponseBaseModel
}

func deleteProblem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	responseModel := deleteProblemResponseModel{}
	responseModel.Code = model.PublicFail
	defer func() {
		if data, err := json.Marshal(responseModel); err == nil {
			_, _ = w.Write(data)
		}
	}()

	requestModel := deleteProblemRequestModel{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	if err = json.Unmarshal(body, &requestModel); err != nil {
		return
	}
	if service.AuthCheck(requestModel.User.Uid, requestModel.User.Token) != service.AuthorityAdmin {
		responseModel.Message = "没有权限进行该操作"
		return
	}
	problemLock.Lock()
	defer problemLock.Unlock()
	tx, err := mysql.DBConn.Begin()
	if err != nil {
		return
	}
	defer func() {
		if responseModel.Code != model.Success {
			_ = tx.Rollback()
		}
	}()
	// delete problem
	if _, err := tx.Exec("DELETE FROM problem WHERE PID = ?", requestModel.Pid); err != nil {
		return
	}
	if _, err := tx.Exec("DELETE  FROM contest_problem_mapping WHERE PID = ?", requestModel.Pid); err != nil {
		return
	}
	if _, err := tx.Exec("DELETE FROM testcase_mapping WHERE PID = ?", requestModel.Pid); err != nil {
		return
	}
	if err = tx.Commit(); err != nil {
		return
	}
	responseModel.Code = model.Success
}

type getContestResponseModel struct {
	Contest  *dao.ContestTableModel    `json:"contest"`
	Problems *[]dao.ProblemDetailModel `json:"problems"`
	model.ResponseBaseModel
}

func getContest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}
	responseModel := getContestResponseModel{}
	responseModel.Code = model.PublicFail
	defer func() {
		if stream, err := json.Marshal(responseModel); err == nil {
			_, _ = w.Write(stream)
		}
	}()
	query := r.URL.Query()
	cid := utils.StringConstraint(query.Get("cid"), 1, math.MaxInt64, math.MaxInt64)
	if cid == math.MaxInt64 {
		return
	}
	responseModel.Contest = dao.GetContestWithCid(cid)
	responseModel.Problems = dao.GetProblemsWithCid(cid)
	responseModel.Code = model.Success
}

type createContestRequestModel struct {
	Title     string    `json:"title"`
	StartTime time.Time `json:"start_time"`
	Duration  int64     `json:"duration"`
	EndTime   time.Time `json:"end_time"`
	Problems  []struct {
		Index int64 `json:"index"`
		Pid   int64 `json:"pid"`
	} `json:"problems"`
	User model.User `json:"user"`
}

type createContestResponseModel struct {
	model.ResponseBaseModel
}

var contestLock sync.Mutex

func createContest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	response := createContestResponseModel{}
	response.Code = model.PublicFail
	defer func() {
		if stream, err := json.Marshal(response); err == nil {
			_, _ = w.Write(stream)
		}
	}()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	requestModel := createContestRequestModel{}
	if err := json.Unmarshal(body, &requestModel); err != nil {
		return
	}
	if service.AuthCheck(requestModel.User.Uid, requestModel.User.Token) != service.AuthorityAdmin {
		return
	}
	fmt.Println(requestModel)
	contestLock.Lock()
	defer contestLock.Unlock()
	tx, err := mysql.DBConn.Begin()
	if err != nil {
		return
	}
	defer func() {
		if response.Code != model.Success {
			_ = tx.Rollback()
		}
	}()
	result, err := tx.Exec("INSERT INTO contest (TITLE, BEGIN_TIME, DURATION) VALUES (?,?,?)", requestModel.Title, requestModel.StartTime, requestModel.Duration)
	if err != nil {
		return
	}
	cid, err := result.LastInsertId()
	if err != nil {
		return
	}
	for i := 0; i < len(requestModel.Problems); i++ {
		problem := requestModel.Problems[i]
		if _, err = tx.Exec("INSERT INTO contest_problem_mapping (CID, `INDEX`, PID) VALUES (?,?,?)", cid, problem.Index, problem.Pid); err != nil {
			return
		}
	}
	if err := tx.Commit(); err != nil {
		return
	}
	response.Code = model.Success
}

type updateContestRequestModel struct {
	Cid       int64     `json:"cid"`
	Title     string    `json:"title"`
	StartTime time.Time `json:"start_time"`
	Duration  int64     `json:"duration"`
	EndTime   time.Time `json:"end_time"`
	Problems  []struct {
		Index int64 `json:"index"`
		Pid   int64 `json:"pid"`
	} `json:"problems"`
	User model.User `json:"user"`
}

type updateContestResponseModel struct {
	model.ResponseBaseModel
}

func updateContest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	response := updateContestResponseModel{}
	response.Code = model.PublicFail
	defer func() {
		if stream, err := json.Marshal(response); err == nil {
			_, _ = w.Write(stream)
		}
	}()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	requestModel := updateContestRequestModel{}
	if err := json.Unmarshal(body, &requestModel); err != nil {
		return
	}
	if requestModel.Cid <= 0 || service.AuthCheck(requestModel.User.Uid, requestModel.User.Token) != service.AuthorityAdmin {
		return
	}
	fmt.Println(requestModel)
	contestLock.Lock()
	defer contestLock.Unlock()
	tx, err := mysql.DBConn.Begin()
	if err != nil {
		return
	}
	defer func() {
		if response.Code != model.Success {
			_ = tx.Rollback()
		}
	}()
	_, err = tx.Exec("UPDATE contest SET TITLE = ?, BEGIN_TIME = ?, DURATION = ? WHERE CID = ?", requestModel.Title, requestModel.StartTime, requestModel.Duration, requestModel.Cid)
	if err != nil {
		return
	}
	if _, err = tx.Exec("DELETE FROM contest_problem_mapping WHERE CID = ?", requestModel.Cid); err != nil {
		return
	}
	for i := 0; i < len(requestModel.Problems); i++ {
		problem := requestModel.Problems[i]
		if _, err = tx.Exec("INSERT INTO contest_problem_mapping (CID, `INDEX`, PID) VALUES (?,?,?)", requestModel.Cid, problem.Index, problem.Pid); err != nil {
			return
		}
	}
	if err := tx.Commit(); err != nil {
		return
	}
	response.Code = model.Success
}

type deleteContestRequestModel struct {
	Cid  int64      `json:"cid"`
	User model.User `json:"user"`
}

type deleteContestResponseModel struct {
	model.ResponseBaseModel
}

func deleteContest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}
	response := deleteContestResponseModel{}
	response.Code = model.PublicFail
	defer func() {
		if stream, err := json.Marshal(response); err == nil {
			_, _ = w.Write(stream)
		}
	}()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	requestModel := deleteContestRequestModel{}
	if err := json.Unmarshal(body, &requestModel); err != nil {
		return
	}
	if requestModel.Cid <= 0 || service.AuthCheck(requestModel.User.Uid, requestModel.User.Token) != service.AuthorityAdmin {
		return
	}
	contestLock.Lock()
	defer contestLock.Unlock()
	fmt.Println(requestModel)
	tx, err := mysql.DBConn.Begin()
	if err != nil {
		return
	}
	defer func() {
		if response.Code != model.Success {
			_ = tx.Rollback()
		}
	}()
	if _, err = tx.Exec("DELETE FROM contest WHERE CID = ?", requestModel.Cid); err != nil {
		return
	}
	if _, err = tx.Exec("DELETE FROM contest_problem_mapping WHERE CID = ?", requestModel.Cid); err != nil {
		return
	}

	if err := tx.Commit(); err != nil {
		return
	}
	response.Code = model.Success
}

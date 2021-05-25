package controller

//type CreateContestController struct {
//}
//
//type CreateContestRequestModel struct {
//	Title     string    `json:"title"`
//	StartTime time.Time `json:"start_time"`
//	Duration  int64     `json:"duration"`
//	EndTime   time.Time `json:"end_time"`
//	Problems  []struct {
//		Index int64 `json:"index"`
//		Pid   int64 `json:"pid"`
//	} `json:"problems"`
//}
//
//type CreateContestResponseModel struct {
//	model.ResponseBaseModel
//}
//
//func (c *CreateContestController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	if r.Method != "POST" {
//		return
//	}
//	response := CreateContestResponseModel{}
//	response.Code = model.PublicFail
//	defer func() {
//		if stream, err := json.Marshal(response); err == nil {
//			_, _ = w.Write(stream)
//		}
//	}()
//	body, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		return
//	}
//	requestModel := CreateContestRequestModel{}
//	if err := json.Unmarshal(body, &requestModel); err != nil {
//		return
//	}
//
//	fmt.Println(requestModel)
//	tx, err := mysql.DBConn.Begin()
//	if err != nil {
//		return
//	}
//	result, err := tx.Exec("INSERT INTO contest (TITLE, BEGIN_TIME, DURATION) VALUES (?,?,?)", requestModel.Title, requestModel.StartTime, requestModel.Duration)
//	if err != nil {
//		_ = tx.Rollback()
//		return
//	}
//	cid, err := result.LastInsertId()
//	if err != nil {
//		_ = tx.Rollback()
//		return
//	}
//	for i := 0; i < len(requestModel.Problems); i++ {
//		problem := requestModel.Problems[i]
//		if _, err = tx.Exec("INSERT INTO contest_problem_mapping (CID, `INDEX`, PID) VALUES (?,?,?)", cid, problem.Index, problem.Pid); err != nil {
//			_ = tx.Rollback()
//			return
//		}
//	}
//	if err := tx.Commit(); err != nil {
//		return
//	}
//	response.Code = model.Success
//}
//
//func init() {
//	RegisterController("/createContest/", new(CreateContestController))
//}

package controller

//type CreateProblemRequestModel struct {
//	Title        string          `json:"title"`
//	Description  string          `json:"description"`
//	Difficulty   int64           `json:"difficulty"`
//	Input        string          `json:"input"`
//	Output       string          `json:"output"`
//	Source       string          `json:"source"`
//	TimeLimit    int64           `json:"time_limit"`
//	MemoryLimit  int64           `json:"memory_limit"`
//	FilenameList []string        `json:"filename_list"`
//	Samples      []model.Example `json:"samples"`
//}

//type CreateProblemResponseModel struct {
//	Pid int64 `json:"pid,omitempty"`
//	model.ResponseBaseModel
//}
//
//type CreateProblemController struct {
//}

//var problemLock sync.Mutex

//func (c *CreateProblemController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	if r.Method != "POST" {
//		return
//	}
//	responseModel := CreateProblemResponseModel{}
//	responseModel.Code = model.PublicFail
//	defer func() {
//		if data, err := json.Marshal(responseModel); err == nil {
//			_, _ = w.Write(data)
//		}
//	}()
//
//	requestModel := CreateProblemRequestModel{}
//	body, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		return
//	}
//	if err = json.Unmarshal(body, &requestModel); err != nil {
//		return
//	}
//	tx, err := mysql.DBConn.Begin()
//	if err != nil {
//		return
//	}
//	// create problem
//	problemLock.Lock()
//	defer problemLock.Unlock()
//	var pid int64
//	if err := mysql.DBConn.Get(&pid, "SELECT MAX(PID) FROM problem"); err != nil {
//		return
//	}
//	pid++
//	sql := "INSERT INTO problem (TITLE, DESCRIPTION, DIFF, INPUT, OUTPUT, SOURCE, TIME_LIMIT, MEMORY_LIMIT) VALUES (?,?,?,?,?,?,?,?)"
//	if _, err = tx.Exec(sql, requestModel.Title, requestModel.Description, requestModel.Difficulty, requestModel.Input, requestModel.Output, requestModel.Source, requestModel.TimeLimit, requestModel.MemoryLimit); err != nil {
//		_ = tx.Rollback()
//		return
//	}
//
//	// examples
//	for index, sample := range requestModel.Samples {
//		if _, err = tx.Exec("INSERT INTO sample (PID, SID, INPUT, OUTPUT) VALUES (?,?,?,?)", pid, index+1, sample.Input, sample.Output); err != nil {
//			_ = tx.Rollback()
//			return
//		}
//	}
//
//	// add to contest_problem_mapping
//	if _, err = tx.Exec("INSERT INTO contest_problem_mapping (CID, `INDEX`, PID) VALUES (?,?,?)", 0, pid, pid); err != nil {
//		_ = tx.Rollback()
//		return
//	}
//
//	// mapping file
//	_, err = os.Lstat(fmt.Sprintf("%s%d", utils.GlobalConfig.Path.Data, pid))
//	if os.IsNotExist(err) {
//		if err := os.MkdirAll(fmt.Sprintf("%s%d", utils.GlobalConfig.Path.Data, pid), os.ModePerm); err != nil {
//			_ = tx.Rollback()
//			return
//		}
//	}
//	for _, hash := range requestModel.FilenameList {
//		resourceModel := model.ResourceMappingModel{}
//		if err = mysql.DBConn.Get(&resourceModel, "SELECT * FROM resource_mapping WHERE SHA_KEY = ?", hash); err != nil {
//			_ = tx.Rollback()
//			return
//		}
//
//		realPath := fmt.Sprintf("%s%d/%s", utils.GlobalConfig.Path.Data, pid, resourceModel.Filename)
//		data, err := ioutil.ReadFile(resourceModel.Path)
//		if err != nil {
//			w.WriteHeader(500)
//			_ = tx.Rollback()
//			return
//		}
//		if err := ioutil.WriteFile(realPath, data, os.ModePerm); err != nil {
//			w.WriteHeader(500)
//			_ = tx.Rollback()
//			return
//		}
//
//		if _, err = tx.Exec("INSERT INTO testcase_mapping (PID, FILENAME, `KEY`, PATH) VALUES (?,?,?,?)", pid, resourceModel.Filename, hash, realPath); err != nil {
//			_ = tx.Rollback()
//			return
//		}
//	}
//
//	if err = tx.Commit(); err != nil {
//		_ = tx.Rollback()
//		return
//	}
//	responseModel.Pid = pid
//	responseModel.Code = model.Success
//}
//
//func init() {
//	RegisterController("/createProblem/", new(CreateProblemController))
//}

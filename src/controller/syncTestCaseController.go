package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/mysql"
	"main/utils"
	"math"
	"net/http"
)

type syncTestCaseRequestModel struct {
	Filenames []string `json:"filenames"`
	Pid       int64    `json:"pid"`
}

type syncTestCaseResponseModel struct {
	Filenames       []string `json:"filenames"`
	RemoveFilenames []string `json:"remove_filenames"`
}

type CheckTestCaseController struct {
}

func (c *CheckTestCaseController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	responseModel := syncTestCaseResponseModel{}
	defer func() {
		if data, err := json.Marshal(responseModel); err == nil {
			_, _ = w.Write(data)
		}
	}()
	requestModel := syncTestCaseRequestModel{}
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		_ = json.Unmarshal(body, &requestModel)
	}

	var filenames []string
	if err := mysql.DBConn.Select(&filenames, "SELECT FILENAME FROM testcase_mapping WHERE PID = ?", requestModel.Pid); err != nil {
		return
	}

	filenameMap := make(map[string]bool)
	for _, filename := range filenames {
		filenameMap[filename] = true
	}

	// filename that will be removed
	removeFiles := make(map[string]bool)
	for _, filename := range requestModel.Filenames {
		removeFiles[filename] = true
	}
	for filename := range filenameMap {
		delete(removeFiles, filename)
	}
	for filename := range removeFiles {
		responseModel.RemoveFilenames = append(responseModel.RemoveFilenames, filename)
	}

	// filename that need to be synced
	for _, filename := range requestModel.Filenames {
		delete(filenameMap, filename)
	}
	for fileName := range filenameMap {
		responseModel.Filenames = append(responseModel.Filenames, fileName)
	}
}

type DownloadTestCaseController struct {
}

func (c *DownloadTestCaseController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	query := r.URL.Query()
	pid := utils.StringConstraint(query.Get("pid"), 1, math.MaxInt64, math.MaxInt64)
	filename := query.Get("filename")
	data, err := ioutil.ReadFile(fmt.Sprintf("%s%d\\%s", utils.GlobalConfig.Path.Data, pid, filename))
	if err == nil {
		_, err = w.Write(data)
		if err == nil {
			return
		}
	}
	w.WriteHeader(404)
}

func init() {
	RegisterController("/downloadTestCase/", new(DownloadTestCaseController))
	RegisterController("/checkTestCase/", new(CheckTestCaseController))
}

package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/utils"
	"math"
	"net/http"
	"regexp"
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

	testCases, err := ioutil.ReadDir(fmt.Sprintf("%s%d", utils.GlobalConfig.Path.Data, requestModel.Pid))
	if err == nil {
		files := make(map[string]bool)
		removeFiles := make(map[string]bool)
		inputFileRegex := regexp.MustCompile("^.+\\.in$")
		outputFileRegex := regexp.MustCompile("^.+\\.out$")
		for _, file := range testCases {
			if inputFileRegex.MatchString(file.Name()) || outputFileRegex.MatchString(file.Name()) {
				files[file.Name()] = true
			}
		}

		// filename that will be removed
		for _, filename := range requestModel.Filenames {
			removeFiles[filename] = true
		}
		for filename := range files {
			delete(removeFiles, filename)
		}
		for filename := range removeFiles {
			responseModel.RemoveFilenames = append(responseModel.RemoveFilenames, filename)
		}

		// filename that need to be synced
		for _, filename := range requestModel.Filenames {
			delete(files, filename)
		}
		for fileName := range files {
			responseModel.Filenames = append(responseModel.Filenames, fileName)
		}
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

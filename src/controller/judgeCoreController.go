package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/model"
	"main/mysql"
	"main/service"
	"net/http"
	"time"
)

type judgerRequestModel struct {
	Status       []model.Submit `json:"status"`
	JudgingCount int64          `json:"judging_count"`
}

type judgerResponseModel struct {
	Problems []model.JudgeSubmitModel `json:"problems"`
}

type JudgerController struct {
}

func (c *JudgerController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("POST")
	if r.Method != "POST" {
		return
	}
	response := judgerResponseModel{}
	defer func() {
		bf := bytes.NewBuffer([]byte{})
		js := json.NewEncoder(bf)
		js.SetEscapeHTML(false)
		if err := js.Encode(response); err != nil {
			return
		}
		_, _ = w.Write(bf.Bytes())
	}()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	//fmt.Println(string(body))
	request := judgerRequestModel{}
	if err := json.Unmarshal(body, &request); err != nil {
		return
	}
	for i := 0; i < len(request.Status); i++ {
		status := request.Status[i]
		fmt.Printf("[%s] Receive: [Rid:%d] [Status:%d]\n", time.Now().Format("2006-01-02 15:04:05"), status.Rid, status.Status)
		sql := "UPDATE submit SET STATUS = ?, RUN_TIME = ? ,RUN_MEMORY = ?, COMPILATION_MESSAGE = ? WHERE RID = ?"
		_, _ = mysql.DBConn.Exec(sql, status.Status, status.TimeCost, status.MemoryCost, status.CompilationMessage, status.Rid)
		service.UpdateRank(status.Rid)
	}
	if request.JudgingCount < 6 {
		//for ; request.JudgingCount < 6; request.JudgingCount++ {
		//	var submit *model.JudgeSubmitModel
		if submit := service.FetchSubmit(); submit != nil {
			fmt.Printf("[%s] Send: [Rid:%d]\n", time.Now().Format("2006-01-02 15:04:05"), submit.Rid)
			response.Problems = append(response.Problems, *submit)
		}
	}
}

func init() {
	RegisterController("/api/core/j2s/", new(JudgerController))
}

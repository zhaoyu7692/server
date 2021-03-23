package controller

import (
	"encoding/json"
	"io/ioutil"
	"main/model"
	"main/mysql"
	"net/http"
	"time"
)

type submitResponseModel struct {

}

type SubmitController struct {
}

func (c *SubmitController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	var submit model.Submit
	if err := json.Unmarshal(body, &submit); err != nil {
		return
	}
	if _, err := mysql.DBConn.Exec("INSERT INTO submit (PID, UID, CODE, STATUS, LANGUAGE, GMT_CREATED, GMT_MODIFIED) values (?,?,?,0,?,?,?)", submit.Pid, submit.Uid, submit.Code, submit.Language, time.Now(), time.Now()); err != nil {
		return
	}
	_, _ = w.Write([]byte("{code:1}"))
}

func init() {
	RegisterController("/submit/", new(SubmitController))
}

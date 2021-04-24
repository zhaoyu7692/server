package controller

import (
	"encoding/json"
	"io/ioutil"
	"main/model"
	"main/mysql"
	"net/http"
)

type checkUsernameResponseModel struct {
	model.ResponseBaseModel
}

type CheckUsernameController struct {
}

func (c *CheckUsernameController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	var response checkUsernameResponseModel
	defer func() {
		stream, err := json.Marshal(response)
		if err == nil {
			_, _ = w.Write(stream)
		}
	}()
	response.Code = model.PublicFail
	var user model.User
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.Message = "内部错误"
		return
	}
	if err := json.Unmarshal(body, &user); err != nil {
		response.Message = "内部错误"
		return
	}
	var count int64
	if err := mysql.DBConn.Get(&count, "SELECT COUNT(*) FROM user WHERE USERNAME = ?", user.Username); err != nil {
		response.Message = "内部错误"
		return
	}
	if count == 0 {
		response.Code = model.Success
		response.Message = "用户名可用"
	} else {
		response.Message = "用户名已存在"
	}
}

func init() {
	RegisterController("/checkUsername/", new(CheckUsernameController))
}

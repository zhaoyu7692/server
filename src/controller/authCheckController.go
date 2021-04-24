package controller

import (
	"encoding/json"
	"io/ioutil"
	"main/model"
	"main/service"
	"net/http"
)

type authCheckResponseModel struct {
	model.ResponseBaseModel
}

type AuthCheckController struct {
}

func (c *AuthCheckController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	var response authCheckResponseModel
	defer func() {
		if stream, err := json.Marshal(response); err == nil {
			_, _ = w.Write(stream)
		}
	}()
	response.Code = model.JumpLogin
	user := &model.User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	if err := json.Unmarshal(body, user); err != nil {
		return
	}
	code := service.AuthLogin(user.Uid, user.Token)
	if code == service.Authority {
		response.Code = model.Success
	}
}

func init() {
	RegisterController("/authCheck/", new(AuthCheckController))
}

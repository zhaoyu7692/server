package controller

import (
	"encoding/json"
	"io/ioutil"
	"main/model"
	"main/service"
	"net/http"
)

func init() {

	RegisterHandler("/authCheck/", authCheck)

	RegisterHandler("/login/", login)

	RegisterHandler("/register/", register)
}

// Auth Check
type authCheckResponseModel struct {
	model.ResponseBaseModel
}

func authCheck(w http.ResponseWriter, r *http.Request) {
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
	if code == service.Authority || code == service.AuthorityAdmin {
		response.Code = model.Success
	}
}

type loginResponseModel struct {
	User *model.User `json:"user,omitempty"`
	model.ResponseBaseModel
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	var response loginResponseModel
	// 返回
	defer func() {
		if stream, err := json.Marshal(response); err == nil {
			_, _ = w.Write(stream)
		}
	}()
	response.Code = model.PublicFail

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	var user *model.User
	if err := json.Unmarshal(body, &user); err != nil {
		return
	}
	user, status := service.Login(user.Username, user.Password)
	if status == service.LoginSuccess {
		response.Code = model.Success
		response.User = user
	} else {
		switch status {
		case service.LoginFail:
			{
				response.Message = "登录失败，请稍后重试！"
			}
		case service.LoginWrongPassword:
			{
				response.Message = "用户名或密码错误，请确认后重试！"
			}
		}
	}
}

type RegisterAccountResponseModel struct {
	model.ResponseBaseModel
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}
	var response RegisterAccountResponseModel
	defer func() {
		if stream, err := json.Marshal(response); err == nil {
			_, _ = w.Write(stream)
		}
	}()
	response.Code = model.PublicFail
	user := &model.User{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.Message = "注册失败，请稍后重试！"
		return
	}
	if err := json.Unmarshal(body, &user); err != nil {
		response.Message = "注册失败，请稍后重试！"
		return
	}

	user, code := service.Register(user.Username, user.Password)
	switch code {
	case service.RegisterSuccess:
		{
			response.Code = model.Success
			response.Message = "注册成功！"
		}
	case service.RegisterFail:
		{
			response.Message = "注册失败，请稍后重试！"
		}
	case service.RegisterRepetitiveUsername:
		{
			response.Message = "用户名已存在，请重试！"
		}
	}
}

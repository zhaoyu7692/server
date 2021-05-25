package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/model"
	"main/service"
	"net/http"
)

type loginResponseModel struct {
	User *model.User `json:"user,omitempty"`
	model.ResponseBaseModel
}

type LoginController struct {
}

func (c *LoginController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func init() {
	RegisterController("/login/", new(LoginController))
	//routers.InjectRouter("/login", login)
	fmt.Println("loginController init")
}

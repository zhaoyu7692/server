package controller

import (
	"encoding/json"
	"io/ioutil"
	"main/model"
	"main/service"
	"net/http"
)

type RegisterAccountResponseModel struct {
	model.ResponseBaseModel
}

type RegisterAccountController struct {
}

func (c *RegisterAccountController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func init() {
	RegisterController("/register/", new(RegisterAccountController))
}

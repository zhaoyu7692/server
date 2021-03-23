package controller

import (
	"encoding/json"
	"fmt"
	"main/model"
	"net/http"
)

type LoginController struct {
	//Controller
}

func (c *LoginController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := model.User{}
	query := r.URL.Query()
	fmt.Print(query)
	msg, _ := json.Marshal(user)
	_, _ = w.Write(msg)
}

func init() {
	RegisterController("/login/", new(LoginController))
	//routers.InjectRouter("/login", login)
	fmt.Println("loginController init")
}

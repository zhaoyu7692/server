package controller

import (
	"bytes"
	"encoding/json"
	"main/model"
	"main/service"
	"math/rand"
	"net/http"
)

type fetchSubmitController struct {

}

func (c *fetchSubmitController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		return
	}
	submit := service.FetchSubmit()
	bf := bytes.NewBuffer([]byte{})
	js := json.NewEncoder(bf)
	js.SetEscapeHTML(false)
	_ = js.Encode(submit)
	_, _ = w.Write(bf.Bytes())
}
type XXXController struct {

}

func (c *XXXController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		return
	}
	submit := &model.Submit{Pid: rand.Int63()}
	service.StashSubmit(submit)
}
func init() {
	RegisterController("/api/core/fetchSubmit/", new(fetchSubmitController))
	RegisterController("/api/core/stashSubmit/", new(XXXController))
}





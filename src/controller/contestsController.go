package controller

import (
	"net/http"
)

type ContestsController struct {

}

func (c *ContestsController) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func init() {
	RegisterController("/contests/", new(ContestsController))
}

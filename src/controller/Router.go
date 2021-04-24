package controller

import (
	"fmt"
	"net/http"
)

var controllers map[string]http.Handler

func init() {
	controllers = map[string]http.Handler{}
	fmt.Println("router.go" + " init")
}

func RegisterController(pattern string, controller http.Handler) {
	controllers[pattern] = controller
	fmt.Println(pattern)
}

func AddRouter(mux *http.ServeMux) *http.ServeMux {
	for key := range controllers {
		mux.Handle(key, controllers[key])
	}
	return mux
}

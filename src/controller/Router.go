package controller

import (
	"fmt"
	"net/http"
)

var controllers map[string]http.Handler

var handlers map[string]func(w http.ResponseWriter, r *http.Request)

func init() {
	controllers = map[string]http.Handler{}
	handlers = map[string]func(w http.ResponseWriter, r *http.Request){}
	fmt.Println("router.go" + " init")
}

func RegisterController(pattern string, controller http.Handler) {
	controllers[pattern] = controller
	fmt.Println(pattern)
}

func RegisterHandler(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	handlers[pattern] = handler
	fmt.Println(pattern)
}

func AddRouter(mux *http.ServeMux) *http.ServeMux {
	for key := range controllers {
		mux.Handle(key, controllers[key])
	}
	for key := range handlers {
		mux.HandleFunc(key, handlers[key])
	}
	return mux
}

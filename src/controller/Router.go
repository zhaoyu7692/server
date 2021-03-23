package controller

import "net/http"

var controllers map[string]http.Handler

func init() {
	controllers = map[string]http.Handler{}
}

func RegisterController(pattern string, controller http.Handler) {
	controllers[pattern] = controller
}

func AddRouter(mux *http.ServeMux) *http.ServeMux {
	for key := range controllers {
		mux.Handle(key, controllers[key])
	}
	return mux
}
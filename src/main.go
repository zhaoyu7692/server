package main

import (
	"fmt"
	"log"
	"main/controller"
	"main/utils"
	"net/http"
	"os"
)

func A() {
	fmt.Println('A', os.Getpid())
	x := 1000
	for i := 0; i < 100000; i++ {
		x += 1
	}
	fmt.Println('A', os.Getpid())
}

func B() {
	fmt.Println('B', )
	x := 1000
	for i := 0; i < 100000000; i++ {
		x += 1
	}
	fmt.Println('B', os.Getpid())
}

func main() {

	//pbc := controller.ProblemController{}
	//pbc.

	//var mux = http.DefaultServeMux
	//mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
	//	_, _ = w.Write([]byte("ping ping"))
	//	//w.Write(byte[]("ping ping\n"))
	//})
	//server := http.Server{
	//	Addr: 	"0.0.0.0:8887",
	//	Handler: mux,
	//}
	//log.Fatal(server.ListenAndServe())

	//for i := 0; i < 1000000000; i++ {
	//	go A()
	//	go B()
	//}
	//fmt.Println(runtime.NumCPU())
	openJudge := utils.GlobalConfig.OpenJudge
	address := fmt.Sprintf("%v:%v", openJudge.Host, openJudge.Port)
	server := http.Server{
		Addr: address,
		Handler: controller.AddRouter(http.NewServeMux()),
	}
	log.Fatal(server.ListenAndServe())
}

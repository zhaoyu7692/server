package main

import (
	"fmt"
	"log"
	"main/controller"
	"main/service"
	"main/utils"
	"net/http"
)

func main() {
	service.UpdateRank(1)
	openJudge := utils.GlobalConfig.OpenJudge
	address := fmt.Sprintf("%v:%v", openJudge.Host, openJudge.Port)
	server := http.Server{
		Addr:    address,
		Handler: controller.AddRouter(http.NewServeMux()),
	}
	log.Fatal(server.ListenAndServe())
}

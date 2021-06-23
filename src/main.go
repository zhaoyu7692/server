package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"main/controller"
	"main/mysql"
	"main/utils"
	"net/http"
	"strconv"
)

func updateDataBase() {
	path := utils.GlobalConfig.Path
	problems, err := ioutil.ReadDir(path.Data)
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(problems); i++ {
		if problems[i].IsDir() {
			if pid, err := strconv.ParseInt(problems[i].Name(), 10, 64); err == nil {
				if cases, err := ioutil.ReadDir(fmt.Sprintf("%s%d", path.Data, pid)); err == nil {
					//fmt.Println(cases)
					for _, file := range cases {
						data, _ := ioutil.ReadFile(fmt.Sprintf("%s%d/%s", path.Data, pid, file.Name()))
						hash := sha256.New()
						hash.Write(data)
						hashBytes := hash.Sum(nil)
						_, _ = mysql.DBConn.Exec("INSERT INTO testcase_mapping (PID, FILENAME, `KEY`, PATH) VALUES (?,?,?,?)", pid, file.Name(), hex.EncodeToString(hashBytes), fmt.Sprintf("%s%d/%s", utils.GlobalConfig.Path.Data, i, file.Name()))
						//fmt.Println(fmt.Sprintf("%s%d/%s", path.Data, pid, file.Name()))
						//fmt.Println(err)
					}
				}
			}
		}
	}
}

func main() {
	openJudge := utils.GlobalConfig.OpenJudge
	updateDataBase()
	address := fmt.Sprintf("%v:%v", openJudge.Host, openJudge.Port)
	server := http.Server{
		Addr:    address,
		Handler: controller.AddRouter(http.NewServeMux()),
	}
	log.Fatal(server.ListenAndServe())
}

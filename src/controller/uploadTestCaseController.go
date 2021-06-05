package controller

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"main/model"
	"main/mysql"
	"main/utils"
	"net/http"
	"os"
)

type UploadTestCaseController struct {
}

func (c *UploadTestCaseController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Test")
	if r.Method != "POST" {
		w.WriteHeader(500)
		return
	}
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer file.Close()
	dataBuffer := new(bytes.Buffer)
	if _, err = dataBuffer.ReadFrom(file); err != nil {
		w.WriteHeader(500)
		return
	}
	hash := sha256.New()
	hash.Write(dataBuffer.Bytes())
	hashBytes := hash.Sum(nil)
	hashCode := hex.EncodeToString(hashBytes)

	// 复制文件到目标目录
	filePath := utils.GlobalConfig.Path.Data + hashCode
	//+ fileHeader.Filename
	fmt.Println(filePath)
	err = ioutil.WriteFile(filePath, dataBuffer.Bytes(), os.ModePerm)
	if err != nil {
		log.Printf("UploadHandler: 文件复制失败！ -> {%s}", err)
		w.WriteHeader(500)
		return
	}
	item := model.ResourceMappingModel{}
	if err := mysql.DBConn.Get(&item, "SELECT * FROM resource_mapping WHERE SHA_KEY = ?", hashCode); err == nil {
		_, err = mysql.DBConn.Exec("UPDATE resource_mapping SET PATH = ?, FILENAME  = ? WHERE SHA_KEY = ?", filePath, fileHeader.Filename, hashCode)
		if err != nil {
			w.WriteHeader(500)
			return
		}
	} else {
		_, err = mysql.DBConn.Exec("INSERT INTO resource_mapping (SHA_KEY, PATH, FILENAME) VALUES (?,?,?)", hashCode, filePath, fileHeader.Filename)
		if err != nil {
			w.WriteHeader(500)
			return
		}
	}
	w.WriteHeader(200)
}

func init() {
	//for i := 1; i < 17; i++ {
	//	files, _ := ioutil.ReadDir(fmt.Sprintf("%s%d", utils.GlobalConfig.Path.Data, i))
	//	for _, file := range files {
	//		data, _ := ioutil.ReadFile(fmt.Sprintf("%s%d/%s", utils.GlobalConfig.Path.Data, i, file.Name()))
	//		hash := sha256.New()
	//		hash.Write(data)
	//		hashBytes := hash.Sum(nil)
	//		_, err := mysql.DBConn.Exec("INSERT INTO testcase_mapping (PID, FILENAME, `KEY`, PATH) VALUES (?,?,?,?)", i, file.Name(), hex.EncodeToString(hashBytes), fmt.Sprintf("%s%d/%s", utils.GlobalConfig.Path.Data, i, file.Name()))
	//		fmt.Println(err)
	//	}
	//}
	RegisterController("/uploadTestCase/", new(UploadTestCaseController))
}

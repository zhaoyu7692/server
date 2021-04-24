package controller

import (
	"fmt"
	"io"
	"main/utils"
	"net/http"
	"os"
)

type UploadTestCaseController struct {
}

func (c *UploadTestCaseController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Test")
	if r.Method != "POST" {
		return
	}
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return
	}
	defer file.Close()

	//headerByte, _ := json.Marshal(fileHeader.Header)
	//log.Printf("当前文件：Filename - >{%s}, Size -> {%v}, FileHeader -> {%s}", fileHeader.Filename, fileHeader.Size, string(headerByte))

	newFile, err := os.Create(utils.GlobalConfig.OpenJudge.DataPath + "/" + fileHeader.Filename)
	if err != nil {
		return
	}
	defer newFile.Close()

	// 复制文件到目标目录
	_, _ = io.Copy(newFile, file)
	//if errCopy != nil {
	//log.Printf("UploadHandler: 文件复制失败！ -> {%s}", err)
	//_, _ = io.WriteString(w, "服务器错误！")
	//return
	//}

}

func init() {
	RegisterController("/uploadTestCase/", new(UploadTestCaseController))
}

package controller

import (
	"encoding/json"
	"io/ioutil"
	"main/model"
	"main/mysql"
	"net/http"
)

type submitStatusModel struct {
	ItemList []model.SubmitStatus `json:"item_list"`
	model.ResponsePaginationModel
}

type SubmitStatusController struct {

}

func (c *SubmitStatusController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	var page model.RequestPaginationModel
	if err = json.Unmarshal(body, &page); err != nil {
		return
	}


	var submitStatus []model.SubmitStatus
	if err := mysql.DBConn.Select(&submitStatus, "SELECT RID, UID, LANGUAGE, STATUS, RUN_TIME, RUN_MEMORY, GMT_CREATED FROM SUBMIT LIMIT ?, ?", page.Offset, page.Size); err != nil {
		return
	}
	for index, _ := range submitStatus {
		if err := mysql.DBConn.Get(&submitStatus[index], "SELECT USERNAME FROM USER WHERE ID = ?", submitStatus[index].Uid); err != nil {
			continue
		}
	}
	var responseModel submitStatusModel
	if err := mysql.DBConn.Get(&responseModel.Total, "SELECT COUNT(RID) FROM SUBMIT"); err != nil {
		return
	}
	responseModel.ItemList = submitStatus
	responseModel.Size = page.Size

	stream, err := json.Marshal(responseModel)
	if err != nil {
		return
	}
	_, _ = w.Write(stream)
}

func init() {
	RegisterController("/submitStatus/", new(SubmitStatusController))
}

package controller

import (
	"encoding/json"
	"main/model"
	"main/mysql"
	"main/utils"
	"math"
	"net/http"
)

type submitStatusResponseModel struct {
	SubmitStatus struct {
		ItemList []model.Submit `json:"item_list"`
		model.ResponsePaginationModel
	} `json:"submit_status"`
	model.ResponseBaseModel
}

type SubmitStatusController struct {
}

func (c *SubmitStatusController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	var responseModel submitStatusResponseModel
	defer func() {
		if stream, err := json.Marshal(responseModel); err == nil {
			_, _ = w.Write(stream)
		}
	}()
	responseModel.Code = model.PublicFail
	query := r.URL.Query()
	cid := utils.StringConstraint(query.Get("cid"), 0, math.MaxInt64, math.MaxInt64)
	page := utils.StringConstraint(query.Get("page"), 1, math.MaxInt64, 1)
	size := utils.StringConstraint(query.Get("size"), 20, 50, 20)
	offset := (page - 1) * size
	sql := "SELECT RID, CID,`INDEX`, s.UID, USERNAME, LANGUAGE, STATUS, COMPILATION_MESSAGE, RUN_TIME, RUN_MEMORY, SUBMIT_TIME FROM submit as s, user as u WHERE CID = ? AND s.UID = u.UID ORDER BY RID DESC LIMIT ?, ?"
	if err := mysql.DBConn.Select(&responseModel.SubmitStatus.ItemList, sql, cid, offset, size); err != nil {
		return
	}
	if err := mysql.DBConn.Get(&responseModel.SubmitStatus.Total, "SELECT COUNT(RID) FROM submit WHERE CID = ?", cid); err != nil {
		return
	}
	responseModel.SubmitStatus.Size = size
	responseModel.Code = model.Success
}

func init() {
	RegisterController("/submitStatus/", new(SubmitStatusController))
}

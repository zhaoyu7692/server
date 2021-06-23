package controller

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	"main/model"
	"main/mysql"
	"main/redispool"
	"net/http"
	"time"
)

func init() {
	RegisterHandler("/systemStatistics/", systemStatistics)
}

type StatisticsItemModel struct {
	Name string `json:"name"`

	Data []int64 `json:"data"`
}

type StatisticsModel struct {
	Legend []string              `json:"legend"`
	Items  []StatisticsItemModel `json:"items"`
}

type SystemStatisticsResponseModel struct {
	ActivityData StatisticsModel `json:"activity_data"`
	SubmitData   StatisticsModel `json:"submit_data"`
	model.ResponseBaseModel
}

func systemStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}
	conn := redispool.Get()
	defer func() {
		conn.Close()
	}()
	if reply, err := redis.Bytes(conn.Do("GET", "openjudge_statistics")); err == nil {
		_, _ = w.Write(reply)
		return
	}
	response := SystemStatisticsResponseModel{}
	defer func() {
		if response.Code != model.Success {
			return
		}
		if stream, err := json.Marshal(response); err == nil {
			_, _ = conn.Do("SET", "openjudge_statistics", stream)
			_, _ = conn.Do("EXPIRE", "openjudge_statistics", 30)
			_, _ = w.Write(stream)
		}
	}()
	response.Code = model.PublicFail
	now := time.Now()
	day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	submitStatusCategory := make(map[string][]model.JudgeStatus)
	submitStatusCategory["答案正确"] = []model.JudgeStatus{model.JudgeStatusAccept}
	//submitStatusCategory["编译错误"] = []model.JudgeStatus{model.JudgeStatusCompilationError}
	//submitStatusCategory["编译超时"] = []model.JudgeStatus{model.JudgeStatusCompilationTimeLimitExceeded}
	//submitStatusCategory["格式错误"] = []model.JudgeStatus{model.JudgeStatusPresentationError}
	submitStatusCategory["运行错误"] = []model.JudgeStatus{model.JudgeStatusRuntimeError}
	submitStatusCategory["答案错误"] = []model.JudgeStatus{model.JudgeStatusWrongAnswer}
	submitStatusCategory["时间超限"] = []model.JudgeStatus{model.JudgeStatusTimeLimitExceeded}
	submitStatusCategory["内存超限"] = []model.JudgeStatus{model.JudgeStatusMemoryLimitExceeded}
	//submitStatusCategory["输出超限"] = []model.JudgeStatus{model.JudgeStatusOutputLimitExceeded}
	//submitStatusCategory["全部提交"] = []model.JudgeStatus{-1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
	//fmt.Println(submitStatusCategory)

	for key, value := range submitStatusCategory {
		tempDay := day.Add(-6 * 24 * time.Hour)
		item := StatisticsItemModel{Name: key}
		count := int64(0)
		for i := 0; i < 7; i++ {
			if query, args, err := sqlx.In("SELECT COUNT(1) FROM submit WHERE SUBMIT_TIME >= ? AND SUBMIT_TIME < ? AND STATUS IN (?)", tempDay, tempDay.Add(24*time.Hour), value); err == nil {
				if err := mysql.DBConn.Get(&count, query, args...); err != nil {
					return
				}
				item.Data = append(item.Data, count)
			} else {
				return
			}
			tempDay = tempDay.Add(24 * time.Hour)
		}
		response.SubmitData.Items = append(response.SubmitData.Items, item)
	}
	newUserItemModel := StatisticsItemModel{Name: "新增用户"}
	totalUserItemModel := StatisticsItemModel{Name: "累积用户"}
	hotUserItemModel := StatisticsItemModel{Name: "活跃用户"}
	var userLegend []string
	tempDay := day.Add(-6 * 24 * time.Hour)
	for i := 0; i < 7; i++ {
		userLegend = append(userLegend, tempDay.Format("01-02"))
		count := int64(0)
		if err := mysql.DBConn.Get(&count, "SELECT COUNT(1) FROM user WHERE GMT_CREATED < ?", tempDay.Add(24*time.Hour)); err != nil {
			return
		}
		totalUserItemModel.Data = append(totalUserItemModel.Data, count)
		if err := mysql.DBConn.Get(&count, "SELECT COUNT(1) FROM user WHERE GMT_CREATED >= ? AND GMT_CREATED < ?", tempDay, tempDay.Add(24*time.Hour)); err != nil {
			return
		}
		newUserItemModel.Data = append(newUserItemModel.Data, count)
		if err := mysql.DBConn.Get(&count, "SELECT COUNT(DISTINCT UID) FROM submit WHERE SUBMIT_TIME >= ? AND SUBMIT_TIME < ?", tempDay, tempDay.Add(24*time.Hour)); err != nil {
			return
		}
		hotUserItemModel.Data = append(hotUserItemModel.Data, count)
		tempDay = tempDay.Add(24 * time.Hour)
	}
	response.ActivityData = StatisticsModel{Items: []StatisticsItemModel{newUserItemModel, totalUserItemModel, hotUserItemModel}, Legend: userLegend}
	response.SubmitData.Legend = userLegend
	response.Code = model.Success

	//{
	//	tempDay := day.Add(-6 * 24 * time.Hour)
	//	item := StatisticsItemModel{Name:}
	//	count := int64(0)
	//	for i := 0; i < 7; i++ {
	//		if err := mysql.DBConn.Get(&count, "SELECT COUNT(1) FROM submit WHERE SUBMIT_TIME >= ?", tempDay); err != nil {
	//			return
	//		}
	//		item.Data = append(item.Data, count)
	//		if i > 0 {
	//			item.Data[i-1] -= count
	//		}
	//		tempDay = tempDay.Add(24 * time.Hour)
	//	}
	//	response.SubmitData.Items = append(response.SubmitData.Items, item)
	//}

}

package service

import (
	"fmt"
	"main/model"
	"main/mysql"
	"sync"
)

var submitPool []*model.JudgeSubmitModel
var submitPoolMutex sync.Mutex

func StashSubmit(submit *model.JudgeSubmitModel) {
	submitPoolMutex.Lock()
	defer submitPoolMutex.Unlock()
	submitPool = append(submitPool, submit)
	fmt.Println(submitPool)
}

func FetchSubmit() *model.JudgeSubmitModel {
	//fmt.Println(submitPool)
	submitPoolMutex.Lock()
	defer submitPoolMutex.Unlock()
	if len(submitPool) <= 0 {
		return nil
	}
	res := submitPool[0]
	submitPool = submitPool[1:]
	return res
}

func init() {
	var submits []model.JudgeSubmitModel
	if err := mysql.DBConn.Select(&submits, "SELECT RID, CID, `INDEX`, CODE, STATUS, LANGUAGE FROM submit WHERE STATUS = ?", 0); err != nil {
		return
	}
	for i := 0; i < len(submits); i++ {
		if err := mysql.DBConn.Get(&submits[i], "SELECT p.PID, p.TIME_LIMIT, p.MEMORY_LIMIT FROM problem as p, contest_problem_mapping as cp WHERE CID = ? AND `INDEX` = ? AND p.PID = cp.PID", submits[i].Cid, submits[i].Index); err == nil {
			submitPool = append(submitPool, &submits[i])
		}
	}
}

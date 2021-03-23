package service

import (
	"fmt"
	"main/model"
	"main/mysql"
	"sync"
)

var submitPool []*model.Submit
var mutex sync.Mutex

func StashSubmit(submit *model.Submit) {
	mutex.Lock()
	defer mutex.Unlock()
	submitPool = append(submitPool, submit)
	fmt.Println(submitPool)
}

func FetchSubmit() *model.Submit {
	if len(submitPool) <= 0 {
		return nil
	}
	fmt.Println(submitPool)
	mutex.Lock()
	defer mutex.Unlock()
	res := submitPool[0]
	submitPool = submitPool[1:]
	return res
}

func init() {
	var submits []model.Submit
	if err := mysql.DBConn.Select(&submits, "SELECT * FROM submit WHERE STATUS = 0"); err != nil {
		return
	}
	for i := 0; i < len(submits); i++ {
		submitPool = append(submitPool, &submits[i])
	}
}
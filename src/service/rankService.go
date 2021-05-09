package service

import (
	"fmt"
	"main/model"
	"main/mysql"
	"main/redispool"
)

func UpdateRank(rid int64) {
	go func() {
		submit := model.Submit{}
		if err := mysql.DBConn.Get(&submit, "SELECT CID, `INDEX`, UID FROM submit WHERE RID = ?", rid); err != nil {
			return
		}
		// 废弃相关 redis 缓存
		_, _ = redispool.Get().Do("DEL", fmt.Sprintf("contest_rank_key_cid_%d", submit.Cid))

		//var acceptSubmits []model.Submit
		//if err := mysql.DBConn.Select(&acceptSubmits, "SELECT RID, SUBMIT_TIME FROM submit WHERE CID = ? AND `INDEX` = ? AND UID = ? ORDER BY RID LIMIT 0, 1", submit.Cid, submit.Index, submit.Uid); err != nil {
		//	return
		//}
		//if len(acceptSubmits) == 0 {
		//	_, _ = mysql.DBConn.Exec("UPDATE contest_rank SET TRY_COUNT = (SELECT COUNT(*) FROM submit WHERE CID = ? AND `INDEX` = ? AND UID = ?) WHERE CID = ? AND `INDEX` = ? AND UID = ?", submit.Cid, submit.Index, submit.Uid, submit.Cid, submit.Index, submit.Uid)
		//} else {
		//	acceptSubmit := acceptSubmits[0]
		//	_, _ = mysql.DBConn.Exec("UPDATE contest_rank SET TRY_COUNT = (SELECT COUNT(*) FROM submit WHERE CID = ? AND `INDEX` = ? AND UID = ? AND RID < ?), ACCEPT_TIME = ? WHERE CID = ? AND `INDEX` = ? AND UID = ?", submit.Cid, submit.Index, submit.Uid, acceptSubmit.Rid, acceptSubmit.SubmitTime, submit.Cid, submit.Index, submit.Uid)
		//}
	}()
}

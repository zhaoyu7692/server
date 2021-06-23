package service

import "main/dao"

type ContestDetailModel struct {
	Contest  dao.ContestTableModel
	Problems []dao.ProblemDetailModel
}

func GetContestInfo(cid int64) *ContestDetailModel {
	if contest := dao.GetContestWithCid(cid); contest != nil {
		if problems := dao.GetProblemsWithCid(cid); problems != nil {
			return &ContestDetailModel{
				Contest:  *contest,
				Problems: *problems,
			}
		}
	}
	return nil
}


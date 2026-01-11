package api

import (
	"github.com/gin-gonic/gin"
	"go-novel/app/service/api/rank_service"
	"go-novel/utils"
)

type Rank struct{}

func (rank *Rank) RankList(c *gin.Context) {
	list, err := rank_service.GetRankList()
	if err != nil {
		utils.FailEncrypt(c, err, "")
		return
	}

	utils.SuccessEncrypt(c, list, "ok")
}

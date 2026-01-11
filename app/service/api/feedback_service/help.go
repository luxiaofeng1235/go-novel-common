package feedback_service

import (
	"go-novel/app/models"
	"go-novel/global"
)

func GetFeedbackHelpById(id int64) (help *models.McFeedbackHelp, err error) {
	err = global.DB.Model(models.McFeedbackHelp{}).Where("id", id).First(&help).Error
	return
}

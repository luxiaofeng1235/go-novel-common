package withdraw_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/common/setting_service"
	"go-novel/global"
	"go-novel/utils"
	"gorm.io/gorm"
)

func GetWithdrawLimit() (limits []*models.WithdrawLimitRes, err error) {
	err = global.DB.Model(models.WithdrawLimitRes{}).Order("sort desc").Where("status = 1").Find(&limits).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}

	for _, limit := range limits {
		var rate, cion int64
		rate, cion, err = setting_service.GetCionByMoney(limit.Money)
		if err != nil {
			continue
		}
		limit.Rate = rate
		limit.Cion = cion
	}
	return
}

func GetAccountDetailById(req *models.WithdrawAccountDetailReq) (alipay *models.AccountDetailRes, wxpay *models.AccountDetailRes, err error) {
	userId := req.UserId
	var accountAli *models.McWithdrawAccount
	var accountWx *models.McWithdrawAccount
	accountAli, _ = GetAccountByPayType(userId, 1)

	alipay = new(models.AccountDetailRes)
	alipay.McWithdrawAccount = accountAli
	if accountAli.CardPic != "" {
		alipay.CardPicUrl = utils.GetFileUrl(alipay.CardPic)
	}

	wxpay = new(models.AccountDetailRes)
	accountWx, _ = GetAccountByPayType(userId, 2)
	wxpay.McWithdrawAccount = accountWx
	if accountWx.CardPic != "" {
		wxpay.CardPicUrl = utils.GetFileUrl(wxpay.CardPic)
	}
	return
}

func AccountSave(req *models.WithdrawAccountSaveReq) (err error) {
	payType := req.PayType
	cardName := req.CardName
	cardNumber := req.CardNumber
	cardPic := req.CardPic
	userId := req.UserId
	if payType <= 0 {
		err = fmt.Errorf("%v", "账号类型不能为空")
		return
	}
	if cardName == "" {
		err = fmt.Errorf("%v", "姓名不能为空")
		return
	}
	if cardNumber == "" {
		err = fmt.Errorf("%v", "账号不能为空")
		return
	}
	if cardPic == "" {
		err = fmt.Errorf("%v", "收款码不能为空")
		return
	}
	if userId <= 0 {
		err = fmt.Errorf("%v", "未登录")
		return
	}

	var count int64
	count = GetCountByPayType(payType, userId)
	if count <= 0 {
		account := models.McWithdrawAccount{
			Uid:        userId,
			PayType:    payType,
			CardName:   cardName,
			CardNumber: cardNumber,
			CardPic:    cardPic,
			Addtime:    utils.GetUnix(),
			Uptime:     utils.GetUnix(),
		}
		if err = global.DB.Create(&account).Error; err != nil {
			err = fmt.Errorf("记录失败，稍后再试 err=%v", err.Error())
			return
		}
	} else {
		data := make(map[string]interface{})
		data["card_name"] = cardName
		data["card_number"] = cardNumber
		data["card_pic"] = cardPic
		data["uptime"] = utils.GetUnix()
		err = UpdateAccountByPayType(payType, userId, data)
		if err != nil {
			return
		}
	}
	return
}

func AccountDel(req *models.WithdrawAccountDelReq) (err error) {
	accountId := req.AccountId
	userId := req.UserId
	if accountId <= 0 {
		err = fmt.Errorf("%v", "账户ID不能为空")
		return
	}
	if userId <= 0 {
		err = fmt.Errorf("%v", "未登录")
		return
	}
	err = DeleteAccountByAccountId(accountId, userId)
	return
}

func Apply(req *models.WithdrawApplyReq) (err error) {
	limitId := req.LimitId
	if limitId <= 0 {
		err = fmt.Errorf("%v", "请选择提现金额")
		return
	}
	accountId := req.AccountId
	if accountId <= 0 {
		err = fmt.Errorf("%v", "请选择提现账号")
		return
	}
	account, err := GetWithdrawAccountById(accountId)
	if err != nil {
		err = fmt.Errorf("%v", "提现账号不存在")
		return
	}
	if account.Id <= 0 {
		err = fmt.Errorf("%v", "提现账号不存在")
		return
	}

	userId := req.UserId
	if userId <= 0 {
		err = fmt.Errorf("%v", "账号未登录")
		return
	}

	limit, err := GetWithdrawLimitById(limitId)
	if err != nil {
		return
	}
	if limit.Id <= 0 {
		err = fmt.Errorf("%v", "提现金额异常")
		return
	}
	if limit.Status != 1 {
		err = fmt.Errorf("%v", "提现维护中")
		return
	}
	withdrawMoney := limit.Money

	_, costCion, err := setting_service.GetCionByMoney(withdrawMoney)

	user, err := getUserById(userId)
	if err != nil {
		err = fmt.Errorf("获取用户信息失败 %v", err.Error())
		return
	}
	if user.Status == 0 {
		err = fmt.Errorf("%v", "该用户已被禁用")
		return
	}
	if costCion > user.Cion {
		err = fmt.Errorf("%v", "账户金币余额不足")
		return
	}

	readSecond := limit.Read * 60
	if readSecond > 0 {
		todayUnix := utils.GetTodayUnix()
		var second int64
		err = global.DB.Model(models.McBookTime{}).Select("coalesce(sum(second), 0)").Where("uid = ? and addtime >= ?", userId, todayUnix).Scan(&second).Error
		if err != nil {
			return
		}
		if readSecond > second {
			err = fmt.Errorf("%v", "请先完成阅读任务")
			return
		}
	}

	//获取转换比例 校验余额 判断要求是否完成 扣取等值金币 添加金币变动记录 添加提现记录
	tx := global.DB.Begin()
	err = tx.Model(models.McUser{}).Where("id = ?", userId).Update("cion", gorm.Expr("cion - ?", costCion)).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		tx.Rollback()
		return
	}

	change := models.McCionChange{
		Tid:        0,
		Uid:        userId,
		Cion:       costCion,
		ChangeType: 2,
		OperatType: 5,
		Addtime:    utils.GetUnix(),
	}
	err = tx.Model(models.McCionChange{}).Create(&change).Error
	if err != nil {
		tx.Rollback()
		global.Sqllog.Errorf("%v", err.Error())
		return
	}

	apply := models.McWithdrawApply{
		Uid:           userId,
		PayType:       account.PayType,
		AccountId:     accountId,
		WithdrawMoney: withdrawMoney,
		Cion:          costCion,
		CardName:      account.CardName,
		CardNumber:    account.CardNumber,
		CardPic:       account.CardPic,
		Status:        0,
		Addtime:       utils.GetUnix(),
	}
	err = tx.Model(models.McWithdrawApply{}).Create(&apply).Error
	if err != nil {
		tx.Rollback()
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	tx.Commit()
	return
}

func WithdrawApplyList(req *models.WithdrawApplyListReq) (applys []*models.WithdrawApplyListRes, total int64, err error) {
	userId := req.UserId
	if userId <= 0 {
		err = fmt.Errorf("%v", "账号未登录")
		return
	}
	var list []*models.McWithdrawApply
	db := global.DB.Model(&models.McWithdrawApply{}).Order("id desc")

	db = db.Where("uid = ?", userId)

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	pageNum := req.Page
	pageSize := req.Size

	if pageSize == 0 || pageSize > 100 {
		pageSize = 15
	}
	if pageSize == 0 {
		pageSize = 1
	}

	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}

	if len(list) <= 0 {
		return
	}

	for _, val := range list {
		var typeName string
		if val.PayType == 1 {
			typeName = "支付宝"
		} else if val.PayType == 2 {
			typeName = "微信"
		}
		var statusName string
		if val.Status == 0 {
			statusName = "待处理"
		} else if val.Status == 1 {
			statusName = "已处理"
		} else {
			statusName = "异常"
		}
		apply := &models.WithdrawApplyListRes{
			Id:            val.Id,
			UserId:        val.Uid,
			PayType:       val.PayType,
			TypeName:      typeName,
			WithdrawMoney: val.WithdrawMoney,
			Cion:          val.Cion,
			Status:        val.Status,
			StatusName:    statusName,
			Addtime:       val.Addtime,
		}
		applys = append(applys, apply)
	}
	return
}

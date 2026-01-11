package vip_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/api/book_service"
	"go-novel/app/service/common/setting_service"
	"go-novel/global"
	"go-novel/utils"
	"gorm.io/gorm"
)

func GetVipBookStore(req *models.VipBookStoreReq) (isVip int64, vipCount int64, choices []*models.McBook, hots []*models.McBook, news []*models.McBook, err error) {
	userId := req.UserId
	if userId > 0 {
		isVip = GetUserVipById(userId)
		if isVip > 0 {
			isVip = 1
		}
	}
	fmt.Println("11111111111111111111", req.Ip)
	vipCount = GetVipBookCountById() //获取vip统计的总数
	choices, err = GetVipBooks(1, 0, 0, 6, req.DeviceType, req.PackageName, req.Ip, req.Mark)
	hots, err = GetVipBooks(0, 1, 0, 3, req.DeviceType, req.PackageName, req.Ip, req.Mark)
	news, err = GetVipBooks(0, 1, 0, 3, req.DeviceType, req.PackageName, req.Ip, req.Mark)
	return
}

// 获取其他书籍统计
func GetVipBooks(isChoice, isHot, isNew, size int, deviceType string, packageName string, ip string, mark string) (list []*models.McBook, err error) {
	db := global.DB.Model(&models.McBook{})
	//.Order("uptime desc")
	fmt.Println("11111111111111111111", ip)
	bookStatus, err := book_service.GetBookCopyright(deviceType, packageName, ip, mark)
	if err != nil {
		global.Sqllog.Errorf("GetBookCopyright err:%v", err)
	}
	if bookStatus == 1 {
		db = db.Where("status = 1 and is_banquan = 1").Debug()
	} else {
		db = db.Where("status = 1").Debug()
	}

	//到时候需要再放开
	//db = db.Where("is_pay = 2")
	if size <= 0 {
		size = 6
	}

	if isChoice > 0 {
		db = db.Where("is_choice = ?", isChoice)
	}

	if isHot > 0 {
		db = db.Where("is_hot = ?", isHot)
	}

	if isNew > 0 {
		db = db.Where("is_new = ?", isNew)
	}

	var total int64
	db.Count(&total) //统计总数
	global.Requestlog.Infof("VipBook推荐的书籍总数 total = %v", total)
	//获取随机的数量
	offsetNum := utils.GetBookRandPosition(total)
	fmt.Println("22222222222222", offsetNum)
	//if pageNum > 0 && size > 0 {
	//	//请求分页的情况数据加载
	//	err = db.Offset((pageNum - 1) * size).Limit(size).Find(&list).Error
	//} else {
	//不需要分页的数据加载
	//err = db.Limit(size).Find(&list).Error
	//根据随机数来进行随机进行跳转选择
	err = db.Offset(offsetNum).Limit(size).Find(&list).Error
	//}
	if len(list) <= 0 {
		return
	}
	//处理图片和相关数据
	for _, val := range list {
		if val.Pic != "" {
			val.Pic = utils.GetFileUrl(val.Pic)
		}
	}
	return
}

func GetVipMessage(req *models.VipMessageReq) (records []*models.McVipMessage, err error) {
	userId := req.UserId
	db := global.DB.Model(models.McVipMessage{}).Order("id desc")
	if req.UserId > 0 {
		db = db.Where("uid = ?", userId)
	}
	err = db.Find(&records).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	for _, record := range records {
		record.Pic = utils.GetFileUrl(record.Pic)
	}
	return
}

func GetVipCardRmb() (records []*models.McVipCard, err error) {
	err = global.DB.Model(models.McVipCard{}).Order("sort asc").Where("status = 1 and is_rmb = 1").Find(&records).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetVipCardCion() (cions []*models.VipCardCionRes, err error) {
	var cards []*models.McVipCard
	err = global.DB.Model(models.McVipCard{}).Order("sort asc").Where("status = 1 and is_cion = 1").Find(&cards).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	if len(cards) <= 0 {
		return
	}
	for _, val := range cards {
		disPrice := val.DisPrice
		if disPrice <= 0 {
			continue
		}
		var rate, gold int64
		rate, gold, err = setting_service.GetCionByMoney(disPrice)
		if err != nil {
			continue
		}
		cion := &models.VipCardCionRes{
			Id:       val.Id,
			CardName: val.CardName,
			Price:    val.Price,
			DisRate:  val.DisRate,
			DisPrice: val.DisPrice,
			CionRate: rate,
			Cion:     gold,
		}
		cions = append(cions, cion)
	}
	return
}

func ExchangeVip(req *models.ExchangeVipReq) (err error) {
	userId := req.UserId
	cardId := req.VipCardId
	if userId <= 0 {
		err = fmt.Errorf("%v", "账号未登录")
		return
	}
	vip, err := GetVipById(cardId)
	if err != nil {
		return
	}
	if vip.Id <= 0 {
		return
	}
	if vip.IsCion != 1 {
		err = fmt.Errorf("%v", "暂时无法兑换")
		return
	}
	if vip.Status != 1 {
		err = fmt.Errorf("%v", "暂时无法兑换")
		return
	}
	var gold int64
	_, gold, err = setting_service.GetCionByMoney(vip.DisPrice)
	if err != nil {
		return
	}

	user, err := getUserById(userId)
	if err != nil {
		err = fmt.Errorf("获取用户信息失败 %v", err.Error())
		return
	}
	if user.Status == 0 {
		err = fmt.Errorf("%v", "该用户已被禁用")
		return
	}
	if gold > user.Cion {
		err = fmt.Errorf("%v", "账户金币余额不足")
		return
	}
	//金币扣款 增加vip 余额改变记录 兑换记录
	tx := global.DB.Begin()
	err = tx.Model(models.McUser{}).Where("id = ?", userId).Update("cion", gorm.Expr("cion - ?", gold)).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		tx.Rollback()
		return
	}

	day := vip.Day
	if day > 0 {
		us := make(map[string]interface{})
		us["vip"] = 1
		timeStamp := utils.GetUnix()
		if user.Viptime > timeStamp {
			us["viptime"] = user.Viptime + 86400*day
		} else {
			us["viptime"] = utils.GetUnix() + 86400*day
		}
		err = tx.Model(models.McUser{}).Where("id = ?", user.Id).Updates(us).Error
		if err != nil {
			tx.Rollback()
			return
		}
	}

	change := models.McCionChange{
		Tid:        0,
		Uid:        userId,
		Cion:       gold,
		ChangeType: 2,
		OperatType: 6,
		Addtime:    utils.GetUnix(),
	}
	err = tx.Model(models.McCionChange{}).Create(&change).Error
	if err != nil {
		tx.Rollback()
		global.Sqllog.Errorf("%v", err.Error())
		return
	}

	card := models.McVipMessage{
		Uid:      userId,
		Nickname: user.Nickname,
		Pic:      user.Pic,
		CardId:   cardId,
		CardName: vip.CardName,
		PayType:  2,
		Addtime:  utils.GetUnix(),
	}
	err = tx.Model(models.McVipMessage{}).Create(&card).Error
	if err != nil {
		tx.Rollback()
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	tx.Commit()
	return
}

package adver_service

import (
	"fmt"
	"go-novel/app/models"
	"go-novel/global"
	"go-novel/utils"
	"gorm.io/gorm"
	"log"
	"strings"
)

func GetAdverMap() (madver map[string]interface{}, err error) {
	//书架广告
	var adverBookshelf *models.McAdver
	err = global.DB.Model(models.McAdver{}).Order("weight desc").Where("status = 1 and adver_type = 1").First(&adverBookshelf).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
	}
	if adverBookshelf.Pic != "" {
		adverBookshelf.Pic = utils.GetAdminFileUrl(adverBookshelf.Pic)
	}

	//开屏广告
	var adverOpenScreen *models.McAdver
	err = global.DB.Model(models.McAdver{}).Order("weight desc").Where("status = 1 and adver_type = 2").First(&adverOpenScreen).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
	}
	if adverOpenScreen.Pic != "" {
		adverOpenScreen.Pic = utils.GetAdminFileUrl(adverOpenScreen.Pic)
	}

	//阅读中
	var adverReading *models.McAdver
	err = global.DB.Model(models.McAdver{}).Order("weight desc").Where("status = 1 and adver_type = 3").First(&adverReading).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
	}
	if adverReading.Pic != "" {
		adverReading.Pic = utils.GetAdminFileUrl(adverReading.Pic)
	}

	//小说详情页
	var adverBookDetail *models.McAdver
	err = global.DB.Model(models.McAdver{}).Order("weight desc").Where("status = 1 and adver_type = 4").First(&adverBookDetail).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
	}
	if adverBookDetail.Pic != "" {
		adverBookDetail.Pic = utils.GetAdminFileUrl(adverBookDetail.Pic)
	}

	//分类下小说列表广告
	var adverClassBookList *models.McAdver
	err = global.DB.Model(models.McAdver{}).Order("weight desc").Where("status = 1 and adver_type = 5").First(&adverClassBookList).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
	}
	if adverClassBookList.Pic != "" {
		adverClassBookList.Pic = utils.GetAdminFileUrl(adverClassBookList.Pic)
	}

	madver = make(map[string]interface{})
	madver["adverBookshelf"] = adverBookshelf
	madver["adverOpenScreen"] = adverOpenScreen
	madver["adverReading"] = adverReading
	madver["adverBookDetail"] = adverBookDetail
	madver["adverClassBookList"] = adverClassBookList
	return
}

func UpdateClickCount(adverValue string) (err error) {
	if adverValue == "" {
		err = fmt.Errorf("%v", "广告id不能为空")
		return
	}
	mdata := make(map[string]interface{})
	mdata["click_count"] = gorm.Expr("click_count + 1")
	err = global.DB.Model(models.McAdver{}).Where("adver_value = ?", adverValue).Updates(&mdata).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

/*
* @note 根据当前ID获取广告信息
* @param id integer 广告ID
* @return object,err
 */
func GetAdverInfoById(id int64) (adver *models.McAdver, err error) {
	err = global.DB.Model(models.McAdver{}).Where("id", id).First(&adver).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

/*
* @note 根据当前的报名来获取小说的基本信息
* @param id int64 包ID
* @return object ,total , err
 */
func GetAdvertPackageInfo(req *models.AdverPackageAqiReq) (adver *models.AdvertProjectInfoRes, err error) {
	packageName := req.PackageName
	deviceType := req.DeviceType
	mark := req.Mark //获取渠道

	db := global.DB.Model(&models.McAdverPackage{})
	if packageName != "" {
		db = db.Where("package_name = ?", packageName)
	}
	if deviceType != "" {
		db = db.Where("device_type = ?", deviceType)
	}
	//获取渠道
	if mark != "" {
		db = db.Where("mark = ?", mark)
	}
	err = db.Debug().Find(&adver).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

// 通过包名反查获取数据信息
func GetPackageByPackageName(package_name string, device string) (info *models.McAdverPackage, err error) {
	if package_name == "" {
		return
	}
	package_name = strings.TrimSpace(package_name) //包名和渠道进行查询
	device = strings.TrimSpace(device)             //渠道设备
	//设置qdh的参数
	db := global.DB.Model(&models.McAdverPackage{}).Debug().Where("package_name = ? and device_type = ? ", package_name, device)
	err = db.Debug().Find(&info).Error
	if err != nil {
		return nil, err
	}
	return info, nil
}

// 通过包名反查获取数据信息
func GetPackageByPackageQdh(device string, package_name string, mark string) (info *models.McAdverPackage, err error) {
	if package_name == "" {
		return
	}
	package_name = strings.TrimSpace(package_name) //包名和渠道进行查询
	device = strings.TrimSpace(device)             //渠道设备
	mark = strings.TrimSpace(mark)
	//设置qdh的参数
	log.Printf("****************************获取到对应的渠道标识11111111111111111 :%s*********************\n", mark)
	db := global.DB.Model(&models.McAdverPackage{}).Debug().Where(" device_type = ? and  package_name = ? and mark = ?", device, package_name, mark)
	err = db.Debug().Find(&info).Error
	if err != nil {
		return nil, err
	}
	return info, nil
}

/*
* @note 获取所有的配置列表信息
* @param id int64 包ID
* @return object ,total , err
 */
func GetAllProjectList(packge_id int64) (list []*models.McAdverProject, err error) {
	if packge_id <= 0 {
		return
	}
	err = global.DB.Model(&models.McAdverProject{}).Where("package_id", packge_id).Debug().Find(&list).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

/*
* @note 广告列表搜索
* @param req object 搜索参数
* @return object ,total , err
 */
func AdverApiSearch(req *models.AdverListReq) (list []*models.McAdver, total int64, err error) {
	db := global.DB.Model(&models.McAdver{}).Order("id desc")

	adverName := strings.TrimSpace(req.AdverName)
	if adverName != "" {
		db = db.Where("adver_name = ?", adverName)
	}

	adverCode := strings.TrimSpace(req.AdverCode)
	if adverCode != "" {
		db = db.Where("adver_code = ?", adverCode)
	}

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	err = db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := req.PageNum
	pageSize := req.PageSize

	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	if len(list) <= 0 {
		return
	}
	//处理广告的图片封面设置
	for _, adver := range list {
		adver.Pic = utils.GetFileUrl(adver.Pic)
	}
	return list, total, err
}

package version_service

import (
	"go-novel/app/models"
	"go-novel/app/service/api/adver_service"
	"go-novel/global"
	"go-novel/utils"
	"log"
	"strings"
)

func GetVersionByDevice(device string) (version *models.McAppVersion, err error) {
	err = global.DB.Model(models.McAppVersion{}).Where("device = ?", device).First(&version).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return
}

func GetVersionByDeviceNew(device string, package_name string) (version *models.McAppVersion, err error) {
	log.Println("11122", device, package_name)

	//查询对应的报名关联的信息 =通过包名+端号+渠道号来关联
	info, err := adver_service.GetPackageByPackageName(package_name, device)
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		//return nil, err
	}
	db := global.DB.Model(models.McAppVersion{})
	packageId := info.Id
	log.Printf("匹配到的应用ID package_id = %d", packageId)
	if packageId > 0 {
		db = db.Where("package_id = ?", packageId)
	} else {
		//没有给一个默认的防止溢出
		db = db.Where("package_id = ?", -1)
	}
	if device != "" {
		db = db.Where("device = ?", device)
	}
	err = db.Debug().Find(&version).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	return version, nil
}

// 通过渠道号来进行计算
func GetVersionByQdh(device string, package_name string, mark string) (version *models.McAppVersion, err error) {
	log.Println("11122", device, package_name, mark)
	if device == "" || package_name == "" || mark == "" {
		return
	}
	//查询对应的报名关联的信息 =通过包名+端号+渠道号来关联
	info, err := adver_service.GetPackageByPackageQdh(device, package_name, mark)
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return nil, err
	}
	db := global.DB.Model(models.McAppVersion{}).Debug()
	packageId := info.Id //包名ID
	log.Printf("匹配到的应用ID package_id = %d", packageId)
	if packageId > 0 {
		db = db.Where("package_id = ?", packageId)
	} else {
		//没有给一个默认的防止溢出
		db = db.Where("package_id = ?", -1)
	}
	if device != "" {
		db = db.Where("device = ?", device)
	}
	err = db.Debug().Find(&version).Error
	if err != nil {
		global.Sqllog.Errorf("%v", err.Error())
		return
	}
	//处理返回的路径
	if version.ForceUrl != "" {
		version.ForceUrl = strings.ReplaceAll(version.ForceUrl, utils.REPLACEAPK, "") //替换对应的字符冗余信息
		//没有地址的时候进行拼接
		if strings.Contains(version.ForceUrl, "https") == false || strings.Contains(version.ForceUrl, "http") == false {
			version.ForceUrl = utils.GetApkFileUrl(version.ForceUrl)
		}
	}
	return version, nil
}

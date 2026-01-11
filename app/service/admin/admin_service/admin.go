package admin_service

import (
	"errors"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-novel/app/models"
	"go-novel/app/service/admin/casbin_adapter_service"
	"go-novel/app/service/admin/menu_service"
	"go-novel/app/service/admin/monitor_service"
	"go-novel/app/service/admin/role_service"
	"go-novel/global"
	"go-novel/utils"
	"log"
	"strconv"
	"strings"
)

// 根据用户id获取管理员信息
func GetAdminById(id int64) (*models.SysAdmin, error) {
	var sysadmin *models.SysAdmin
	err := global.DB.Model(models.SysAdmin{}).Where("id", id).First(&sysadmin).Error
	return sysadmin, err
}

// 根据用户名获取管理员信息
func GetAdminByUsername(username string) (*models.SysAdmin, error) {
	var sysadmin *models.SysAdmin
	err := global.DB.Model(models.SysAdmin{}).Where("username", username).First(&sysadmin).Error
	return sysadmin, err
}

// 获取单前登录用户的信息
func GetCurrentUserInfo(c *gin.Context) (map[string]interface{}, error) {
	username, _ := c.Get("username")
	entity, err := GetAdminByUsername(username.(string))
	if err != nil {
		return nil, err
	}
	//处理获取对应的信息
	if entity.Avatar != "" {
		entity.Avatar = utils.GetAdminFileUrl(entity.Avatar)
	}
	//userInfo,err:=utils.StructToMap(&entity)
	userInfo := structs.Map(&entity)
	delete(userInfo, "password")
	userInfo["roles"] = make([]string, 0)
	allRoles, err := role_service.GetRoleList()
	if err != nil {
		return nil, err
	}
	roles, err := GetAdminRole(entity.Id, allRoles)
	if err != nil {
		return nil, err
	}
	//角色
	userInfo["roles"] = roles
	return userInfo, nil
}

// Login 登录
func Login(login *models.Login, ip string) (token string, expireTime int64, err error) {
	//数据验证
	if len(login.Password) < 6 {
		err = errors.New("密码不能少于6位")
		return
	}

	//判断验证码是否正确
	if !utils.VerifyString(login.IdKeyC, login.IdValueC) {
		err = errors.New("验证码输入错误")
		return
	}

	err = monitor_service.GetLoginError(login.Username)
	if err != nil {
		return
	}

	var sysadmin models.SysAdmin
	//log.Println(login.Username,login.Password)
	//log.Println(db.DB)
	//判断用户名是否存在
	err = global.DB.Model(models.SysAdmin{}).Where("username = ?", login.Username).First(&sysadmin).Error

	if sysadmin.Id == 0 {
		err = errors.New("用户名不存在")
		return
	}
	if err != nil || utils.Md5(login.Password) != sysadmin.Password {
		err = errors.New("密码错误")
		return
	}

	if sysadmin.Status == 0 {
		err = errors.New("用户已被冻结")
		return
	}

	//发送token
	token, expireTime, err = utils.GenerateToken(login.Username, login.Password, 1)
	if err != nil {
		log.Println(token, err)
		return
	}

	//修改最后登录时间和最新ip
	mapUpdate := models.SysAdmin{
		LastLoginTime: utils.GetUnix(),
		LastLoginIp:   ip,
	}

	err = global.DB.Model(models.SysAdmin{}).Where("id = ?", sysadmin.Id).Updates(mapUpdate).Error

	return
}

// 修改密码
func UpdatePwd(req *models.ChangePwdReq, username string) (res bool, err error) {
	adminEntity, err := GetAdminByUsername(username)
	if err != nil || adminEntity.Id <= 0 {
		err = errors.New("获取管理员数据失败")
		return
	}

	if req.OldPassword == "" {
		err = errors.New("旧密码不能为空")
		return
	}

	if utils.Md5(req.OldPassword) != adminEntity.Password {
		err = errors.New("旧密码输入错误")
		return
	}

	if len(req.NewPassword) < 6 {
		err = errors.New("密码长度至少为6位")
		return
	}

	admin := models.SysAdmin{
		Password:   utils.Md5(req.NewPassword),
		UpdateTime: utils.GetUnix(),
	}

	if err = global.DB.Where("id", adminEntity.Id).Updates(&admin).Error; err != nil {
		return false, err
	}

	return true, nil
}

// 修改个人信息
func EditProfile(req *models.EditProfile, username string) (res bool, err error) {
	adminEntity, err := GetAdminByUsername(username)
	if err != nil || adminEntity.Id <= 0 {
		err = errors.New("获取管理员数据失败")
		return
	}

	var mapData = make(map[string]interface{})
	if req.Avatar != "" {
		mapData["avatar"] = req.Avatar
	}

	if req.Nickname != "" {
		mapData["nickname"] = req.Nickname
	}

	if req.Mobile != "" {
		mapData["mobile"] = req.Mobile
	}

	if req.Email != "" {
		mapData["email"] = req.Email
	}

	if req.Remark != "" {
		mapData["remark"] = req.Remark
	}

	if req.Sex >= 0 {
		mapData["sex"] = req.Sex
	}

	if err = global.DB.Model(models.SysAdmin{}).Where("id", adminEntity.Id).Updates(&mapData).Error; err != nil {
		return false, err
	}

	return true, nil
}

// 修改个人头像
func EditAvatar(username, url string) (res bool, err error) {
	var mapData = make(map[string]interface{})
	mapData["avatar"] = url
	if err = global.DB.Model(models.SysAdmin{}).Where("username", username).Updates(&mapData).Error; err != nil {
		return false, err
	}
	return true, nil
}

// Login 登录
func GetUserInfo(username string) (err error, admin models.SysAdmin, roles []string, permissions []string, rolesInfo []*models.SysRole) {
	if username == "" {
		err = errors.New("用户名不能为空")
		return
	}

	err = global.DB.Model(models.SysAdmin{}).Where("username = ?", username).First(&admin).Error
	if err != nil {
		return
	}
	if admin.Id == 0 {
		err = errors.New("用户名不存")
		return
	}
	//获取拼装后台的路径信息
	if admin.Avatar != "" {
		admin.Avatar = utils.GetAdminFileUrl(admin.Avatar)
	}

	//获取用户角色信息
	allRoles, err := role_service.GetRoleList()
	userId := admin.Id

	rolesList := make([]string, 0, 10)
	var permissionsList []string
	if err == nil {
		rolesInfo, err = GetAdminRole(userId, allRoles)
		if err == nil {
			name := make([]string, len(rolesInfo))
			roleIds := make([]uint, len(rolesInfo))
			for k, v := range rolesInfo {
				name[k] = v.Name
				roleIds[k] = uint(v.Id)
			}
			isSuperAdmin := false
			//获取无需验证权限的用户id
			NotCheckAuthAdminIds := viper.GetIntSlice("adminInfo.notCheckAuthAdminIds")
			for _, v := range NotCheckAuthAdminIds {
				if int64(v) == userId {
					isSuperAdmin = true
					break
				}
			}
			if isSuperAdmin {
				permissionsList = []string{"*/*/*"}
			} else {
				permissionsList, err = menu_service.GetPermissionsName(roleIds)
			}
			rolesList = name
		} else {
			return
		}
	} else {
		log.Println(err)
		return
	}
	return nil, admin, rolesList, permissionsList, rolesInfo
}

// 获取管理员的角色信息
func GetAdminRole(adminId int64, allRoleList []*models.SysRole) (roles []*models.SysRole, err error) {
	roleIds, err := GetAdminRoleIds(adminId)
	if err != nil {
		log.Println("ERROR:", err.Error())
		return
	}
	roles = make([]*models.SysRole, 0, len(allRoleList))
	for _, v := range allRoleList {
		for _, id := range roleIds {
			if id == uint(v.Id) {
				roles = append(roles, v)
			}
		}
		if len(roles) == len(roleIds) {
			break
		}
	}
	return
}

// 获取管理员对应的角色ids
func GetAdminRoleIds(adminId int64) (roleIds []uint, err error) {
	groupPolicy, err := casbin_adapter_service.GetAdminRole(adminId)
	if err != nil {
		return
	}
	//查询关联角色规则
	if len(groupPolicy) > 0 {
		roleIds = make([]uint, len(groupPolicy))
		//得到角色id的切片
		for k, v := range groupPolicy {
			countSplit := strings.Split(v[1], "_")
			roleId, _ := strconv.ParseUint(countSplit[1], 10, 64)
			roleIds[k] = uint(roleId)
		}
	}
	return
}

// 获取管理员列表
func AdminListSearch(req *models.AdminListReq) ([]*models.SysAdmin, int64, error) {
	var list []*models.SysAdmin
	db := global.DB.Model(&models.SysAdmin{}).Order("id DESC")

	username := strings.TrimSpace(req.Username)
	if username != "" {
		db = db.Where("username LIKE ?", "%"+username+"%")
	}
	nickname := strings.TrimSpace(req.Nickname)
	if nickname != "" {
		db = db.Where("nickname LIKE ?", "%"+nickname+"%")
	}
	mobile := strings.TrimSpace(req.Mobile)
	if mobile != "" {
		db = db.Where("mobile LIKE ?", "%"+mobile+"%")
	}
	status := strings.TrimSpace(req.Status)
	if status != "" {
		db = db.Where("status = ?", status)
	}
	if req.BeginTime != "" {
		db = db.Where("create_time >=?", utils.DateToUnix(req.BeginTime))
	}
	if req.EndTime != "" {
		db = db.Where("create_time <=?", utils.DateToUnix(req.EndTime))
	}

	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
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
	return list, total, err
}

// 创建管理员
func CreateAdmin(req *models.CreateAdminReq) (InsertId int64, err error) {
	if req.Username == "" {
		err = errors.New("用户名不能为空")
		return
	}

	if req.Nickname == "" {
		req.Nickname = req.Username
	}

	if req.Mobile != "" {
		isMobile := utils.CheckMobile(req.Mobile)
		if isMobile == false {
			err = errors.New("手机号格式不正确")
			return
		}
	}

	if req.Email != "" {
		isEmail := utils.CheckEmail(req.Email)
		if isEmail == false {
			err = errors.New("邮箱格式不正确")
			return
		}
	}

	if req.Password == "" {
		req.Password = "123456"
	}

	if len(req.Password) < 6 {
		err = errors.New("密码长度至少为6位")
		return
	}

	if len(req.RoleIds) <= 0 {
		err = errors.New("角色不能为空")
		return
	}

	admin := models.SysAdmin{
		Username:   req.Username,
		Password:   utils.Md5(req.Password),
		Mobile:     req.Mobile,
		Email:      req.Email,
		Avatar:     req.Avatar,
		Nickname:   req.Nickname,
		Sex:        req.Sex,
		Remark:     req.Remark,
		Status:     req.Status,
		CreateTime: utils.GetUnix(),
	}

	if err = global.DB.Create(&admin).Error; err != nil {
		return 0, err
	} else {
		if len(req.RoleIds) > 0 {
			//设置用户所属角色信息
			err = casbin_adapter_service.AddAdminRole(req.RoleIds, admin.Id)
			if err != nil {
				err = errors.New("设置用户权限失败")
				return
			}
		}
	}

	return admin.Id, nil
}

// 修改管理员
func UpdateAdmin(req *models.UpdateAdminReq) (res bool, err error) {
	if req.Username == "" {
		err = errors.New("用户名不能为空")
		return
	}

	if req.Nickname == "" {
		req.Nickname = req.Username
	}

	if req.Mobile != "" {
		isMobile := utils.CheckMobile(req.Mobile)
		if isMobile == false {
			err = errors.New("手机号格式不正确")
			return
		}
	}

	if req.Email != "" {
		isEmail := utils.CheckEmail(req.Email)
		if isEmail == false {
			err = errors.New("邮箱格式不正确")
			return
		}
	}

	id := req.AdminId

	var mapData = make(map[string]interface{})
	mapData["username"] = req.Username
	mapData["password"] = utils.Md5(req.Password)
	mapData["mobile"] = req.Mobile
	mapData["email"] = req.Email
	mapData["nickname"] = req.Nickname
	mapData["sex"] = req.Sex
	mapData["remark"] = req.Remark
	mapData["status"] = req.Status
	mapData["update_time"] = utils.GetUnix()

	//0值会忽略更新 解决select("*")或者用map
	//admin := models.SysAdmin{
	//	Username:     req.Username,
	//	Password:     utils.Md5(req.Password),
	//	Mobile:       req.Mobile,
	//	Email:        req.Email,
	//	Avatar:       req.Avatar,
	//	Nickname:     req.Nickname,
	//	Sex:          req.Sex,
	//	Remark:       req.Remark,
	//	Status:       req.Status,
	//	UpdateTime:   utils.GetUnix(),
	//}

	if err = global.DB.Model(models.SysAdmin{}).Where("id", id).Updates(&mapData).Error; err != nil {
		return false, err
	} else {
		if len(req.RoleIds) > 0 {
			//设置用户所属角色信息
			err = casbin_adapter_service.EditAdminRole(req.RoleIds, id)
			if err != nil {
				err = errors.New("设置用户权限失败")
				return
			}
		}
	}

	return true, nil
}

// 删除管理员
func DeleteAdmin(req *models.DeleteAdminReq) (res bool, err error) {
	ids := req.AdminIds
	if len(ids) > 0 {
		err = global.DB.Where("id in(?)", ids).Delete(&models.SysAdmin{}).Error
		if err != nil {
			return
		}
	} else {
		err = errors.New("删除失败，参数错误")
		return
	}
	casbin_adapter_service.DeleteAdminRole(ids)
	return true, nil
}

// 重置密码
func ResetAdminPwd(req *models.ResetPwdReq) error {
	//密码加密
	mapData := make(map[string]interface{})
	mapData["password"] = utils.Md5(req.Password)
	if err := global.DB.Model(models.SysAdmin{}).Where("id", req.AdminId).Updates(&mapData).Error; err != nil {
		return err
	}
	return nil
}

// 修改用户状态
func ChangeAdminStatus(req *models.AdminStatusReq) error {
	//密码加密
	mapData := make(map[string]interface{})
	mapData["status"] = req.Status
	if err := global.DB.Model(models.SysAdmin{}).Where("id", req.AdminId).Updates(&mapData).Error; err != nil {
		return err
	}
	return nil
}

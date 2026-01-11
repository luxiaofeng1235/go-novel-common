package role_service

import (
	"errors"
	"go-novel/app/models"
	"go-novel/app/service/admin/casbin_adapter_service"
	"go-novel/global"
	"go-novel/utils"
	"strings"
)

// 获取用户组(角色)列表
func GetRoleList() (list []*models.SysRole, err error) {
	err = global.DB.Model(models.SysRole{}).Order("sort asc,id asc").Find(&list).Error
	return
}

// 根据角色id获取角色信息
func GetRoleById(id int64) (*models.SysRole, error) {
	var sysrole *models.SysRole
	err := global.DB.Model(models.SysRole{}).Where("id", id).First(&sysrole).Error
	return sysrole, err
}

// 获取角色列表
func RoleListSearch(req *models.RoleListReq) ([]*models.SysRole, int64, error) {
	var list []*models.SysRole
	db := global.DB.Model(&models.SysRole{}).Order("id DESC")

	name := strings.TrimSpace(req.Name)
	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
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
	pageNum := int(req.PageNum)
	pageSize := int(req.PageSize)

	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	return list, total, err
}

// 创建角色
func CreateRole(req *models.CreateRoleReq) (InsertId int64, err error) {
	if req.Name == "" {
		err = errors.New("角色名称不能为空")
		return
	}

	role := models.SysRole{
		Name:       req.Name,
		Sort:       req.Sort,
		Remark:     req.Remark,
		Status:     req.Status,
		CreateTime: utils.GetUnix(),
	}

	//tx:=db.DB.Begin()
	if err = global.DB.Create(&role).Error; err != nil {
		//tx.Rollback()
		return 0, err
	} else {
		if len(req.MenuIds) > 0 {
			//添加角色权限
			err = casbin_adapter_service.AddRoleRule(req.MenuIds, role.Id)
			if err != nil {
				err = errors.New("设置角色权限失败")
				//tx.Rollback()
				return
			}
		}
	}
	//tx.Commit()

	return role.Id, nil
}

// 修改角色
func UpdateRole(req *models.UpdateRoleReq) (res bool, err error) {
	//获取path中的userId
	if req.RoleId <= 0 {
		err = errors.New("角色ID不正确")
		return
	}
	if req.Name == "" {
		err = errors.New("角色名称不能为空")
		return
	}
	id := req.RoleId

	//role := models.SysRole{
	//	Name:         req.Name,
	//	Sort:         req.Sort,
	//	Remark:       req.Remark,
	//	Status:       req.Status,
	//	UpdateTime:   utils.GetUnix(),
	//}

	var mapData = make(map[string]interface{})
	mapData["name"] = req.Name
	mapData["sort"] = req.Sort
	mapData["remark"] = req.Remark
	mapData["status"] = req.Status
	mapData["update_time"] = utils.GetUnix()

	//tx:=db.DB.Begin()
	if err = global.DB.Model(models.SysRole{}).Where("id", id).Updates(&mapData).Error; err != nil {
		//tx.Rollback()
		return false, err
	} else {
		if len(req.MenuIds) > 0 {
			//设置用户所属角色信息
			err = casbin_adapter_service.EditRoleRule(req.MenuIds, id)
			if err != nil {
				err = errors.New("设置角色权限失败")
				//tx.Rollback()
				return
			}
		}
	}
	//tx.Commit()
	return true, nil
}

// 删除角色
func DeleteRole(req *models.DeleteRoleReq) (res bool, err error) {
	//tx:=db.DB.Begin()

	ids := req.RoleIds
	if len(ids) > 0 {
		err = global.DB.Where("id in(?)", ids).Delete(&models.SysRole{}).Error
		if err != nil {
			//tx.Rollback()
			return
		}
	} else {
		err = errors.New("删除失败，参数错误")
		return
	}

	//删除角色的权限
	for _, v := range ids {
		//删除对应权限
		err = casbin_adapter_service.DeleteRoleRule(v)
		if err != nil {
			err = errors.New("删除失败")
			//tx.Rollback()
			return
		}
	}
	//tx.Commit()
	return true, nil
}

// 修改角色状态
func ChangeRoleStatus(req *models.RoleStatusReq) error {
	//密码加密
	mapData := make(map[string]interface{})
	mapData["status"] = req.Status
	if err := global.DB.Model(models.SysRole{}).Where("id", req.RoleId).Updates(&mapData).Error; err != nil {
		return err
	}
	return nil
}

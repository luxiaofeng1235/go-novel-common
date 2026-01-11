package menu_service

import (
	"errors"
	"fmt"
	"go-novel/app/models"
	"go-novel/app/service/admin/casbin_adapter_service"
	"go-novel/global"
	"go-novel/utils"
	"strconv"
	"strings"
)

// 根据菜单id获取菜单信息
func GetMenuById(id int64) (*models.SysAuthRule, error) {
	var menu *models.SysAuthRule
	err := global.DB.Model(models.SysAuthRule{}).Where("id", id).First(&menu).Error
	return menu, err
}

// 获取所有菜单
func GetMenuList() (list []*models.SysAuthRule, err error) {
	//从数据库获取
	err = global.DB.Order("weigh desc,id asc").Find(&list).Error
	if err != nil {
		return
	}
	return
}

// 超级管理员获取菜单树
func GetAllMenusTree() (list []*models.SysAuthRule, err error) {
	//获取所有开启的菜单
	allMenus, err := GetIsMenuStatusList()
	if err != nil {
		return
	}
	// parentId为0的是根菜单
	return GenMenuTree(0, allMenus), err
}

func GenMenuTree(parentId int64, menus []*models.SysAuthRule) []*models.SysAuthRule {
	tree := make([]*models.SysAuthRule, 0)

	for _, m := range menus {
		if m.Pid == parentId {
			children := GenMenuTree(m.Id, menus)
			m.Children = children
			tree = append(tree, m)
		}
	}
	return tree
}

// 普通管理员获取菜单树
func GetAdminMenusByRoleIds(roleIds []uint) (list []*models.SysAuthRule, err error) {
	//获取所有开启的菜单
	partMenus, err := GetPermissions(roleIds, 0)
	if err != nil {
		return
	}
	// parentId为0的是根菜单
	return GenMenuTree(0, partMenus), err
}

// 根据角色id获取对应的角色权限 菜单权限mtype=0 按钮=1
func GetPermissions(roleIds []uint, mtype int) (list []*models.SysAuthRule, err error) {
	menuIds, _ := casbin_adapter_service.GetRoleRule(roleIds)
	if len(menuIds) <= 0 {
		return list, nil
	}
	var mList = []*models.SysAuthRule{}
	if mtype == 0 {
		//获取所有的按钮
		mList, err = GetIsMenuStatusList()
		if err != nil {
			return nil, err
		}
	} else {
		//获取所有开启的按钮
		mList, err = GetIsButtonStatusList()
		if err != nil {
			return nil, err
		}
	}

	userButtons := make([]*models.SysAuthRule, 0, 0)
	for _, button := range mList {
		if _, ok := menuIds[int64(button.Id)]; strings.EqualFold(button.Condition, "nocheck") || ok {
			userButtons = append(userButtons, button)
		}
	}
	return userButtons, nil
}

// 检查菜单规则是否存在
func CheckMenuNameUnique(name string, id int64) bool {
	var count int64
	model := global.DB.Model(models.SysAuthRule{}).Where("name=?", name)
	if id != 0 {
		model = model.Where("id!=?", id)
	}
	model.Count(&count)
	return count == 0
}

// 检查菜单路由地址是否已经存在
func CheckMenuPathUnique(path string, id int64) bool {
	var count int64
	model := global.DB.Model(models.SysAuthRule{}).Where("path=?", path).Where("menu_type<>?", 2)
	if id != 0 {
		model = model.Where("id!=?", id)
	}
	model.Count(&count)
	return count == 0
}

// 获取MenuType==0,1菜单列表
func GetIsMenuList() ([]*models.SysAuthRule, error) {
	list, err := GetMenuList()
	if err != nil {
		return nil, err
	}
	var gList = make([]*models.SysAuthRule, 0, len(list))
	for _, v := range list {
		if v.MenuType == 0 {
			v.Component = "Layout"
			v.Path = "/" + v.Path
		}
		if v.AlwaysShow == 1 && v.MenuType == 0 {
			v.AlwaysShow = 1
		} else {
			v.AlwaysShow = 0
		}
		if v.MenuType == 0 || v.MenuType == 1 {
			gList = append(gList, v)
		}
	}
	return gList, nil
}

// 获取isMenu=0|1且status=1的菜单列表
func GetIsMenuStatusList() ([]*models.SysAuthRule, error) {
	list, err := GetMenuList()
	if err != nil {
		return nil, err
	}
	var gList = make([]*models.SysAuthRule, 0, len(list))
	for _, v := range list {
		if v.MenuType == 0 {
			v.Component = "Layout"
			v.Path = "/" + v.Path
		}
		if v.AlwaysShow == 1 && v.MenuType == 0 {
			v.AlwaysShow = 1
		} else {
			v.AlwaysShow = 0
		}
		if (v.MenuType == 0 || v.MenuType == 1) && v.Status == 1 {
			gList = append(gList, v)
		}
	}
	return gList, nil
}

// 获取所有按钮isMenu=2 且status=1的菜单列表
func GetIsButtonStatusList() ([]*models.SysAuthRule, error) {
	list, err := GetMenuList()
	if err != nil {
		return nil, err
	}
	var gList = make([]*models.SysAuthRule, 0, len(list))
	for _, v := range list {
		if v.MenuType == 2 && v.Status == 1 {
			gList = append(gList, v)
		}
	}
	return gList, nil
}

// 获取status==1的菜单列表
func GetMenuIsStatusList() ([]*models.SysAuthRule, error) {
	list, err := GetMenuList()
	if err != nil {
		return nil, err
	}
	var gList = make([]*models.SysAuthRule, 0, len(list))
	for _, v := range list {
		if v.MenuType == 0 {
			v.Component = "Layout"
			v.Path = "/" + v.Path
		}
		if v.AlwaysShow == 1 && v.MenuType == 0 {
			v.AlwaysShow = 1
		} else {
			v.AlwaysShow = 0
		}
		if v.Status == 1 {
			gList = append(gList, v)
		}
	}
	return gList, nil
}

// 获取角色名称 获取所欲按钮权限
func GetPermissionsName(roleIds []uint) (userButtons []string, err error) {
	Permissions, err := GetPermissions(roleIds, 1)
	if err != nil {
		return userButtons, err
	}
	userButtons = make([]string, 0)
	for _, permission := range Permissions {
		if permission.Name != "" {
			userButtons = append(userButtons, permission.Name)
		}
	}
	return userButtons, nil
}

func GetMenuListSearch(req *models.SysAuthRuleReqSearch) (list []*models.SysAuthRule, err error) {
	list, err = GetMenuList()
	if err != nil {
		return
	}
	if req != nil {
		tmpList := make([]*models.SysAuthRule, 0, len(list))
		for _, entity := range list {
			status, _ := strconv.Atoi(req.Status)
			if (req.Title == "" || strings.Contains(strings.ToUpper(entity.Title), strings.ToUpper(req.Title))) &&
				(req.Status == "" || status == entity.Status) {
				tmpList = append(tmpList, entity)
			}
		}
		list = tmpList
	}
	return
}

// 获取菜单
func GetMapMenus() (menus []map[string]interface{}, err error) {
	//获取所有开启的菜单
	allMenus, err := GetIsMenuStatusList()
	if err != nil {
		return
	}
	menus = make([]map[string]interface{}, len(allMenus))
	for k, v := range allMenus {
		menu := make(map[string]interface{})
		menu = setMenuMap(menu, v)
		menus[k] = menu
	}
	menus = PushSonToParent(menus, 0, "pid", "children", true)
	return
}

// 组合返回menu前端数据
func setMenuMap(menu map[string]interface{}, entity *models.SysAuthRule) map[string]interface{} {
	menu["id"] = entity.Id
	menu["pid"] = entity.Pid
	menu["index"] = entity.Name
	menu["name"] = utils.FirstUpper(entity.Path)
	menu["menuName"] = entity.Title
	if entity.MenuType != 0 {
		menu["component"] = entity.Component
		menu["path"] = entity.Path
	} else {
		menu["path"] = "/" + entity.Path
		menu["component"] = "Layout"
	}
	menu["meta"] = map[string]string{
		"icon":  entity.Icon,
		"title": entity.Title,
	}
	if entity.AlwaysShow == 1 {
		menu["hidden"] = false
	} else {
		menu["hidden"] = true
	}
	if entity.AlwaysShow == 1 && entity.MenuType == 0 {
		menu["alwaysShow"] = true
	} else {
		menu["alwaysShow"] = false
	}
	return menu
}

//有层级关系的数组 ,将子级压入到父级（树形结构）
/*
pid         int      父级id
pidField    string   父级id键名
childFied   string   子级数组键名
showNoChild bool     是否显示不存在的子级健
*/
func PushSonToParent(list []map[string]interface{}, pid int64, pidField string, childFied string, showNoChild bool) []map[string]interface{} {
	var returnList []map[string]interface{}
	for _, v := range list {
		if v[pidField] == pid {
			titlePrefix := ""
			titlePrefix = "├" + titlePrefix
			v["title_prefix"] = titlePrefix
			v["title_show"] = fmt.Sprintf("%s%s", v["title_prefix"], v["name"])
			child := PushSonToParent(list, v["id"].(int64), pidField, childFied, showNoChild)
			if child != nil || showNoChild {
				v[childFied] = child
			}
			returnList = append(returnList, v)
		}
	}
	return returnList
}

// 获取管理员列表
func MenuListSearch(req *models.MenuListReq) ([]*models.SysAuthRule, int64, error) {
	var list []*models.SysAuthRule
	db := global.DB.Model(&models.SysAuthRule{}).Order("weigh desc,id asc")

	title := strings.TrimSpace(req.Title)
	if title != "" {
		db = db.Where("title like ?", "%"+title+"%")
	}

	name := strings.TrimSpace(req.Name)
	if name != "" {
		db = db.Where("name like ?", "%"+name+"%")
	}

	status := strings.TrimSpace(req.Status)
	if status != "" {
		db = db.Where("status = ?", status)
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

// 创建菜单
func CreateMenu(req *models.CreateMenuReq) (InsertId int64, err error) {
	if req.Name == "" {
		err = errors.New("菜单规则名称不能为空")
		return
	}

	//判断菜单规则是否存在
	if !CheckMenuNameUnique(req.Name, 0) {
		err = errors.New("菜单规则名称已经存在")
		return
	}

	//判断路由是否已经存在
	if !CheckMenuPathUnique(req.Path, 0) {
		err = errors.New("路由地址已经存在")
		return
	}

	entity := new(models.SysAuthRule)
	entity.Title = req.Title
	entity.Status = req.Status
	entity.MenuType = req.MenuType
	entity.Path = req.Path
	entity.Component = req.Component
	entity.AlwaysShow = req.AlwaysShow
	entity.Icon = req.Icon
	entity.Name = req.Name
	entity.IsFrame = req.IsFrame
	entity.Pid = req.Pid
	entity.CreateTime = utils.GetUnix()
	entity.Weigh = req.Weigh

	if err = global.DB.Create(&entity).Error; err != nil {
		return 0, err
	}

	return entity.Id, nil
}

// 修改菜单
func UpdateMenu(req *models.UpdateMenuReq) (res bool, err error) {
	//获取path中的userId
	if req.MenuId <= 0 {
		err = errors.New("菜单ID不正确")
		return
	}
	if req.Name == "" {
		err = errors.New("菜单规则名称不能为空")
		return
	}
	id := req.MenuId
	//判断菜单规则是否存在
	if !CheckMenuNameUnique(req.Name, id) {
		err = errors.New("菜单规则名称已经存在")
		return
	}

	//判断路由是否已经存在
	if !CheckMenuPathUnique(req.Path, id) {
		err = errors.New("路由地址已经存在")
		return
	}

	var mapData = make(map[string]interface{})
	mapData["title"] = req.Title
	mapData["icon"] = req.Icon
	mapData["status"] = req.Status
	mapData["menu_type"] = req.MenuType
	mapData["path"] = req.Path
	mapData["component"] = req.Component
	mapData["always_show"] = req.AlwaysShow
	mapData["name"] = req.Name
	mapData["is_frame"] = req.IsFrame
	mapData["pid"] = req.Pid
	mapData["update_time"] = utils.GetUnix()
	mapData["weigh"] = req.Weigh

	//entity := new(models.SysAuthRule)
	//entity.Title = req.Title
	//entity.Status = req.Status
	//entity.MenuType = req.MenuType
	//entity.Path = req.Path
	//entity.Component = req.Component
	//entity.AlwaysShow = req.AlwaysShow
	//entity.Icon = req.Icon
	//entity.Name = req.Name
	//entity.IsFrame = req.IsFrame
	//entity.Pid = req.Pid
	//entity.UpdateTime = utils.GetUnix()
	//entity.Weigh = req.Weigh

	if err = global.DB.Model(models.SysAuthRule{}).Where("id", id).Updates(&mapData).Error; err != nil {
		return false, err
	}

	return true, nil
}

// 删除菜单
func DeleteMenu(req *models.DeleteMenuReq) (res bool, err error) {
	ids := req.MenuIds

	//获取菜单数据
	menus, err := GetMenuList()
	if err != nil {
		return
	}

	//将数据转为map数组
	var mList []map[string]interface{}
	for _, entity := range menus {
		m := map[string]interface{}{
			"id":   entity.Id,
			"pid":  entity.Pid,
			"name": entity.Title,
		}
		mList = append(mList, m)
	}

	//递归遍历获取所有子数据
	son := make([]map[string]interface{}, 0, len(menus))
	for _, id := range ids {
		son = append(son, utils.FindSonByParentId(mList, id, "pid", "id")...)
	}

	//取出所有的子数据id
	for _, v := range son {
		ids = append(ids, v["id"].(int64))
	}

	if len(ids) > 0 {
		err = global.DB.Where("id in(?)", ids).Delete(&models.SysAuthRule{}).Error
		if err == nil {
			res = true
		}
	} else {
		err = errors.New("删除失败，参数错误")
	}
	return
}

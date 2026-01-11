package casbin_adapter_service

import "fmt"

//查询管理员角色信息
func GetAdminRole(adminId int64) (ids [][]string,err error) {
	enforcer, err := GetEnforcer()
	if err != nil{
		return
	}
	//查询关联角色规则
	groupPolicy := enforcer.GetFilteredGroupingPolicy(0, fmt.Sprintf("u_%d", adminId))
	return groupPolicy,nil
}
//添加管理员角色信息
func AddAdminRole(roleIds []int, adminId int64) (err error) {
	enforcer, e := GetEnforcer()
	if e != nil {
		err = e
		return
	}
	for _, v := range roleIds {
		_, err = enforcer.AddGroupingPolicy(fmt.Sprintf("u_%d", adminId), fmt.Sprintf("g_%d", v))
		if err != nil {
			return
		}
	}
	return
}

//修改用户角色信息
func EditAdminRole(roleIds []int, adminId int64) (err error) {
	enforcer, e := GetEnforcer()
	if e != nil {
		err = e
		return
	}
	//删除用户旧角色信息
	enforcer.RemoveFilteredGroupingPolicy(0, fmt.Sprintf("u_%d", adminId))
	for _, v := range roleIds {
		_, err = enforcer.AddGroupingPolicy(fmt.Sprintf("u_%d", adminId), fmt.Sprintf("g_%d", v))
		if err != nil {
			return
		}
	}
	return
}
//添加管理员角色信息
func DeleteAdminRole(adminIds []int64){
	//删除对应权限
	enforcer, err := GetEnforcer()
	if err == nil {
		for _, v := range adminIds {
			enforcer.RemoveFilteredGroupingPolicy(0, fmt.Sprintf("u_%d", v))
		}
	}
}
package casbin_adapter_service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//根据角色id获取对应的角色权限
func GetRoleRule(roleIds []uint) (map[int64]int64, error) {
	//获取角色对应的菜单id
	enforcer, err := GetEnforcer()
	if err != nil {
		return nil, err
	}
	menuIds := map[int64]int64{}
	for _, roleId := range roleIds {
		//查询当前权限
		gp := enforcer.GetFilteredPolicy(0, fmt.Sprintf("g_%d", roleId))
		if len(gp) <= 0{
			return menuIds,nil
		}
		for _, p := range gp {
			countSplit := strings.Split(p[1], "_")
			mid, _ := strconv.ParseInt(countSplit[1], 10, 64)
			menuIds[mid] = mid
		}
	}
	return menuIds,nil
}

func GetRoleIdRule(roleId int) (arr []int,err error) {
	//获取角色关联的菜单规则
	enforcer, err := GetEnforcer()
	if err != nil {
		return
	}
	gp := enforcer.GetFilteredNamedPolicy("p", 0, fmt.Sprintf("g_%d", roleId))
	gpSlice := make([]int, len(gp))
	for k, v := range gp {
		countSplit := strings.Split(v[1], "_")
		mid, _ := strconv.ParseInt(countSplit[1], 10, 64)
		gpSlice[k] = int(mid)
	}
	return gpSlice,nil
}

//添加角色授权规则
func AddRoleRule(iRule []int, roleId int64) (err error) {
	if len(iRule) <= 0{
		return errors.New("规则数组为空")
	}
	enforcer, e := GetEnforcer()
	if e != nil {
		err = e
		return
	}
	//strRule, err := json.Marshal(iRule)
	//if err != nil {
	//	return
	//}
	//rule:=strings.Split(string(strRule),",")
	go func() {
		for _, v := range iRule {
			_, err = enforcer.AddPolicy(fmt.Sprintf("g_%d", roleId), fmt.Sprintf("r_%d", v), "All")
			if err != nil {
				return
			}
		}
	}()
	return nil
}

//修改角色的授权规则
func EditRoleRule(iRule []int, roleId int64) (err error) {
	enforcer, e := GetEnforcer()
	if e != nil {
		err = e
		return
	}
	if len(iRule) <= 0{
		return errors.New("规则数组为空")
	}
	//strRule, err := json.Marshal(iRule)
	//if err != nil {
	//	return
	//}
	//rule:=strings.Split(string(strRule),",")
	go func() {
		//查询当前权限
		gp := enforcer.GetFilteredPolicy(0, fmt.Sprintf("g_%d", roleId))
		//删除旧权限
		for _, v := range gp {
			_, e = enforcer.RemovePolicy(v)
			if e != nil {
				err = e
				return
			}
		}
		for _, v := range iRule {
			_, err = enforcer.AddPolicy(fmt.Sprintf("g_%d", roleId), fmt.Sprintf("r_%d", v), "All")
			if err != nil {
				return
			}
		}
	}()
	return
}

//删除角色权限操作
func DeleteRoleRule(roleId int64) (err error) {
	enforcer, e := GetEnforcer()
	if e != nil {
		err = e
		return
	}
	//查询当前权限
	gp := enforcer.GetFilteredNamedPolicy("p", 0, fmt.Sprintf("g_%d", roleId))
	//删除旧权限
	for _, v := range gp {
		_, e = enforcer.RemovePolicy(v)
		if e != nil {
			err = e
			return
		}
	}
	return
}
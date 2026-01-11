package middleware

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go-novel/app/models"
	"go-novel/app/service/admin/admin_service"
	"go-novel/app/service/admin/casbin_adapter_service"
	"go-novel/app/service/admin/menu_service"
	"go-novel/utils"
	"strings"
)

// 权限判断处理中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessParams := c.GetStringSlice("accessParams")
		accessParamsStr := ""
		if len(accessParams) > 0 && accessParams[0] != "undefined" {
			accessParamsStr = "?" + strings.Join(accessParams, "&")
		}
		url := strings.TrimLeft(c.FullPath(), "/") + accessParamsStr

		username, _ := c.Get("username")
		if username == nil {
			return
		}

		admininfo, err := admin_service.GetAdminByUsername(username.(string))
		if err != nil {
			utils.Fail(c, err, "获取用户信息失败")
			return
		}
		adminId := admininfo.Id
		//获取无需验证权限的用户id
		NotCheckAuthAdminIds := viper.GetIntSlice("adminInfo.notCheckAuthAdminIds")
		for _, v := range NotCheckAuthAdminIds {
			if int64(v) == adminId {
				c.Next()
				return
			}
		}

		//获取地址对应的菜单id
		menuList, err := menu_service.GetMenuIsStatusList()
		if err != nil {
			utils.Fail(c, err, "请求数据失败")
			c.Abort()
			return
		}
		var menu *models.SysAuthRule
		for _, m := range menuList {
			ms := utils.SubStr(m.Name, 0, strings.Index(m.Name, "?"))
			if m.Name == url || ms == url {
				menu = m
				break
			}
		}
		//只验证存在数据库中的规则
		if menu != nil {
			//若存在不需要验证的条件则跳过
			if strings.Contains(menu.Condition, "nocheck") {
				c.Next()
				return
			}
			menuId := menu.Id
			if menuId != 0 {
				//判断权限操作
				var enforcer *casbin.SyncedEnforcer
				enforcer, err = casbin_adapter_service.GetEnforcer()
				if err != nil {
					utils.Fail(c, err, "获取权限失败")
					c.Abort()
					return
				}
				groupPolicy := enforcer.GetFilteredGroupingPolicy(0, fmt.Sprintf("u_%d", adminId))

				if len(groupPolicy) == 0 {
					utils.Fail(c, err, "没有访问权限")
					c.Abort()
					return
				}
				hasAccess := false
				for _, v := range groupPolicy {
					if enforcer.HasPolicy(v[1], fmt.Sprintf("r_%d", menuId), "All") {
						hasAccess = true
						break
					}
				}
				if !hasAccess {
					utils.Fail(c, err, "没有访问权限")
					c.Abort()
					return
				}
			}
		} else if menu == nil && accessParamsStr != "" {
			utils.Fail(c, err, "没有访问权限")
			c.Abort()
			return
		}

		c.Next()
	}
}

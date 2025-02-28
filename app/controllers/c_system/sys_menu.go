package c_system

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"hrkGo/app/model/sys_model"
	"hrkGo/app/service/sys_service"
	"hrkGo/utils/global/consts"
	"hrkGo/utils/global/variable"
	"hrkGo/utils/response"
	"strconv"
	"time"
)

type MenuController struct {
	MenuService sys_service.MenuService
	RoleService sys_service.RoleService
}

// GetRouters 获取全部路由
func (m *MenuController) GetRouters(c *gin.Context) {
	trees, err := m.MenuService.GetMenuTreeAll()
	RouterVo := m.MenuService.BuildMenus(trees)
	if err != nil {

	}
	response.Success(c, "查询成功", RouterVo)
}

// TreeSelect 菜单下拉列表
func (m *MenuController) TreeSelect(c *gin.Context) {
	var menu sys_model.SysMenu
	menus, err := m.MenuService.SelectMenuList(menu, c.Keys["userId"].(string))
	if err != nil {
		response.BusinessFail(c, "查询失败，请稍后再试")
	}
	TreeSelect := m.MenuService.BuildMenuTree(menus)
	response.Success(c, "查询成功", TreeSelect)
}

// SelectMenuList 查询系统菜单列表
func (m *MenuController) SelectMenuList(c *gin.Context) {

	var menu sys_model.SysMenu
	if err := c.ShouldBindQuery(&menu); err != nil {
		response.ValidateFail(c, response.GetErrorMsg(menu, err))
		return
	}

	menus, err := m.MenuService.SelectMenuList(menu, c.Keys["userId"].(string))
	if err != nil {
		response.BusinessFail(c, "查询失败，请稍后再试")
	}
	response.Success(c, "查询成功", menus)
}

// RoleMenuTreeSelect 查询系统菜单列表
func (m *MenuController) RoleMenuTreeSelect(c *gin.Context) {
	roleId := c.Param("roleId")

	role, _ := m.RoleService.SelectRoleDataById(roleId)

	menuIds, _ := m.MenuService.SelectMenuListByRole(role)

	var menu sys_model.SysMenu
	menus, _ := m.MenuService.SelectMenuList(menu, c.Keys["userId"].(string))

	TreeSelect := map[string]interface{}{
		"checkedKeys": menuIds,
		"menus":       m.MenuService.BuildMenuTree(menus),
	}
	response.Success(c, "查询成功", TreeSelect)

}

// SelectMenuById 查询系统菜单列表
func (m *MenuController) SelectMenuById(c *gin.Context) {
	menuId := c.Param("menuId")
	dict, err := m.MenuService.SelectMenuById(menuId)
	if err != nil {
		response.BusinessFail(c, consts.SQLERROR)
		return
	}
	response.Success(c, "查询成功", dict)
}

// InsertMenu 新增菜单
func (m *MenuController) InsertMenu(c *gin.Context) {

	var form sys_model.SysMenu

	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, response.GetErrorMsg(form, err))
		return
	}

	exists := m.MenuService.CheckMenuNameUnique(form)

	if exists {
		response.BusinessFail(c, "菜单名称已存在")
		return
	} else if !m.MenuService.IsFace(form) {
		response.BusinessFail(c, "新增菜单失败，地址必须以http(s)://开头")
		return
	}
	// 方法二：
	now := time.Now()

	menu := &sys_model.SysMenu{
		MenuId:     variable.SnowFlake.GetIdStr(),
		Icon:       form.Icon,
		IsCache:    form.IsCache,
		IsFrame:    form.IsFrame,
		MenuName:   form.MenuName,
		MenuType:   form.MenuType,
		OrderNum:   form.OrderNum,
		ParentId:   form.ParentId,
		Perms:      form.Perms,
		Query:      form.Query,
		Component:  form.Component,
		Path:       form.Path,
		Status:     form.Status,
		Visible:    form.Visible,
		CreateTime: &now,
		CreateBy:   c.Keys["userId"].(string),
	}

	err := m.MenuService.InsertMenu(menu)
	if err != nil {
		response.BusinessFail(c, "新增失败,请重试")
		return
	}
	response.SuccessNil(c, "新增成功")
}

// UpdateMenu 新增菜单
func (m *MenuController) UpdateMenu(c *gin.Context) {

	var form sys_model.SysMenu

	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, response.GetErrorMsg(form, err))
		return
	}
	exists := m.MenuService.CheckMenuNameUnique(form)

	if exists {
		response.BusinessFail(c, "菜单名称已存在")
		return
	} else if !m.MenuService.IsFace(form) {
		response.BusinessFail(c, "新增菜单失败，地址必须以http(s)://开头")
		return
	}

	// 方法二：
	now := time.Now()

	menu := &sys_model.SysMenu{
		MenuId:     form.MenuId,
		Icon:       form.Icon,
		IsCache:    form.IsCache,
		IsFrame:    form.IsFrame,
		MenuName:   form.MenuName,
		MenuType:   form.MenuType,
		OrderNum:   form.OrderNum,
		ParentId:   form.ParentId,
		Perms:      form.Perms,
		Query:      form.Query,
		Component:  form.Component,
		Path:       form.Path,
		Status:     form.Status,
		Visible:    form.Visible,
		UpdateBy:   strconv.FormatUint(uint64(c.Keys["userId"].(uint)), 10),
		UpdateTime: &now,
	}

	err := m.MenuService.UpdateMenu(menu)
	if err != nil {
		fmt.Println(err)
		response.BusinessFail(c, "修改失败,请重试")

		return
	}
	response.SuccessNil(c, "操作成功")
}

// DeleteMenuById 通过id删除
func (m *MenuController) DeleteMenuById(c *gin.Context) {
	menuId := c.Param("menuId")
	if m.MenuService.HasChildByMenuId(menuId) {
		response.BusinessFail(c, "存在子菜单,不允许删除")
		return
	}
	if m.MenuService.CheckMenuExistRole(menuId) {
		response.BusinessFail(c, "菜单已分配,不允许删除")
		return
	}

	// 获取路径参数

	err := m.MenuService.DeleteMenuById(menuId)
	if err != nil {
		response.BusinessFail(c, "删除失败,请重试")
		return
	}
	response.SuccessNil(c, "操作成功")
}

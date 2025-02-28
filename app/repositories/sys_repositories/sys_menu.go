package sys_repositories

import (
	"errors"
	"fmt"
	"hrkGo/app/model/sys_model"
	"hrkGo/utils/global/variable"
	"time"
)

// MenuCrud 实现 IMenuCrud 接口
type menuCrud struct{}

var MenuCrud = new(menuCrud)

// SelectMenuList 获取菜单列表
func (m *menuCrud) SelectMenuList(menu sys_model.SysMenu) ([]*sys_model.SysMenu, error) {
	var menuList []*sys_model.SysMenu
	query := variable.GormDbMysql.Model(&sys_model.SysMenu{})
	// 构建查询条件
	if menu.MenuName != "" {
		query = query.Where("menu_name LIKE ?", "%"+menu.MenuName+"%")
	}
	if menu.Status != "" {
		query = query.Where("status = ?", menu.Status)
	}

	// 添加排序
	err := query.Order("parent_id, order_num").Find(&menuList).Error
	return menuList, err
}

// SelectMenuListByUserId 根据userId获取菜单列表
func (m *menuCrud) SelectMenuListByUserId(menu sys_model.SysMenu) ([]*sys_model.SysMenu, error) {
	var menuList []*sys_model.SysMenu

	query := variable.GormDbMysql.Model(&sys_model.SysMenu{}).
		Joins("left join sys_role_menu rm on sys_menu.menu_id = rm.menu_id").
		Joins("left join sys_user_role ur on rm.role_id = ur.role_id").
		Joins("left join sys_role ro on ur.role_id = ro.role_id").
		Where("ur.user_id = ?", menu.Params["userId"])

	if menu.MenuName != "" {
		query = query.Where("menu_name LIKE ?", "%"+menu.MenuName+"%")
	}
	if menu.Status != "" {
		query = query.Where("status = ?", menu.Status)
	}

	err := query.Order("parent_id, order_num").Find(&menuList).Error
	return menuList, err
}

// SelectMenuListByRole 通过role查询菜单 错的
func (m *menuCrud) SelectMenuListByRole(role sys_model.SysRole) (menuIds []string, err error) {
	err = variable.GormDbMysql.Table("sys_menu m").
		Select("m.menu_id").
		Joins("left join sys_role_menu rm on m.menu_id = rm.menu_id").
		Where("rm.role_id = ?", role.RoleId).
		Order("m.parent_id, m.order_num").
		Pluck("m.menu_id", &menuIds). // Pluck the menu_id into the menuIds slice
		Error

	if err != nil {
		return nil, err
	}
	return menuIds, nil
}

// SelectMenuById 通过id查询菜单
func (m *menuCrud) SelectMenuById(menuId string) (menuData sys_model.SysMenu, err error) {
	err = variable.GormDbMysql.Where("menu_id = ?", menuId).First(&menuData).Error
	return menuData, err
}

// HasChildByMenuId 查询有没有子菜单
func (m *menuCrud) HasChildByMenuId(menuId string) (count int64) {
	variable.GormDbMysql.Model(&sys_model.SysMenu{}).Where("parent_id = ?", menuId).Count(&count)
	return count
}

// CheckMenuNameUnique 查询是否有同名菜单
func (m *menuCrud) CheckMenuNameUnique(menu sys_model.SysMenu) (menuData sys_model.SysMenu, exists bool) {

	result := variable.GormDbMysql.Where("menu_name = ? AND parent_id = ?", menu.MenuName, menu.ParentId).
		First(&menuData)

	if result.RowsAffected == 0 {
		return menuData, false
	}
	return menuData, true
}

// InsertMenu 新建菜单
func (m *menuCrud) InsertMenu(menu *sys_model.SysMenu) (err error) {
	return variable.GormDbMysql.Create(menu).Error
}

// DeleteMenuById 删除菜单
func (m *menuCrud) DeleteMenuById(menuId string) (err error) {
	return variable.GormDbMysql.Where("menu_id = ?", menuId).Delete(&sys_model.SysMenu{}).Error
}

// UpdateMenu 修改菜单
func (m *menuCrud) UpdateMenu(menu *sys_model.SysMenu) (err error) {
	// 创建一个map来存储需要更新的字段
	updates := make(map[string]interface{})
	fmt.Println("1111111111111111")
	// 只更新非空字段
	if menu.MenuName != "" {
		updates["menu_name"] = menu.MenuName
	}
	if menu.ParentId != "" {
		updates["parent_id"] = menu.ParentId
	}
	if menu.OrderNum != 0 {
		updates["order_num"] = menu.OrderNum
	}
	if menu.Path != "" {
		updates["path"] = menu.Path
	}
	if menu.Component != "" {
		updates["component"] = menu.Component
	}
	if menu.Query != "" {
		updates["query"] = menu.Query
	}
	if menu.IsFrame != "" {
		updates["is_frame"] = menu.IsFrame
	}
	if menu.IsCache != "" {
		updates["is_cache"] = menu.IsCache
	}
	if menu.MenuType != "" {
		updates["menu_type"] = menu.MenuType
	}
	if menu.Visible != "" {
		updates["visible"] = menu.Visible
	}
	if menu.Status != "" {
		updates["status"] = menu.Status
	}
	if menu.Perms != "" {
		updates["perms"] = menu.Perms
	}
	if menu.Icon != "" {
		updates["icon"] = menu.Icon
	}
	if menu.UpdateBy != "" {
		updates["update_by"] = menu.UpdateBy
	}
	if menu.Remark != "" {
		updates["remark"] = menu.Remark
	}

	// 更新时间总是会更新
	now := time.Now()
	updates["update_time"] = &now

	// 执行更新操作，只更新非零值字段
	result := variable.GormDbMysql.Model(&sys_model.SysMenu{}).
		Where("menu_id = ?", menu.MenuId).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	// 如果没有记录被更新
	if result.RowsAffected == 0 {
		return errors.New("menu not found")
	}
	return err
}

// GetMenuTreeAll 查询所有树状图
func (m *menuCrud) GetMenuTreeAll() ([]*sys_model.SysMenu, error) {
	var menus []*sys_model.SysMenu
	// 一次性查询所有菜单
	err := variable.GormDbMysql.Where("status = '0' and menu_type != ?", "F").
		Order("order_num").
		Find(&menus).Error
	return menus, err
}

// SelectMenuPermsByUserId 查询所有树状图
func (m *menuCrud) SelectMenuPermsByUserId(userId string) ([]string, error) {
	var perms []string

	err := variable.GormDbMysql.Model(&sys_model.SysMenu{}).
		Distinct("m.perms").
		Select("sys_menu.perms").
		Joins("LEFT JOIN sys_role_menu rm ON sys_menu.menu_id = rm.menu_id").
		Joins("LEFT JOIN sys_user_role ur ON rm.role_id = ur.role_id").
		Joins("LEFT JOIN sys_role r ON r.role_id = ur.role_id").
		Where("sys_menu.status = ? AND r.status = ? AND ur.user_id = ?", "0", "0", userId).
		Find(&perms).Error

	return perms, err
}

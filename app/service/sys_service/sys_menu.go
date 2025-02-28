package sys_service

import (
	"hrkGo/app/model/sys_model"
	"hrkGo/app/repositories/sys_repositories"
	"hrkGo/utils/StringUtils"
	"hrkGo/utils/global/consts"
	"strings"
)

type MenuService struct {
}

// SelectMenuList 获取菜单下拉树列表
func (m *MenuService) SelectMenuList(menu sys_model.SysMenu, userId string) (menuList []*sys_model.SysMenu, err error) {
	// 管理员显示所有菜单信息
	if sys_model.IsAdmin(userId) {
		menuList, err = sys_repositories.MenuCrud.SelectMenuList(menu)
	} else {
		menu.Params = make(map[string]interface{})
		menu.Params["userId"] = userId
		menuList, err = sys_repositories.MenuCrud.SelectMenuListByUserId(menu)
	}
	return menuList, err
}

// SelectMenuListByRole 根据角色ID查询菜单树信息
func (m *MenuService) SelectMenuListByRole(role sys_model.SysRole) (menuIds []string, err error) {
	menuIds, err = sys_repositories.MenuCrud.SelectMenuListByRole(role)
	return menuIds, err
}

// SelectMenuById 获取菜单下拉树列表
func (m *MenuService) SelectMenuById(menuId string) (menuData sys_model.SysMenu, err error) {
	menuData, err = sys_repositories.MenuCrud.SelectMenuById(menuId)
	return menuData, err
}

// HasChildByMenuId 查询是否有子菜单
func (m *MenuService) HasChildByMenuId(menuId string) bool {
	count := sys_repositories.MenuCrud.HasChildByMenuId(menuId)
	return count > 0
}

// CheckMenuExistRole 查询菜单是否存在角色
//
// Parameters:
//   - menuId: 菜单ID
//
// Returns:
//   - bool: true表示菜单被角色使用，false表示菜单未被角色使用
func (m *MenuService) CheckMenuExistRole(menuId string) bool {
	count := sys_repositories.RomeMenuCrud.CheckMenuExistRole(menuId)
	return count > 0
}

// InsertMenu 新增菜单
func (m *MenuService) InsertMenu(menu *sys_model.SysMenu) (err error) {
	return sys_repositories.MenuCrud.InsertMenu(menu)
}

// CheckMenuNameUnique 查询是否有同名菜单
func (m *MenuService) CheckMenuNameUnique(menu sys_model.SysMenu) (exists bool) {
	menuData, exists := sys_repositories.MenuCrud.CheckMenuNameUnique(menu)
	if exists && menuData.MenuId == menu.MenuId {
		return false
	}
	return exists
}

// UpdateMenu 修改
func (m *MenuService) UpdateMenu(menu *sys_model.SysMenu) (err error) {
	return sys_repositories.MenuCrud.UpdateMenu(menu)
}

// DeleteMenuById 删除
func (m *MenuService) DeleteMenuById(menuId string) (err error) {
	return sys_repositories.MenuCrud.DeleteMenuById(menuId)
}

//// BuildMenuTreeSelect 下拉树结构列表
//func (m *MenuService) BuildMenuTreeSelect(menus []*sys_model.SysMenu) (menuList []*sys_model.TreeSelect) {
//	//menuTrees := buildMenuTree(menus)
//	//return menuTrees
//}

// LoginSelectMenuPermsByUserId 用户权限列表
func (m *MenuService) LoginSelectMenuPermsByUserId(userId string) (permsSet []string) {
	// 管理员拥有所有权限
	if sys_model.IsAdmin(userId) {
		permsSet = append(permsSet, "*:*:*")
	} else {
		perms, _ := sys_repositories.MenuCrud.SelectMenuPermsByUserId(userId)
		permsSet = SplitPerms(perms)
	}
	return permsSet
}

// IsFace 是否外链
func (m *MenuService) IsFace(menu sys_model.SysMenu) (status bool) {
	if menu.IsFrame == "0" && ternary(menu.Path) != "" {
		return true
	} else if menu.IsFrame == "1" {
		return true
	}
	return false
}

// SplitPerms 工具函数：将权限字符串切片转换为去重后的权限集合
func SplitPerms(perms []string) []string {
	// 使用 map 来去重
	permMap := make(map[string]struct{})

	// 遍历权限
	for _, perm := range perms {
		// 判断非空
		if strings.TrimSpace(perm) != "" {
			// 分割权限字符串
			splitPerms := strings.Split(perm, ",")
			// 遍历分割后的权限
			for _, p := range splitPerms {
				// 去除首尾空格并添加到 map 中
				trimmed := strings.TrimSpace(p)
				if trimmed != "" {
					permMap[trimmed] = struct{}{}
				}
			}
		}
	}

	// 将 map 转换为切片
	result := make([]string, 0, len(permMap))
	for p := range permMap {
		result = append(result, p)
	}

	return result
}

// GetMenuTreeAll 定义菜单树结构 一次性查询所有数据，然后组装树（推荐，性能更好
func (m *MenuService) GetMenuTreeAll() ([]*sys_model.MenuTree, error) {
	menus, err := sys_repositories.MenuCrud.GetMenuTreeAll()
	if err != nil {
		return nil, err
	}
	// 转换为map，方便查找
	menuMap := make(map[string]*sys_model.MenuTree)
	var trees []*sys_model.MenuTree
	// 先转换为MenuTree对象
	for _, menu := range menus {
		menuMap[menu.MenuId] = &sys_model.MenuTree{
			SysMenu:  menu,
			Children: make([]*sys_model.MenuTree, 0),
		}
	}
	// 组装树
	for _, menu := range menus {
		tree := menuMap[menu.MenuId]
		if menu.ParentId == "0" {
			// 根节点
			trees = append(trees, tree)
		} else {
			// 非根节点
			if parent, ok := menuMap[menu.ParentId]; ok {
				parent.Children = append(parent.Children, tree)
			}
		}
	}

	return trees, nil
}

// BuildMenus  拼接路由
func (m *MenuService) BuildMenus(menus []*sys_model.MenuTree) []*sys_model.RouterVo {
	routers := make([]*sys_model.RouterVo, 0)
	for _, menu := range menus {
		router := &sys_model.RouterVo{}
		router.Hidden = menu.Visible == "1"

		// 设置路由基本信息
		router.Name = getRouteName(menu)
		router.MenuType = menu.MenuType
		router.Path = getRouterPath(menu)
		router.Query = menu.Query
		router.Component = getComponent(menu)

		// 设置 Meta 信息
		router.Meta = sys_model.MetaVo{
			Title:   menu.MenuName,
			Icon:    menu.Icon,
			NoCache: menu.IsCache == "1",
			Link:    ternary(menu.Path),
		}
		//if StringUtils.IsNotEmpty(menu.Children) && consts.TYPE_DIR == menu.MenuType {
		//	router.AlwaysShow = true
		//	//router.setChildren(buildMenus(cMenus))
		//}
		if menu.IsFrame == "1" {
			router.Redirect = "noRedirect"
		}

		// 处理子菜单
		if len(menu.Children) > 0 {
			if menu.MenuType == "M" && len(menu.Children) > 0 {
				router.AlwaysShow = true
			} else {
				router.AlwaysShow = false
			}

			// 递归处理子菜单
			children := m.BuildMenus(menu.Children)
			if len(children) > 0 {
				router.Children = children // 直接赋值所有子菜单
			}
		}

		routers = append(routers, router)
	}
	return routers
}

// BuildMenuTree 构建菜单树 - 使用 MenuTree 结构
func (m *MenuService) BuildMenuTree(menus []*sys_model.SysMenu) []*sys_model.TreeSelect {
	// 创建 map 用于快速查找
	menuMap := make(map[string]*sys_model.TreeSelect)
	var rootNodes []*sys_model.TreeSelect

	// 第一次遍历：创建所有节点
	for _, menu := range menus {
		node := &sys_model.TreeSelect{
			Id:       menu.MenuId,
			Label:    menu.MenuName,
			Disabled: menu.Status == "1",
			Children: make([]*sys_model.TreeSelect, 0),
		}
		menuMap[menu.MenuId] = node

		// 如果是根节点，加入 rootNodes
		if menu.ParentId == "0" {
			rootNodes = append(rootNodes, node)
			continue
		}

		// 将当前节点添加到父节点的 children 中
		if parent, exists := menuMap[menu.ParentId]; exists {
			parent.Children = append(parent.Children, node)
		}
	}

	return rootNodes
}

// ternary 使用辅助函数
func ternary(Path string) string {
	if StringUtils.IsHttp(Path) {
		return Path
	}
	return ""
}

// getComponent 判断组件类型
func getComponent(menu *sys_model.MenuTree) string {
	component := consts.LAYOUT
	if StringUtils.IsBlank(menu.Component) {
		component = menu.Component
	} else if StringUtils.IsEmptyStr(menu.Component) && menu.ParentId != "0" && isInnerLink(menu) {
		component = consts.INNER_LINK
	} else if StringUtils.IsEmptyStr(menu.Component) && isParentView(menu) {
		component = consts.PARENT_VIEW
	}

	return component
}

// isParentView 是否为parent_view组件
func isParentView(menu *sys_model.MenuTree) bool {
	return menu.ParentId != "0" && consts.TYPE_DIR == menu.MenuType
}

// getRouteName 获取路由名称
func getRouteName(menu *sys_model.MenuTree) string {
	routerName := StringUtils.Capitalize(menu.Path)
	// 非外链并且是一级目录（类型为目录）
	if isMenuFrame(menu) {
		routerName = consts.EMPTY
	}
	return routerName
}

func getRouterPath(menu *sys_model.MenuTree) string {
	routerPath := menu.Path
	// 内链打开外网方式
	if menu.ParentId != "0" && isInnerLink(menu) {
		routerPath = innerLinkReplaceEach(routerPath)
	}
	// 非外链并且是一级目录（类型为目录）
	if "0" == menu.ParentId && consts.TYPE_DIR == menu.MenuType && consts.NO_FRAME == menu.IsFrame {
		routerPath = "/" + menu.Path
	} else if isMenuFrame(menu) {
		// 非外链并且是一级目录（类型为菜单）
		routerPath = "/"
	}
	return routerPath
}

// isMenuFrame 是否为菜单内部跳转
func isMenuFrame(menu *sys_model.MenuTree) bool {
	return menu.ParentId == "0" && consts.TYPE_MENU == menu.MenuType && menu.IsFrame == consts.NO_FRAME
}

// isInnerLink 是否为内链组件
func isInnerLink(menu *sys_model.MenuTree) bool {
	return menu.IsFrame == consts.NO_FRAME && StringUtils.IsHttp(menu.Path)
}

// isInnerLink 内链域名特殊字符替换 return 替换后的内链域名
func innerLinkReplaceEach(path string) string {
	return strings.NewReplacer(
		"http://", "",
		"https://", "",
		"www.", "",
		".", "/",
		":", "/",
	).Replace(path)
}

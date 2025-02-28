package sys_service

import (
	"hrkGo/app/model/sys_model"
	"hrkGo/app/repositories/sys_repositories"
)

type RoleService struct {
}

// GetRoleList 获取角色列表
func (u *RoleService) GetRoleList(req sys_model.RoleListRequest) (list []sys_model.SysRole, total int64, err error) {
	list, total, err = sys_repositories.RoleCrud.GetRoleList(req)
	return list, total, err
}

// SelectRoleDataById 获取角色id查找角色数据
func (u *RoleService) SelectRoleDataById(roleId string) (role sys_model.SysRole, err error) {
	role, err = sys_repositories.RoleCrud.SelectRoleDataById(roleId)
	return role, err
}

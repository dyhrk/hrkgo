package sys_repositories

import (
	"fmt"
	"hrkGo/app/model/sys_model"
	"hrkGo/utils/global/variable"
	"time"
)

type roleCrud struct{}

var RoleCrud = new(roleCrud)

// GetRoleList 实现接口方法
func (m *roleCrud) GetRoleList(req sys_model.RoleListRequest) (list []sys_model.SysRole, total int64, err error) {
	// 构建查询条件
	db := variable.GormDbMysql.Model(&sys_model.SysRole{})

	// 名称模糊查询
	if req.RoleName != "" {
		db = db.Where("role_name LIKE ?", "%"+req.RoleName+"%")
	}
	if req.RoleKey != "" {
		db = db.Where("role_key LIKE ?", "%"+req.RoleKey+"%")
	}
	// 状态查询
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	// 添加时间范围查询
	if req.BeginTime != "" && req.EndTime != "" {

		start, _ := time.Parse("2006-01-02", req.BeginTime)

		end, _ := time.Parse("2006-01-02", req.EndTime)

		// 将结束日期调整到当天的最后一刻
		end = end.Add(24 * time.Hour).Add(-time.Second)
		fmt.Println(start, end, "end")

		db = db.Where("create_time BETWEEN ? AND ?",
			start,
			end)
	}
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 查询列表
	offset := (req.PageNum - 1) * req.PageSize

	err = db.
		Offset(offset).
		Limit(req.PageSize).
		Find(&list).Error
	return list, total, err
}

// SelectRoleDataById 通过角色ID查询角色
func (m *roleCrud) SelectRoleDataById(roleId string) (role sys_model.SysRole, err error) {
	err = variable.GormDbMysql.Where("role_id = ?", roleId).First(&role).Error
	fmt.Println(role)
	return role, err
}

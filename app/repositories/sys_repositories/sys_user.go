package sys_repositories

import (
	"errors"
	"fmt"
	"hrkGo/app/model/sys_model"
	"hrkGo/utils/global/consts"
	"hrkGo/utils/global/variable"
	"time"
)

type userCrud struct{}

var UserCrud = new(userCrud)

// SelectUserByUserName 根据username查询用户
func (u *userCrud) SelectUserByUserName(username string) (user *sys_model.SysUser, err error) {
	result := variable.GormDbMysql.Where("user_name = ? and del_flag = 0", username).First(&user)
	if result.RowsAffected == 0 {
		return nil, errors.New(consts.EXISTS)
	}
	return
}

// SelectUserList 获取用户列表
func (u *userCrud) SelectUserList(query sys_model.UserListRequest) ([]sys_model.SysUserInfo, int64, error) {

	db := variable.GormDbMysql.Model(&sys_model.SysUserInfo{}).
		Select("sys_user.*, sys_dept.dept_name, sys_dept.leader").
		Joins("LEFT JOIN sys_dept ON sys_user.dept_id = sys_dept.dept_id").
		Where("sys_user.del_flag = ?", "0")

	db = db.Preload("Dept")

	if query.UserName != "" {
		db = db.Where("sys_user.user_name LIKE ?", "%"+query.UserName+"%")
	}

	if query.Status != "" {
		db = db.Where("sys_user.status = ?", query.Status)
	}

	if query.Phonenumber != "" {
		db = db.Where("sys_user.phonenumber LIKE ?", "%"+query.Phonenumber+"%")
	}

	// 添加时间范围查询
	if query.BeginTime != "" && query.EndTime != "" {
		start, _ := time.Parse("2006-01-02", query.BeginTime)
		end, _ := time.Parse("2006-01-02", query.EndTime)
		// 将结束日期调整到当天的最后一刻
		end = end.Add(24 * time.Hour).Add(-time.Second)
		fmt.Println(start, end, "end")

		db = db.Where("sys_user.create_time BETWEEN ? AND ?",
			start,
			end)
	}

	fmt.Println("111111111111")

	// 部门ID查询
	if query.DeptId != "" {
		db = db.Where("(sys_user.dept_id = ? OR sys_user.dept_id IN (SELECT dept_id FROM sys_dept WHERE FIND_IN_SET(?, ancestors)))",
			query.DeptId, query.DeptId)
	}
	// 先获取总数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	var users []sys_model.SysUserInfo
	err = db.Find(&users).Error
	return users, total, err

}

// InsertUser 插入用户方法
func (u *userCrud) InsertUser(user sys_model.SysUser) error {
	// 创建一个map来存储非空字段
	fields := make(map[string]interface{})

	// 只添加非空字段
	if user.UserId != "0" {
		fields["user_id"] = user.UserId
	}
	if user.DeptId != "0" {
		fields["dept_id"] = user.DeptId
	}
	if user.UserName != "" {
		fields["user_name"] = user.UserName
	}
	if user.NickName != "" {
		fields["nick_name"] = user.NickName
	}
	if user.Email != "" {
		fields["email"] = user.Email
	}
	if user.Avatar != "" {
		fields["avatar"] = user.Avatar
	}
	if user.Phonenumber != "" {
		fields["phonenumber"] = user.Phonenumber
	}
	if user.Sex != "" {
		fields["sex"] = user.Sex
	}
	if user.Password != "" {
		fields["password"] = user.Password
	}
	if user.Status != "" {
		fields["status"] = user.Status
	}
	if user.CreateBy != "" {
		fields["create_by"] = user.CreateBy
	}
	if user.Remark != "" {
		fields["remark"] = user.Remark
	}

	// 添加创建时间
	fields["create_time"] = time.Now()

	// 执行插入
	result := variable.GormDbMysql.Model(&sys_model.SysUser{}).Create(fields)
	return result.Error
}

package c_system

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"hrkGo/app/middleware"
	"hrkGo/app/model/sys_model"
	"hrkGo/app/service/sys_service"
	"hrkGo/utils/response"
)

type UserController struct {
	UserService sys_service.UserService
	DeptService sys_service.DeptService
}

// DeptTree 部门列表
func (d UserController) DeptTree(c *gin.Context) {
	var req sys_model.DeptListRequest
	depts, err := d.DeptService.SelectDeptList(req, middleware.DataScopeMiddleware(c.MustGet("User"), "d", "u"))
	if err != nil {
		response.BusinessFail(c, "查询失败，请稍后再试")
	}
	TreeSelect := d.DeptService.BuildDeptTree(depts)
	response.Success(c, "查询成功", TreeSelect)

}

// SelectUserList 用户列表
func (d UserController) SelectUserList(c *gin.Context) {
	var req sys_model.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidateFail(c, response.GetErrorMsg(req, err))
		return
	}

	fmt.Println(req)

	list, total, err := d.UserService.SelectUserList(req)
	if err != nil {

	}
	response.SuccessRow(c, "查询成功", list, total)

}

// InsertUser 新增用户
func (d UserController) InsertUser(c *gin.Context) {
	var req sys_model.SysUser
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidateFail(c, response.GetErrorMsg(req, err))
		return
	}

	err := d.UserService.InsertUser(req)
	if err != nil {

	}
	response.SuccessNil(c, "操作成功")

}

// ResetPwd 修改密码
func (d UserController) ResetPwd(c *gin.Context) {
	var req sys_model.SysUser
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidateFail(c, response.GetErrorMsg(req, err))
		return
	}

	if sys_model.IsAdmin(req.UserId) {
		response.BusinessFail(c, "不允许操作超级管理员用户")
		return
	}

	err := d.UserService.ResetPwd(req)
	if err != nil {

	}
	response.SuccessNil(c, "操作成功")

}

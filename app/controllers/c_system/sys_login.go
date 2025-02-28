package c_system

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"hrkGo/app/model/sys_model"
	"hrkGo/app/service/sys_service"
	"hrkGo/utils/global/consts"
	"hrkGo/utils/global/variable"
	"hrkGo/utils/response"
	"net/http"
)

type Controller struct {
	loginService sys_service.LoginService
	CaptService  sys_service.CaptService
	JwtCurd      sys_service.JwtService
	MenuService  sys_service.MenuService
}

// 1. 直接创建并赋值
func CreateTokenData(user *sys_model.SysUserInfo, perms []string) sys_model.TokenData {
	return sys_model.TokenData{
		User:  user,
		Perms: perms,
	}
}

// WxLoginRequest 接收小程序登录请求的参数
type WxLoginRequest struct {
	Code string `json:"code" binding:"required"`
}

// CaptchaImage 获取验证码
func (l *Controller) CaptchaImage(c *gin.Context) {
	id, b64s, err := l.CaptService.CaptMake()
	if err != nil {
		response.ValidateFail(c, "验证码接口异常")
	}
	response.Success(c, consts.CurdStatusOkMsg, gin.H{
		"uuid": id, "img": b64s,
	})
}

// Login 管理端登录
func (l *Controller) Login(c *gin.Context) {
	var form sys_model.Login
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, response.GetErrorMsg(form, err))
		return
	}
	if err := l.CaptService.CaptVerify(form.Uuid, form.Code); err != true {
		response.ValidateFail(c, consts.CAPTCGAERROR)
		return
	}
	if user, err := l.loginService.Login(form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		perms := l.MenuService.LoginSelectMenuPermsByUserId(user.UserId)
		UserData, _ := l.loginService.GetUserInfo(user.UserId)
		tokenData := CreateTokenData(UserData, perms)
		token, err := l.JwtCurd.GenerateTokenWithCustomClaims(tokenData)
		if err != nil {
			response.BusinessFail(c, err.Error())
			return
		}
		response.SuccessToken(c, consts.SUCCESS, token)
	}
}

// Info 个人信息
func (l *Controller) Info(c *gin.Context) {
	permissions := sys_service.GetUserData(c.Keys["tokenId"].(string))

	// 创建一个新的数组来存储 roleKey
	var roles []string

	// 遍历 user.roles 数组
	for _, role := range permissions.User.Roles {
		roles = append(roles, role.RoleKey)
	}

	response.SuccessInfo(c, "操作成功", permissions.Perms, roles, permissions.User)
}

// WxLogin 微信登录
func (l *Controller) WxLogin(c *gin.Context) {
	var loginReq WxLoginRequest
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求参数",
		})
		return
	}
	// 请求微信服务器获取 session_key 和 openid
	wxLoginURL := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		variable.ConfigYml.GetString("wechat.MiniProgram.appid"), variable.ConfigYml.GetString("wechat.MiniProgram.secret"), loginReq.Code)
	resp, err := http.Get(wxLoginURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "请求微信服务器失败",
		})
		return
	}
	defer resp.Body.Close()
	// 解析微信返回的数据
	var wxResp struct {
		SessionKey string `json:"session_key"`
		OpenID     string `json:"openid"`
		UnionID    string `json:"unionid"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wxResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "解析微信响应失败",
		})
		return
	}
	if wxResp.ErrCode != 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": wxResp.ErrMsg,
		})
		return
	}
	if user, err := l.loginService.GetOpenIdUser(wxResp.OpenID); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		if user.UserId == "0" {
			response.Success(c, "操作成功", nil)
			return
		}
		// token, err := l.JwtCurd.GenerateTokenWithCustomClaims(user)
		// if err != nil {
		//	response.BusinessFail(c, err.Error())
		//	return
		//}
		// response.SuccessToken(c, "操作成功", token)
	}

}

// Logout 管理端登录
func (l *Controller) Logout(c *gin.Context) {
	response.SuccessNil(c, consts.SUCCESS)
}

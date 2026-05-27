package frontend

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type RegisterReq struct {
	g.Meta       `path:"/user/register" method:"post" tags:"前台用户" summary:"用户注册"`
	Name         string `json:"name"         description:"用户名" v:"required#用户名必填"`
	Password     string `json:"password"     description:"密码" v:"password"`
	Avatar       string `json:"avatar"       description:"头像"`
	UserSalt     string `json:"userSalt"     description:"加密盐 生成密码用"`
	Sex          int    `json:"sex"          description:"1男 2女"`
	Status       int    `json:"status"       description:"1正常 2拉黑冻结"`
	Sign         string `json:"sign"         description:"个性签名"`
	SecretAnswer string `json:"secretAnswer" description:"密保问题的答案"`
}

type RegisterRes struct {
	Id uint `json:"id"`
}

// for jwt
type LoginReq struct {
	g.Meta   `path:"/login" method:"post" tags:"前台用户" summary:"前台用户登录"`
	Name     string `json:"name"         description:"用户名" v:"required#用户名必填"`
	Password string `json:"password"     description:"密码" v:"required#密码必填"`
}

type LoginRes struct {
	Token  string    `json:"token"`
	Expire time.Time `json:"expire"`
}

type UserInfoReq struct {
	g.Meta `path:"/user/info" method:"get" tags:"前台用户" summary:"当前登录用户信息"`
}

type UserInfoRes struct {
	UserInfoBase
}

// 可以复用的，一定要抽取出来
type UserInfoBase struct {
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Sex    uint8  `json:"sex"`
	Sign   string `json:"sign"`
	Status uint8  `json:"status"`
}

// 修改密码
type UpdatePasswordReq struct {
	g.Meta       `path:"/update/password" method:"post" tag:"前台用户" summary:"修改密码"`
	Password     string `json:"password"  v:"password"   description:""`
	UserSalt     string `json:"userSalt,omitempty"     description:"加密盐 生成密码用"`
	SecretAnswer string `json:"secretAnswer" description:"密保问题的答案"`
}

type UpdatePasswordRes struct {
	Id uint `json:"id"`
}

// 登出
type LogoutReq struct {
	g.Meta `path:"/user/logout" method:"post" tags:"前台用户" summary:"用户登出"`
}

type LogoutRes struct {
	Message string `json:"message"`
}

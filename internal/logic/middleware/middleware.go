package middleware

import (
	"strings"

	"bit303_shop/internal/consts"
	"bit303_shop/internal/model"
	"bit303_shop/internal/service"
	"bit303_shop/utility/response"

	"github.com/goflyfox/gtoken/gtoken"
	jwt "github.com/gogf/gf-jwt/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

// 超管专属路径（权限管理模块）
var adminOnlyPaths = []string{
	"/backend/role",
	"/backend/permission",
	"/backend/admin",
	"/backend/user", // 用户管理
}

// 不需要权限校验的路径（登录、注册、登出等）
var noPermissionCheckPaths = []string{
	"/backend/login",
	"/backend/logout",
	"/backend/admin/create",
	"/backend/refresh-token",
}

type sMiddleware struct {
	LoginUrl string // 登录路由地址
}

func init() {
	service.RegisterMiddleware(New())
}

func New() *sMiddleware {
	return &sMiddleware{
		LoginUrl: "/backend/login",
	}
}

const (
	CtxAccountId      = "account_id"       //token获取
	CtxAccountName    = "account_name"     //token获取
	CtxAccountAvatar  = "account_avatar"   //token获取
	CtxAccountSex     = "account_sex"      //token获取
	CtxAccountStatus  = "account_status"   //token获取
	CtxAccountSign    = "account_sign"     //token获取
	CtxAccountIsAdmin = "account_is_admin" //token获取
	CtxAccountRoleIds = "account_role_ids" //token获取
)

type TokenInfo struct {
	Id   int
	Name string
	//Avatar  string
	//Sex     int
	//Status  int
	//Sign    string
	//RoleIds string
	//IsAdmin int
}

// 返回处理中间件
func (s *sMiddleware) ResponseHandler(r *ghttp.Request) {
	r.Middleware.Next()

	// 如果已经有返回内容，那么该中间件什么也不做
	if r.Response.BufferLength() > 0 {
		return
	}

	var (
		err             = r.GetError()
		res             = r.GetHandlerResponse()
		code gcode.Code = gcode.CodeOK
	)
	g.Log().Info(r.Context(), "ResponseHandler res =", res, "err =", err) // 调试
	if err != nil {
		code = gerror.Code(err)
		if code == gcode.CodeNil {
			code = gcode.CodeInternalError
		}
		response.JsonExit(r, code.Code(), err.Error())
		//if r.IsAjaxRequest() {
		//	response.JsonExit(r, code.Code(), err.Error())
		//} else {
		//	service.View().Render500(r.Context(), model.View{
		//		Error: err.Error(),
		//	})
		//}
	} else {
		response.JsonExit(r, code.Code(), "", res)
		//if r.IsAjaxRequest() {
		//	response.JsonExit(r, code.Code(), "", res)
		//} else {
		//	// 什么都不做，业务API自行处理模板渲染的成功逻辑。
		//}
	}
}

// 自定义上下文对象
func (s *sMiddleware) Ctx(r *ghttp.Request) {
	// 初始化，务必最开始执行
	customCtx := &model.Context{
		Session: r.Session,
		Data:    make(g.Map),
	}
	service.BizCtx().Init(r, customCtx)
	if userEntity := service.Session().GetUser(r.Context()); userEntity.Id > 0 {
		customCtx.User = &model.ContextUser{
			Id:   uint(userEntity.Id),
			Name: userEntity.Name,
			//Nickname: userEntity.Nickname,
			//Avatar:   userEntity.Avatar,
			IsAdmin: uint8(userEntity.IsAdmin),
		}
	}
	// 将自定义的上下文对象传递到模板变量中使用
	r.Assigns(g.Map{
		"Context": customCtx,
	})
	// 执行下一步请求逻辑
	r.Middleware.Next()
}

func (s *sMiddleware) CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func (s *sMiddleware) Auth(r *ghttp.Request) {
	service.Auth().MiddlewareFunc()(r)
	r.Middleware.Next()
}

// UserAuth ensures the frontend JWT middleware runs and then continues the handler chain.
func (s *sMiddleware) UserAuth(r *ghttp.Request) {
	// 先运行JWT中间件进行基础验证
	service.UserAuth().MiddlewareFunc()(r)
	if r.IsExited() {
		return
	}

	// 从JWT中获取用户ID并存储到上下文中
	userID := service.UserAuth().GetIdentity(r.Context())
	g.Log().Info(r.Context(), "UserAuth 当前用户ID:", userID)
	r.SetCtxVar(consts.CtxUserId, userID)

	r.Middleware.Next()
}

var GToken *gtoken.GfToken

// PermissionCheck 权限校验中间件
// 检查当前管理员是否有权限访问当前路径
func (s *sMiddleware) PermissionCheck(r *ghttp.Request) {
	ctx := r.Context()
	requestPath := r.URL.Path

	// 检查是否是不需要权限校验的路径
	for _, path := range noPermissionCheckPaths {
		if strings.HasPrefix(requestPath, path) {
			r.Middleware.Next()
			return
		}
	}

	// 从JWT中获取管理员信息
	claims := jwt.ExtractClaims(ctx)
	if claims == nil {
		response.JsonExit(r, 401, "未登录或登录已过期")
		return
	}

	// 获取 is_admin 和 role_ids
	isAdmin := gconv.Int(claims["is_admin"])
	roleIdsStr := gconv.String(claims["role_ids"])

	g.Log().Debug(ctx, "PermissionCheck - path:", requestPath, "is_admin:", isAdmin, "role_ids:", roleIdsStr)

	// 超级管理员拥有所有权限
	if isAdmin == 1 {
		r.Middleware.Next()
		return
	}

	// 检查是否是超管专属路径
	for _, adminPath := range adminOnlyPaths {
		if strings.HasPrefix(requestPath, adminPath) {
			response.JsonExit(r, 403, "权限不足，该功能仅超级管理员可访问")
			return
		}
	}

	// 解析 role_ids（格式如 "1,2,3"）
	var roleIds []int
	if roleIdsStr != "" {
		roleIdStrs := strings.Split(roleIdsStr, ",")
		for _, idStr := range roleIdStrs {
			idStr = strings.TrimSpace(idStr)
			if idStr != "" {
				roleIds = append(roleIds, gconv.Int(idStr))
			}
		}
	}

	// 如果没有角色，则无权限
	if len(roleIds) == 0 {
		response.JsonExit(r, 403, "权限不足，未分配角色")
		return
	}

	// 获取该角色拥有的所有权限路径
	allowedPaths, err := service.Permission().GetPathsByRoleIds(ctx, roleIds)
	if err != nil {
		g.Log().Error(ctx, "获取权限路径失败:", err)
		response.JsonExit(r, 500, "权限校验失败")
		return
	}

	g.Log().Debug(ctx, "PermissionCheck - allowedPaths:", allowedPaths)

	// 检查当前请求路径是否匹配任一权限路径（前缀匹配）
	for _, allowedPath := range allowedPaths {
		if strings.HasPrefix(requestPath, allowedPath) {
			r.Middleware.Next()
			return
		}
	}

	// 没有匹配的权限
	response.JsonExit(r, 403, "权限不足，无法访问该功能")
}

// Gtoken鉴权
func (s *sMiddleware) GTokenSetCtx(r *ghttp.Request) {
	var tokenInfo TokenInfo
	token := GToken.GetTokenData(r)
	err := gconv.Struct(token.GetString("data"), &tokenInfo)
	if err != nil {
		response.Auth(r)
		return
	}
	//账号被冻结拉黑
	//if tokenInfo.Status == 2 {
	//	response.AuthBlack(r)
	//	return
	//}
	r.SetCtxVar(CtxAccountId, tokenInfo.Id)
	r.SetCtxVar(CtxAccountName, tokenInfo.Name)
	//r.SetCtxVar(CtxAccountAvatar, tokenInfo.Avatar)
	//r.SetCtxVar(CtxAccountSex, tokenInfo.Sex)
	//r.SetCtxVar(CtxAccountStatus, tokenInfo.Status)
	//r.SetCtxVar(CtxAccountSign, tokenInfo.Sign)
	//r.SetCtxVar(CtxAccountRoleIds, tokenInfo.RoleIds)
	//r.SetCtxVar(CtxAccountIsAdmin, tokenInfo.Sign)
	r.Middleware.Next()
}

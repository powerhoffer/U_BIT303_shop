package service

import (
	"context"
	"errors"
	"time"

	"bit303_shop/internal/dao"
	"bit303_shop/internal/model/entity"
	"bit303_shop/utility"

	jwt "github.com/gogf/gf-jwt/v2"
	"github.com/gogf/gf/v2/frame/g"
)

var userAuthService *sUserAuth

func UserAuth() *jwt.GfJWTMiddleware {
	return userAuthService.jwt
}

// GetUserAuthService 获取完整的用户认证服务
func GetUserAuthService() *sUserAuth {
	return userAuthService
}

type sUserAuth struct {
	// 自定义 JWT 认证
	jwt *jwt.GfJWTMiddleware

	// 自定义数据
	AuthName    string
	Key         []byte
	MaxAge      int
	IdentityKey string
	LoginPath   string
	LogoutPath  string
	RefreshPath string
}

var (
	// 默认的AuthName
	defaultAuthName = "jwt_frontend"
	// 默认的key
	defaultKey = []byte("secret key user")
	// 默认的maxAge
	defaultMaxAge = 60 * 60 * 24
	// 默认的IdentityKey
	defaultIdentityKey = "user_id"
	// 默认的LoginPath
	defaultLoginPath = "/login"
	// 默认的LogoutPath
	defaultLogoutPath = "/user/logout"
	// 默认的RefreshPath
	defaultRefreshPath = "/refresh"
)

// New 创建
func New() *sUserAuth {
	s := &sUserAuth{
		// 默认配置
		AuthName:    defaultAuthName,
		Key:         defaultKey,
		MaxAge:      defaultMaxAge,
		IdentityKey: defaultIdentityKey,
		LoginPath:   defaultLoginPath,
		LogoutPath:  defaultLogoutPath,
		RefreshPath: defaultRefreshPath,
	}

	// 创建 JWT 中间件实例
	s.jwt = jwt.New(&jwt.GfJWTMiddleware{
		Realm:           "shop-frontend",
		Key:             s.Key,
		Timeout:         time.Duration(s.MaxAge) * time.Second,
		MaxRefresh:      time.Duration(s.MaxAge) * time.Second,
		IdentityKey:     s.IdentityKey,
		TokenLookup:     "header: Authorization, query: token, cookie: jwt_frontend",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
		Authenticator:   UserAuthenticator,
		Unauthorized:    UserUnauthorized,
		PayloadFunc:     UserPayloadFunc,
		IdentityHandler: UserIdentityHandler,
	})

	return s
}

// GetUserAuth 获取单例对象
func GetUserAuth() *sUserAuth {
	return userAuthService
}

// Init 初始化
func (s *sUserAuth) Init() {
	// 这里不需要重新配置，因为JWT实例已经在New()方法中创建和配置好了
}

func init() {
	userAuthService = New()
}

// UserPayloadFunc 生成前台用户 JWT 的 payload
func UserPayloadFunc(data interface{}) jwt.MapClaims {
	claims := jwt.MapClaims{}
	if user, ok := data.(*entity.UserInfo); ok {
		claims[userAuthService.IdentityKey] = user.Id
		claims["name"] = user.Name
		claims["status"] = user.Status
	}
	return claims
}

// UserIdentityHandler 从 JWT 中提取用户 ID
func UserIdentityHandler(ctx context.Context) interface{} {
	claims := jwt.ExtractClaims(ctx)
	g.Log().Info(ctx, "UserIdentityHandler 当前用户ID:", claims[userAuthService.IdentityKey])
	return claims[userAuthService.IdentityKey]
}

// UserUnauthorized 自定义的Unauthorized函数，用于调试JWT验证失败原因
func UserUnauthorized(ctx context.Context, code int, message string) {
	r := g.RequestFromCtx(ctx)
	// 记录JWT验证失败的详细信息
	g.Log().Info(ctx, "UserUnauthorized - Code:", code, "Message:", message)
	g.Log().Info(ctx, "Request Header:", r.Header)
	// 获取所有cookie
	cookies := make(map[string]string)
	for _, cookie := range r.Cookies() {
		cookies[cookie.Name] = cookie.Value
	}
	g.Log().Info(ctx, "Request Cookie:", cookies)

	r.Response.WriteJson(g.Map{
		"code":    code,
		"message": message,
	})
	r.ExitAll()
}

// UserAuthenticator 校验前台用户登录
func UserAuthenticator(ctx context.Context) (interface{}, error) {
	r := g.RequestFromCtx(ctx)
	name := r.Get("name").String()
	password := r.Get("password").String()

	if name == "" || password == "" {
		return nil, jwt.ErrMissingLoginValues
	}

	var user entity.UserInfo
	if err := dao.UserInfo.Ctx(ctx).Where(dao.UserInfo.Columns().Name, name).Scan(&user); err != nil {
		return nil, jwt.ErrFailedAuthentication
	}
	// 使用 EncryptPassword 校验密码
	if utility.EncryptPassword(password, user.UserSalt) != user.Password {
		return nil, jwt.ErrFailedAuthentication
	}

	// 检查用户状态：status=0 表示已冻结
	if user.Status == 0 {
		return nil, ErrUserFrozen
	}

	return &user, nil
}

// ErrUserFrozen 用户已被冻结错误
var ErrUserFrozen = errors.New("账号已被冻结，无法登录")

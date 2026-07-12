package middleware

import (
	"database/sql"
	"errors"
	"strings"

	"bit303_shop/internal/consts"
	"bit303_shop/internal/dao"
	"bit303_shop/internal/model/entity"
	"bit303_shop/internal/service"
	"bit303_shop/utility"
	"bit303_shop/utility/response"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type sMiddleware struct{}

func init() {
	service.RegisterMiddleware(New())
}

func New() *sMiddleware {
	return &sMiddleware{}
}

func (s *sMiddleware) CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}

func (s *sMiddleware) ResponseHandler(r *ghttp.Request) {
	r.Middleware.Next()

	if r.Response.BufferLength() > 0 {
		return
	}

	err := r.GetError()
	res := r.GetHandlerResponse()
	var code gcode.Code = gcode.CodeOK
	if err != nil {
		code = gerror.Code(err)
		if code == gcode.CodeNil {
			code = gcode.CodeInternalError
		}
		message := err.Error()
		if shouldMaskErrorMessage(message) {
			g.Log().Error(r.Context(), err)
			message = "Operation failed"
		}
		response.JsonExit(r, code.Code(), message)
		return
	}
	response.JsonExit(r, code.Code(), "", res)
}

func (s *sMiddleware) EmployeeAuth(r *ghttp.Request) {
	token := bearerToken(r)
	if token == "" {
		response.JsonExit(r, 401, "Unauthorized or login expired")
		return
	}
	claims, err := utility.ParseEmployeeToken(token)
	if err != nil || claims.EmployeeId == 0 {
		response.JsonExit(r, 401, "Unauthorized or login expired")
		return
	}

	columns := dao.EmployeeInfo.Columns()
	var employee entity.EmployeeInfo
	err = dao.EmployeeInfo.Ctx(r.Context()).
		Where(columns.Id, claims.EmployeeId).
		WhereNull(columns.DeletedAt).
		Scan(&employee)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		response.JsonExit(r, 500, "Failed to query employee info")
		return
	}
	if employee.Id == 0 {
		response.JsonExit(r, 401, "Employee account does not exist")
		return
	}
	if employee.Status != consts.EmployeeStatusNormal {
		response.JsonExit(r, 403, "Employee account is disabled")
		return
	}

	r.SetCtxVar(consts.CtxEmployeeId, employee.Id)
	r.SetCtxVar(consts.CtxEmployeeUsername, employee.Username)
	r.Middleware.Next()
}

func (s *sMiddleware) AdminAuth(r *ghttp.Request) {
	token := bearerToken(r)
	if token == "" {
		response.JsonExit(r, 401, "Unauthorized or login expired")
		return
	}
	claims, err := utility.ParseAdminToken(token)
	if err != nil || claims.AdminId == 0 {
		response.JsonExit(r, 401, "Unauthorized or login expired")
		return
	}

	columns := dao.AdminInfo.Columns()
	var admin entity.AdminInfo
	err = dao.AdminInfo.Ctx(r.Context()).
		Where(columns.Id, claims.AdminId).
		WhereNull(columns.DeletedAt).
		Scan(&admin)
	if err != nil {
		response.JsonExit(r, 500, "Failed to query admin info")
		return
	}
	if admin.Id == 0 {
		response.JsonExit(r, 401, "Admin account does not exist")
		return
	}
	if admin.Status != consts.AdminStatusNormal {
		response.JsonExit(r, 403, "Admin account is disabled")
		return
	}

	r.SetCtxVar(consts.CtxAdminId, admin.Id)
	r.SetCtxVar(consts.CtxAdminUsername, admin.Username)
	r.SetCtxVar(consts.CtxAdminIsSuper, admin.IsSuper)
	r.Middleware.Next()
}

func (s *sMiddleware) AdminPermissionAuth(r *ghttp.Request) {
	if r.GetCtxVar(consts.CtxAdminIsSuper).Int() == 1 {
		r.Middleware.Next()
		return
	}
	adminId := r.GetCtxVar(consts.CtxAdminId).Uint()
	if adminId == 0 {
		response.JsonExit(r, 401, "Unauthorized or login expired")
		return
	}
	if !hasAdminPermission(r, adminId, r.Method, r.URL.Path) {
		response.JsonExit(r, 403, "Permission denied")
		return
	}
	r.Middleware.Next()
}

func bearerToken(r *ghttp.Request) string {
	token := strings.TrimSpace(r.Header.Get("Authorization"))
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))
	}
	return token
}

func hasAdminPermission(r *ghttp.Request, adminId uint, method string, path string) bool {
	var roleRelations []entity.AdminRoleRelation
	if err := dao.AdminRoleRelation.Ctx(r.Context()).Where(dao.AdminRoleRelation.Columns().AdminId, adminId).Scan(&roleRelations); err != nil {
		return false
	}
	roleIds := make([]uint, 0, len(roleRelations))
	for _, relation := range roleRelations {
		roleIds = append(roleIds, relation.RoleId)
	}
	if len(roleIds) == 0 {
		return false
	}

	var roles []entity.AdminRole
	roleColumns := dao.AdminRole.Columns()
	if err := dao.AdminRole.Ctx(r.Context()).
		WhereIn(roleColumns.Id, roleIds).
		Where(roleColumns.Status, consts.AdminRoleEnabled).
		WhereNull(roleColumns.DeletedAt).
		Scan(&roles); err != nil {
		return false
	}
	enabledRoleIds := make([]uint, 0, len(roles))
	for _, role := range roles {
		enabledRoleIds = append(enabledRoleIds, role.Id)
	}
	if len(enabledRoleIds) == 0 {
		return false
	}

	var rolePermissions []entity.AdminRolePermission
	if err := dao.AdminRolePermission.Ctx(r.Context()).WhereIn(dao.AdminRolePermission.Columns().RoleId, enabledRoleIds).Scan(&rolePermissions); err != nil {
		return false
	}
	permissionIds := make([]uint, 0, len(rolePermissions))
	for _, relation := range rolePermissions {
		permissionIds = append(permissionIds, relation.PermissionId)
	}
	if len(permissionIds) == 0 {
		return false
	}

	permissionColumns := dao.AdminPermission.Columns()
	count, err := dao.AdminPermission.Ctx(r.Context()).
		WhereIn(permissionColumns.Id, permissionIds).
		Where(permissionColumns.Method, strings.ToUpper(method)).
		Where(permissionColumns.Path, path).
		Where(permissionColumns.Status, consts.AdminPermissionEnabled).
		WhereNull(permissionColumns.DeletedAt).
		Count()
	return err == nil && count > 0
}

func shouldMaskErrorMessage(message string) bool {
	normalized := strings.ToLower(message)
	keywords := []string{
		"select ",
		"insert ",
		"update ",
		"delete ",
		" from ",
		" where ",
		"sql",
		"unknown database",
		"no database selected",
		"error 1049",
		"error 1146",
		"access denied",
		"driver:",
		"dial tcp",
		"connection refused",
	}
	for _, keyword := range keywords {
		if strings.Contains(normalized, keyword) {
			return true
		}
	}
	return false
}

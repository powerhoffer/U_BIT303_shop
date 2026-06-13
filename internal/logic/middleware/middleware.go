package middleware

import (
	"strings"

	"bit303_shop/internal/consts"
	"bit303_shop/internal/dao"
	"bit303_shop/internal/model/entity"
	"bit303_shop/internal/service"
	"bit303_shop/utility"
	"bit303_shop/utility/response"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
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
		response.JsonExit(r, code.Code(), err.Error())
		return
	}
	response.JsonExit(r, code.Code(), "", res)
}

func (s *sMiddleware) EmployeeAuth(r *ghttp.Request) {
	token := strings.TrimSpace(r.Header.Get("Authorization"))
	if token == "" {
		response.JsonExit(r, 401, "Unauthorized or login expired")
		return
	}
	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimSpace(strings.TrimPrefix(token, "Bearer "))
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
	if err != nil {
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

package backend

import "github.com/gogf/gf/v2/frame/g"

type PermissionCreateReq struct {
	g.Meta    `path:"/permission/create" method:"post" tags:"Backend Permission" summary:"Create permission"`
	Name      string `json:"name" v:"required|length:2,128#Permission name is required|Permission name length must be 2 to 128 characters"`
	GroupName string `json:"group_name" v:"max-length:64#Group name must be at most 64 characters"`
	Method    string `json:"method" v:"required|in:GET,POST,PUT,DELETE,PATCH#Method is required|Method is invalid"`
	Path      string `json:"path" v:"required|max-length:255#Path is required|Path must be at most 255 characters"`
}

type PermissionCreateRes struct {
	Permission PermissionBase `json:"permission"`
}

type PermissionListReq struct {
	g.Meta    `path:"/permission/list" method:"get" tags:"Backend Permission" summary:"Permission list"`
	Page      int    `json:"page" in:"query" d:"1" v:"min:1#Page number is invalid"`
	Size      int    `json:"size" in:"query" d:"10" v:"max:100#Page size must be at most 100"`
	Name      string `json:"name" in:"query"`
	GroupName string `json:"group_name" in:"query"`
	Method    string `json:"method" in:"query"`
	Status    int    `json:"status" in:"query" d:"-1"`
}

type PermissionListRes struct {
	List  []PermissionBase `json:"list"`
	Total int              `json:"total"`
	Page  int              `json:"page"`
	Size  int              `json:"size"`
}

type PermissionDetailReq struct {
	g.Meta `path:"/permission/detail" method:"get" tags:"Backend Permission" summary:"Permission detail"`
	Id     uint `json:"id" in:"query" v:"required|min:1#Permission ID is required|Permission ID is invalid"`
}

type PermissionDetailRes struct {
	Permission PermissionBase `json:"permission"`
}

type PermissionUpdateReq struct {
	g.Meta    `path:"/permission/update" method:"post" tags:"Backend Permission" summary:"Update permission"`
	Id        uint   `json:"id" v:"required|min:1#Permission ID is required|Permission ID is invalid"`
	Name      string `json:"name" v:"required|length:2,128#Permission name is required|Permission name length must be 2 to 128 characters"`
	GroupName string `json:"group_name" v:"max-length:64#Group name must be at most 64 characters"`
	Method    string `json:"method" v:"required|in:GET,POST,PUT,DELETE,PATCH#Method is required|Method is invalid"`
	Path      string `json:"path" v:"required|max-length:255#Path is required|Path must be at most 255 characters"`
}

type PermissionUpdateRes struct {
	Permission PermissionBase `json:"permission"`
}

type PermissionStatusReq struct {
	g.Meta `path:"/permission/status" method:"post" tags:"Backend Permission" summary:"Enable or disable permission"`
	Id     uint `json:"id" v:"required|min:1#Permission ID is required|Permission ID is invalid"`
	Status int  `json:"status" v:"required|in:0,1#Status is required|Status must be 0 or 1"`
}

type PermissionStatusRes struct {
	Message string `json:"message"`
}

type PermissionBase struct {
	Id        uint   `json:"id"`
	Name      string `json:"name"`
	GroupName string `json:"group_name"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Status    int    `json:"status"`
}

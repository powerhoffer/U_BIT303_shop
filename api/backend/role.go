package backend

import "github.com/gogf/gf/v2/frame/g"

type RoleCreateReq struct {
	g.Meta      `path:"/role/create" method:"post" tags:"Backend Role" summary:"Create role"`
	Name        string `json:"name" v:"required|length:2,64#Role name is required|Role name length must be 2 to 64 characters"`
	Description string `json:"description" v:"max-length:255#Description must be at most 255 characters"`
}

type RoleCreateRes struct {
	Role RoleBase `json:"role"`
}

type RoleListReq struct {
	g.Meta `path:"/role/list" method:"get" tags:"Backend Role" summary:"Role list"`
	Page   int    `json:"page" in:"query" d:"1" v:"min:1#Page number is invalid"`
	Size   int    `json:"size" in:"query" d:"10" v:"max:50#Page size must be at most 50"`
	Name   string `json:"name" in:"query"`
	Status int    `json:"status" in:"query" d:"-1"`
}

type RoleListRes struct {
	List  []RoleBase `json:"list"`
	Total int        `json:"total"`
	Page  int        `json:"page"`
	Size  int        `json:"size"`
}

type RoleDetailReq struct {
	g.Meta `path:"/role/detail" method:"get" tags:"Backend Role" summary:"Role detail"`
	Id     uint `json:"id" in:"query" v:"required|min:1#Role ID is required|Role ID is invalid"`
}

type RoleDetailRes struct {
	Role RoleBase `json:"role"`
}

type RoleUpdateReq struct {
	g.Meta      `path:"/role/update" method:"post" tags:"Backend Role" summary:"Update role"`
	Id          uint   `json:"id" v:"required|min:1#Role ID is required|Role ID is invalid"`
	Name        string `json:"name" v:"required|length:2,64#Role name is required|Role name length must be 2 to 64 characters"`
	Description string `json:"description" v:"max-length:255#Description must be at most 255 characters"`
}

type RoleUpdateRes struct {
	Role RoleBase `json:"role"`
}

type RoleStatusReq struct {
	g.Meta `path:"/role/status" method:"post" tags:"Backend Role" summary:"Enable or disable role"`
	Id     uint `json:"id" v:"required|min:1#Role ID is required|Role ID is invalid"`
	Status int  `json:"status" v:"required|in:0,1#Status is required|Status must be 0 or 1"`
}

type RoleStatusRes struct {
	Message string `json:"message"`
}

type RolePermissionsReq struct {
	g.Meta        `path:"/role/permissions" method:"post" tags:"Backend Role" summary:"Assign role permissions"`
	Id            uint   `json:"id" v:"required|min:1#Role ID is required|Role ID is invalid"`
	PermissionIds []uint `json:"permission_ids"`
}

type RolePermissionsRes struct {
	Message string `json:"message"`
}

type RoleBase struct {
	Id            uint   `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Status        int    `json:"status"`
	PermissionIds []uint `json:"permission_ids"`
}

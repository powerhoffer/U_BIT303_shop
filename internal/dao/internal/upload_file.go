// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UploadFileDao is the data access object for the table upload_file.
type UploadFileDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  UploadFileColumns  // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// UploadFileColumns defines and stores column names for the table upload_file.
type UploadFileColumns struct {
	Id           string // Upload file ID
	AdminId      string // Upload admin ID
	FileName     string // Saved file name
	OriginalName string // Original file name
	FilePath     string // Local file path
	Url          string // Public access URL
	FileSize     string // File size bytes
	MimeType     string // MIME type
	FileExt      string // File extension
	BizType      string // Business type
	CreatedAt    string // Created time
}

// uploadFileColumns holds the columns for the table upload_file.
var uploadFileColumns = UploadFileColumns{
	Id:           "id",
	AdminId:      "admin_id",
	FileName:     "file_name",
	OriginalName: "original_name",
	FilePath:     "file_path",
	Url:          "url",
	FileSize:     "file_size",
	MimeType:     "mime_type",
	FileExt:      "file_ext",
	BizType:      "biz_type",
	CreatedAt:    "created_at",
}

// NewUploadFileDao creates and returns a new DAO object for table data access.
func NewUploadFileDao(handlers ...gdb.ModelHandler) *UploadFileDao {
	return &UploadFileDao{
		group:    "default",
		table:    "upload_file",
		columns:  uploadFileColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UploadFileDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UploadFileDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UploadFileDao) Columns() UploadFileColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UploadFileDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UploadFileDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *UploadFileDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}

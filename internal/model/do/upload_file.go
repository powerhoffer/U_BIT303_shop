// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UploadFile is the golang structure of table upload_file for DAO operations like Where/Data.
type UploadFile struct {
	g.Meta       `orm:"table:upload_file, do:true"`
	Id           any         // Upload file ID
	AdminId      any         // Upload admin ID
	FileName     any         // Saved file name
	OriginalName any         // Original file name
	FilePath     any         // Local file path
	Url          any         // Public access URL
	FileSize     any         // File size bytes
	MimeType     any         // MIME type
	FileExt      any         // File extension
	BizType      any         // Business type
	CreatedAt    *gtime.Time // Created time
}

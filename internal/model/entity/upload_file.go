// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UploadFile is the golang structure for table upload_file.
type UploadFile struct {
	Id           uint        `json:"id"           orm:"id"            ` // Upload file ID
	AdminId      uint        `json:"adminId"      orm:"admin_id"      ` // Upload admin ID
	FileName     string      `json:"fileName"     orm:"file_name"     ` // Saved file name
	OriginalName string      `json:"originalName" orm:"original_name" ` // Original file name
	FilePath     string      `json:"filePath"     orm:"file_path"     ` // Local file path
	Url          string      `json:"url"          orm:"url"           ` // Public access URL
	FileSize     uint64      `json:"fileSize"     orm:"file_size"     ` // File size bytes
	MimeType     string      `json:"mimeType"     orm:"mime_type"     ` // MIME type
	FileExt      string      `json:"fileExt"      orm:"file_ext"      ` // File extension
	BizType      string      `json:"bizType"      orm:"biz_type"      ` // Business type
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    ` // Created time
}

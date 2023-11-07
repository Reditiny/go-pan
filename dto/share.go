package dto

import "time"

type Share struct {
	FileId     int       `form:"fileId" valid:"Required"`
	ExpireTime time.Time `form:"expireTime" valid:"Required"`
	CodeType   int8      `form:"codeType" valid:"Required"`
	Code       string    `form:"code"`
}

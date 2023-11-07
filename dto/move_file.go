package dto

type MoveFile struct {
	FileIds string `form:"fileIds" valid:"Required"`
	FilePid int    `form:"filePid" valid:"Required"`
}

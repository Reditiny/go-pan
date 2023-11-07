package dto

type DeleteFile struct {
	FileIds string `form:"fileIds" valid:"Required"`
}

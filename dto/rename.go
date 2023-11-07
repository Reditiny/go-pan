package dto

type Rename struct {
	FileId   int    `form:"fileId" valid:"Required"`
	FileName string `form:"fileName" valid:"Required"`
}

package dto

type LoadFileList struct {
	Category string `form:"category" valid:"Required"`
	FilePid  int    `form:"filePid"`
	FileName string `form:"fileNameFuzzy"`
	PageNo   int    `form:"pageNo"`
	PageSize int    `form:"pageSize"`
}

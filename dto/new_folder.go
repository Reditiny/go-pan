package dto

type NewFolder struct {
	FilePid  int    `form:"filePid" valid:"Required"`
	FileName string `form:"fileName" valid:"Required"`
}

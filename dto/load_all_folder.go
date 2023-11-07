package dto

type LoadAllFolder struct {
	FilePid       int `form:"filePid" valid:"Required"`
	CurrentFileId int `form:"currentFileIds"`
}

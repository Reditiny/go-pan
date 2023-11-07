package vo

import "time"

type ShareList struct {
	Id           int       `json:"shareId"`
	FileId       int       `json:"fileId"`
	UserId       int       `json:"userId"`
	ValidType    int8      `json:"validType"`
	ExpireTime   time.Time `json:"expireTime"`
	ShareTime    time.Time `json:"shareTime"`
	Code         string    `json:"code"`
	ShowCount    int       `json:"showCount"`
	FileName     string    `json:"fileName"`
	FolderType   int8      `json:"folderType"`
	FileCategory int8      `json:"fileCategory"`
	FileType     int8      `json:"fileType"`
	FileCover    string    `json:"fileCover"`
}

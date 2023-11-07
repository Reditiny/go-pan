package vo

import "time"

type FileList struct {
	Id             int       `json:"fileId"`
	FilePid        int       `json:"filePid"`
	FileSize       int       `json:"fileSize"`
	FileName       string    `json:"fileName"`
	FileCover      string    `json:"fileCover"`
	CreateTime     time.Time `json:"createTime"`
	LastUpdateTime time.Time `json:"lastUpdateTime"`
	FolderType     int       `json:"folderType"`
	FileCategory   int       `json:"fileCategory"`
	FileType       int       `json:"fileType"`
	Status         int       `json:"status"`
}

package dto

type UploadFile struct {
	FileId     string `form:"fileId"`
	FileName   string `form:"fileName" valid:"Required"`
	FilePid    string `form:"filePid" valid:"Required"`
	FileMd5    string `form:"fileMd5" valid:"Required"`
	ChunkIndex string `form:"chunkIndex" valid:"Required"`
	Chunks     string `form:"chunks" valid:"Required"`
}

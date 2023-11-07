package vo

type Upload struct {
	FileId string `json:"fileId"`
	Status string `json:"status"`
}

func Uploading(fileId string) Upload {
	return Upload{
		FileId: fileId,
		Status: "uploading",
	}
}

func UploadFinish(fileId string) Upload {
	return Upload{
		FileId: fileId,
		Status: "upload_finish",
	}
}

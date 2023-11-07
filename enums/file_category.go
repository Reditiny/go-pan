package enums

type FileCategory struct {
	Code         int8
	CategoryName string
	Desc         string
}

var (
	CATEGORY_VIDEO = FileCategory{1, "video", "视频"}
	CATEGORY_MUSIC = FileCategory{2, "music", "音频"}
	CATEGORY_IMAGE = FileCategory{3, "image", "图片"}
	CATEGORY_DOC   = FileCategory{4, "doc", "文档"}
	CATEGORY_OTHER = FileCategory{5, "other", "其他"}
)

func GetCodeByCategoryName(name string) int8 {
	switch name {
	case CATEGORY_VIDEO.CategoryName:
		return CATEGORY_VIDEO.Code
	case CATEGORY_MUSIC.CategoryName:
		return CATEGORY_MUSIC.Code
	case CATEGORY_IMAGE.CategoryName:
		return CATEGORY_IMAGE.Code
	case CATEGORY_DOC.CategoryName:
		return CATEGORY_DOC.Code
	case CATEGORY_OTHER.CategoryName:
		return CATEGORY_OTHER.Code
	default:
		return 0
	}
}

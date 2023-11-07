package enums

type FileType struct {
	FileCategory
	fileType int8
	Ext      []string
}

var (
	TYPE_VIDEO    = FileType{CATEGORY_VIDEO, 1, []string{"mp4", "avi", "rmvb", "rm", "flv", "wmv", "mkv", "mov", "3gp", "mpeg"}}
	TYPE_MUSIC    = FileType{CATEGORY_MUSIC, 2, []string{"mp3", "wav", "wma", "ogg", "ape", "flac", "aac", "m4a"}}
	TYPE_IMAGE    = FileType{CATEGORY_IMAGE, 3, []string{"jpg", "jpeg", "png", "gif", "bmp", "webp"}}
	TYPE_PDF      = FileType{CATEGORY_DOC, 4, []string{"pdf"}}
	TYPE_WORD     = FileType{CATEGORY_DOC, 5, []string{"doc", "docx"}}
	TYPE_EXCEL    = FileType{CATEGORY_DOC, 6, []string{"xls", "xlsx"}}
	TYPE_TXT      = FileType{CATEGORY_DOC, 7, []string{"txt"}}
	TYPE_PROGRAME = FileType{CATEGORY_OTHER, 8, []string{"go", "java", "c", "cpp", "py", "js", "html", "css", "php", "sql", "sh", "bat"}}
	TYPE_ZIP      = FileType{CATEGORY_OTHER, 9, []string{"zip", "rar", "7z", "tar", "gz"}}
	TYPE_OTHER    = FileType{CATEGORY_OTHER, 10, []string{}}
)

var typeList = []FileType{TYPE_VIDEO, TYPE_MUSIC, TYPE_IMAGE, TYPE_PDF, TYPE_WORD, TYPE_EXCEL, TYPE_TXT, TYPE_PROGRAME, TYPE_ZIP, TYPE_OTHER}

func GetTypeAndCategoryByExt(ext string) (int8, int8) {
	for _, t := range typeList {
		for _, e := range t.Ext {
			if e == ext[1:] {
				return t.FileCategory.Code, t.fileType
			}
		}
	}
	return CATEGORY_OTHER.Code, TYPE_OTHER.fileType
}

package dto

type Path struct {
	Path string `form:"path" valid:"Required"`
}

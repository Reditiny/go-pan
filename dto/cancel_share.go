package dto

type CancelShare struct {
	ShareIds string `form:"shareIds" valid:"Required"`
}

package dto

type Password struct {
	Password string `form:"password" valid:"Required"`
}

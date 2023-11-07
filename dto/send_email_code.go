package dto

type SendEmailCode struct {
	Email     string `form:"email" valid:"Required;Email"`
	Type      string `form:"type" valid:"Required"`
	CheckCode string `form:"checkCode" valid:"Required"`
}

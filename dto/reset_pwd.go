package dto

type ResetPwd struct {
	Email     string `form:"email" valid:"Required;Email"`
	PassWord  string `form:"password" valid:"Required"`
	CheckCode string `form:"checkCode" valid:"Required"`
	EmailCode string `form:"emailCode" valid:"Required"`
}

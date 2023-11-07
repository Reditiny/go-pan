package dto

type Login struct {
	Email     string `form:"email" valid:"Required;Email"`
	Password  string `form:"password" valid:"Required"`
	CheckCode string `form:"checkCode" valid:"Required"`
}

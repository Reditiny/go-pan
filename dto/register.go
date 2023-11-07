package dto

type Register struct {
	Email     string `form:"email" valid:"Required;Email"`
	NickName  string `form:"nickName" valid:"Required"`
	PassWord  string `form:"password" valid:"Required"`
	CheckCode string `form:"checkCode" valid:"Required"`
	EmailCode string `form:"emailCode" valid:"Required"`
}

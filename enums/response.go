package enums

type ResponseCode struct {
	Code int
	Info string
}

var (
	// CODE_200 = ResponseCode{200, "请求成功"}
	CODE_100 = ResponseCode{100, "功能待开发"}
	CODE_404 = ResponseCode{404, "请求地址不存在"}
	CODE_600 = ResponseCode{600, "请求参数错误"}
	CODE_601 = ResponseCode{601, "昵称已存在"}
	CODE_602 = ResponseCode{602, "验证码错误"}
	CODE_603 = ResponseCode{603, "用户名或密码错误"}
	CODE_604 = ResponseCode{604, "邮箱已注册"}
	CODE_605 = ResponseCode{605, "邮箱未注册"}
	CODE_606 = ResponseCode{606, "读取文件失败"}
	CODE_607 = ResponseCode{606, "文件上传失败"}
	CODE_500 = ResponseCode{500, "服务器返回错误，请联系管理员"}
	CODE_901 = ResponseCode{901, "登录已过期，请重新登录"}
	CODE_902 = ResponseCode{902, "分享连接不存在，或者已失效"}
	CODE_903 = ResponseCode{903, "分享验证失效，请重新验证"}
	CODE_904 = ResponseCode{904, "网盘空间不足，请扩容"}
)

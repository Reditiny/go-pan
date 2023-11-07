package controllers

import (
	"context"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	response "gopan/common"
	"gopan/constants"
	"gopan/dto"
	"gopan/enums"
	"gopan/models"
	"gopan/utils"
	"gopan/vo"
	_ "image/jpeg"
	"os"
	"strconv"
	"time"
)

// UserController operations for User
type UserController struct {
	beego.Controller
}

// CheckCode 获取验证码 0:登录注册 1:邮箱验证码发送 默认0
func (c *UserController) CheckCode() {
	codeType := c.GetString("type")
	captcha, s := utils.GenerateImageCaptcha()
	if codeType == "" || codeType == "0" {
		utils.SetToRedis[string](constants.CHECK_CODE_LOGIN_REGISTER, &s, constants.EXPIRE_FIVE_MINUTES)
	} else {
		utils.SetToRedis[string](constants.CHECK_CODE_SEND_EMAIL, &s, constants.EXPIRE_FIVE_MINUTES)
	}
	captcha.WriteTo(c.Ctx.ResponseWriter)
}

// SendEmailCode 发送邮箱验证码
func (c *UserController) SendEmailCode() {
	defer c.ServeJSON()
	err, sendEmailCode := bindAndValidate[dto.SendEmailCode](c.Controller)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_600)
		return
	}
	if !checkCodeType1(sendEmailCode.CheckCode) {
		c.Data["json"] = response.Fail(enums.CODE_602)
		return
	}
	numberCaptcha := utils.GenerateNumberCaptcha(constants.LENGTN_SIX)
	utils.SetToRedis[string](constants.EMAIL_CODE_PREFIX_EMAIL+sendEmailCode.Email, &numberCaptcha, constants.EXPIRE_FIVE_MINUTES)
	go func() {
		utils.SendEmail(sendEmailCode.Email, numberCaptcha)
	}()
	c.Data["json"] = response.Success(nil)
}

// Register 用户注册
func (c *UserController) Register() {
	defer c.ServeJSON()
	err, register := bindAndValidate[dto.Register](c.Controller)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_600)
		return
	}
	if !(checkCodeType0(register.CheckCode) && checkEmailCode(register.Email, register.EmailCode)) {
		c.Data["json"] = response.Fail(enums.CODE_602)
		return
	}
	user := models.User{Email: register.Email}
	err = models.GetUserByEmail(&user)
	if err == nil {
		c.Data["json"] = response.Fail(enums.CODE_604)
		return
	}
	user = models.User{NickName: register.NickName}
	models.GetUserByNickname(&user)
	if err == nil {
		c.Data["json"] = response.Fail(enums.CODE_601)
		return
	}
	err = orm.NewOrm().DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		user := utils.CopyOne[dto.Register, models.User](*register)
		user.JoinTime = time.Now()
		user.QqAvatar = utils.GetDefaultAvatarPath()
		id, err := txOrm.Insert(&user)
		if err != nil {
			return err
		}
		space := models.UserSpace{UserId: int(id), UseSpace: 0, TotalSpace: 1024 * 1024 * 1024}
		_, err = txOrm.Insert(&space)
		return err
	})
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_500)
		return
	}
	c.Data["json"] = response.Success(nil)
}

// Login 用户登录
func (c *UserController) Login() {
	defer c.ServeJSON()
	err, login := bindAndValidate[dto.Login](c.Controller)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_600)
		return
	}
	//if !checkCodeType0(login.CheckCode) {
	//	c.Data["json"] = response.Fail(enums.CODE_602)
	//	return
	//}
	user := models.User{Email: login.Email, Password: login.Password}
	err = models.GetUserByEmailAndPassword(&user)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_603)
		return
	}
	user.LastLoginTime = time.Now()
	models.UpdateUserById(&user)
	loginVo := utils.CopyOne[models.User, vo.Login](user)
	c.Data["json"] = response.Success(loginVo)
	c.SetSession(constants.SESSION_USER_ID_KEY, loginVo.Id)
}

// ResetPwd 重置密码
func (c *UserController) ResetPwd() {
	defer c.ServeJSON()
	err, resetPwd := bindAndValidate[dto.ResetPwd](c.Controller)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_600)
		return
	}
	if !(checkCodeType0(resetPwd.CheckCode) && checkEmailCode(resetPwd.Email, resetPwd.EmailCode)) {
		c.Data["json"] = response.Fail(enums.CODE_602)
		return
	}
	user := models.User{Email: resetPwd.Email}
	err = models.GetUserByEmail(&user)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_605)
		return
	}
	user.Password = resetPwd.PassWord
	err = models.UpdateUserById(&user)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_500)
		return
	}
	c.Data["json"] = response.Success(nil)
}

// GetAvatar 获取用户头像
func (c *UserController) GetAvatar() {
	defer c.ServeJSON()
	userId, err := strconv.Atoi(c.Ctx.Input.Param(":userId"))
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_901)
		return
	}
	user, err := models.GetUserById(userId)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_500)
		return
	}
	fileBytes, err := os.ReadFile(utils.GetAvatarPath() + user.QqAvatar)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_606)
		return
	}
	//c.Ctx.ResponseWriter.Header().Set("Content-Type", "image/jpeg")
	c.Ctx.ResponseWriter.Write(fileBytes)
}

// GetUseSpace 获取用户空间使用情况
func (c *UserController) GetUseSpace() {
	defer c.ServeJSON()
	userId := c.GetSession(constants.SESSION_USER_ID_KEY)
	if userId == nil {
		c.Data["json"] = response.Fail(enums.CODE_901)
		return
	}
	userSpace := models.UserSpace{UserId: userId.(int)}
	err := models.GetUserSpaceByUserId(&userSpace)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_500)
		return
	}
	space := utils.CopyOne[models.UserSpace, vo.Space](userSpace)
	c.Data["json"] = response.Success(space)
}

// Logout 退出登陆
func (c *UserController) Logout() {
	defer c.ServeJSON()
	c.DelSession(constants.SESSION_USER_ID_KEY)
	c.Data["json"] = response.Success(nil)
}

// UpdatePassword 修改密码
func (c *UserController) UpdatePassword() {
	defer c.ServeJSON()
	userId := c.GetSession(constants.SESSION_USER_ID_KEY)
	if userId == nil {
		c.Data["json"] = response.Fail(enums.CODE_901)
		return
	}
	err, newPassword := bindAndValidate[dto.Password](c.Controller)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_600)
		return
	}
	user := models.User{Id: userId.(int)}
	user.Password = newPassword.Password
	err = models.UpdateUserPassword(&user)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_500)
		return
	}
	c.Data["json"] = response.Success(nil)
}

// UpdateAvatar 修改头像
func (c *UserController) UpdateAvatar() {
	defer c.ServeJSON()
	userId := c.GetSession(constants.SESSION_USER_ID_KEY)
	if userId == nil {
		c.Data["json"] = response.Fail(enums.CODE_901)
		return
	}
	f, h, err := c.GetFile("avatar")
	defer f.Close()
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_600)
		return
	}
	fileName := utils.GetTimePath() + "/" + h.Filename
	c.SaveToFile("avatar", utils.ToAvatarPath(fileName))
	user := models.User{Id: userId.(int), QqAvatar: fileName}
	err = models.UpdateUserAvatarById(&user)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_500)
		return
	}
	c.Data["json"] = response.Success(nil)
}

// TODO QQ登录

func checkCodeType0(code string) bool {
	cache, err := utils.GetFromRedis[string](constants.CHECK_CODE_LOGIN_REGISTER)
	if err != nil {
		return false
	}
	return *cache == code
}

func checkCodeType1(code string) bool {
	cache, err := utils.GetFromRedis[string](constants.CHECK_CODE_SEND_EMAIL)
	if err != nil {
		return false
	}
	return *cache == code
}

func checkEmailCode(email, code string) bool {
	cache, err := utils.GetFromRedis[string](constants.EMAIL_CODE_PREFIX_EMAIL + email)
	if err != nil {
		return false
	}
	return *cache == code
}

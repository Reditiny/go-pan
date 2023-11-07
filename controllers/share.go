package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	response "gopan/common"
	"gopan/constants"
	"gopan/dto"
	"gopan/enums"
	"gopan/models"
	"gopan/utils"
	"gopan/vo"
	"strconv"
	"strings"
	"time"
)

// ShareController operations for Share
type ShareController struct {
	beego.Controller
}

// LoadShareList 获取分享列表
func (c *ShareController) LoadShareList() {
	defer c.ServeJSON()
	userId := c.GetSession(constants.SESSION_USER_ID_KEY)
	if userId == nil {
		c.Data["json"] = response.Fail(enums.CODE_901)
		return
	}
	err, page := bindAndValidate[dto.LoadShareList](c.Controller)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_600)
		return
	}
	err, shares := models.GetShareByUserId(userId.(int))
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_500)
		return
	}
	voList := utils.CopyList[models.Share, vo.ShareList](shares)
	for _, shareVo := range voList {
		file, err := models.GetFileById(shareVo.FileId)
		if err != nil {
			c.Data["json"] = response.Fail(enums.CODE_500)
			return
		}
		shareVo.FolderType = file.FolderType
		shareVo.FileName = file.FileName
		shareVo.FileCategory = file.FileCategory
		shareVo.FileType = file.FileType
		shareVo.FileCover = file.FileCover
	}
	pageVo := utils.MakePage[vo.ShareList](voList, page.PageSize, page.PageNo)
	c.Data["json"] = response.Success(pageVo)
}

// ShareFile 分享文件
func (c *ShareController) ShareFile() {
	defer c.ServeJSON()
	userId := c.GetSession(constants.SESSION_USER_ID_KEY)
	if userId == nil {
		c.Data["json"] = response.Fail(enums.CODE_901)
		return
	}
	err, s := bindAndValidate[dto.Share](c.Controller)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_600)
		return
	}
	if s.CodeType == 1 || s.Code == "" {
		s.Code = utils.GenerateNumberCaptcha(5)
	}
	s.ExpireTime = time.Now().AddDate(0, 0, 7)
	share := utils.CopyOne[dto.Share, models.Share](*s)
	share.UserId = userId.(int)
	_, err = models.AddShare(&share)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_500)
		return
	}
	c.Data["json"] = response.Success(share)
}

// CancelShare 取消分享
func (c *ShareController) CancelShare() {
	defer c.ServeJSON()
	err, cancel := bindAndValidate[dto.CancelShare](c.Controller)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_600)
		return
	}
	ids := strings.Split(cancel.ShareIds, ",")
	for _, idStr := range ids {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.Data["json"] = response.Fail(enums.CODE_600)
			return
		}
		err = models.DeleteShare(id)
		if err != nil {
			c.Data["json"] = response.Fail(enums.CODE_500)
			return
		}
	}
	c.Data["json"] = response.Success(nil)
}

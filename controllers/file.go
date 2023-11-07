package controllers

import (
	"context"
	"github.com/beego/beego/v2/client/orm"
	response "gopan/common"
	"gopan/constants"
	"gopan/dto"
	"gopan/enums"
	"gopan/models"
	"gopan/utils"
	"gopan/vo"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

// FileController operations for File
type FileController struct {
	beego.Controller
}

// LoadList 获取文件列表
func (c *FileController) LoadList() {
	defer c.ServeJSON()
	userId := c.GetSession(constants.SESSION_USER_ID_KEY)
	if userId == nil {
		c.Data["json"] = response.Fail(enums.CODE_901)
		return
	}
	err, cond := bindAndValidate[dto.LoadFileList](c.Controller)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_600)
		return
	}
	err, files := models.GetFileList(cond, userId.(int))
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_500)
		return
	}
	list := utils.CopyList[models.File, vo.FileList](files)
	page := utils.MakePage[vo.FileList](list, cond.PageSize, cond.PageNo)
	c.Data["json"] = response.Success(page)
}

// UploadFile 上传文件 TODO 完善文件上传
func (c *FileController) UploadFile() {
	defer c.ServeJSON()
	userId := c.GetSession(constants.SESSION_USER_ID_KEY)
	if userId == nil {
		c.Data["json"] = response.Fail(enums.CODE_901)
		return
	}
	_, h, err1 := c.GetFile("file")
	err2, upload := bindAndValidate[dto.UploadFile](c.Controller)
	chunkIndex, err3 := strconv.Atoi(upload.ChunkIndex)
	chunks, err4 := strconv.Atoi(upload.Chunks)
	pid, err5 := strconv.Atoi(upload.FilePid)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil || chunkIndex >= chunks {
		c.Data["json"] = response.Fail(enums.CODE_600)
		return
	}
	if chunkIndex == 0 {
		// TODO 根据 md5 值实现秒传
		upload.FileId = "1234"
		fileSize := h.Size * int64(chunks)
		userSpace := models.UserSpace{UserId: userId.(int)}
		models.GetUserSpaceByUserId(&userSpace)
		if userSpace.UseSpace+fileSize > userSpace.TotalSpace {
			c.Data["json"] = response.Fail(enums.CODE_904)
			return
		}
	}
	c.SaveToFile("file", utils.GetTempPath()+"/"+upload.ChunkIndex)
	if chunkIndex != chunks-1 {
		c.Data["json"] = response.Success(vo.Uploading(upload.FileId))
	} else {
		fileId, _ := strconv.Atoi(upload.FileId)
		infoChan := make(chan enums.Info, 1)
		go func() {
			mergeFile(infoChan, upload.FileName, chunks, pid, userId.(int), fileId, upload.FileMd5)
		}()
		info := <-infoChan
		if info.Code == enums.FAIL_INFO.Code {
			c.Data["json"] = response.Fail(enums.CODE_607)
		} else {
			c.Data["json"] = response.Success(vo.UploadFinish(upload.FileId))
		}
	}
}

// NewFolder 新建文件夹
func (c *FileController) NewFolder() {
	defer c.ServeJSON()
	userId := c.GetSession(constants.SESSION_USER_ID_KEY)
	if userId == nil {
		c.Data["json"] = response.Fail(enums.CODE_901)
		return
	}
	err, newFolder := bindAndValidate[dto.NewFolder](c.Controller)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_600)
		return
	}
	file := models.File{
		FileName:       newFolder.FileName,
		FilePid:        newFolder.FilePid,
		UserId:         userId.(int),
		FolderType:     1,
		CreateTime:     time.Now(),
		LastUpdateTime: time.Now(),
		RecoveryTime:   time.Now(),
		Status:         2,
	}
	models.AddFile(&file)
}

// GetFolderInfo 获取文件夹信息
func (c *FileController) GetFolderInfo() {
	defer c.ServeJSON()
	userId := c.GetSession(constants.SESSION_USER_ID_KEY)
	if userId == nil {
		c.Data["json"] = response.Fail(enums.CODE_901)
		return
	}
	err, path := bindAndValidate[dto.Path](c.Controller)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_600)
		return
	}
	folders := strings.Split(path.Path, "/")
	var data = make([]*vo.Folder, len(folders))
	for i, folder := range folders {
		id, err := strconv.Atoi(folder)
		if err != nil {
			c.Data["json"] = response.Fail(enums.CODE_600)
			return
		}
		file, _ := models.GetFileById(id)
		vo := utils.CopyOne[models.File, vo.Folder](*file)
		data[i] = &vo
	}
	c.Data["json"] = response.Success(data)
}

// Rename 重命名文件
func (c *FileController) Rename() {
	defer c.ServeJSON()
	userId := c.GetSession(constants.SESSION_USER_ID_KEY)
	if userId == nil {
		c.Data["json"] = response.Fail(enums.CODE_901)
		return
	}
	err, rename := bindAndValidate[dto.Rename](c.Controller)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_600)
		return
	}
	file := models.File{FileName: rename.FileName, Id: rename.FileId, LastUpdateTime: time.Now()}
	err = models.UpdateFileName(&file)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_500)
		return
	}
	c.Data["json"] = response.Success(nil)
}

// LoadAllFolder 获取所有文件夹
func (c *FileController) LoadAllFolder() {
	defer c.ServeJSON()
	userId := c.GetSession(constants.SESSION_USER_ID_KEY)
	if userId == nil {
		c.Data["json"] = response.Fail(enums.CODE_901)
		return
	}
	err, loadAllFolder := bindAndValidate[dto.LoadAllFolder](c.Controller)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_600)
		return
	}
	err, files := models.GetFoldersByPid(loadAllFolder.FilePid)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_500)
		return
	}
	data := utils.CopyList[models.File, vo.FileList](files)
	c.Data["json"] = response.Success(data)
}

// ChangeFileFolder 移动文件
func (c *FileController) ChangeFileFolder() {
	defer c.ServeJSON()
	userId := c.GetSession(constants.SESSION_USER_ID_KEY)
	if userId == nil {
		c.Data["json"] = response.Fail(enums.CODE_901)
		return
	}
	err, changeFileFolder := bindAndValidate[dto.MoveFile](c.Controller)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_600)
		return
	}
	fileIds := strings.Split(changeFileFolder.FileIds, ",")
	for _, fileId := range fileIds {
		id, err := strconv.Atoi(fileId)
		if err != nil {
			c.Data["json"] = response.Fail(enums.CODE_600)
			return
		}
		file := models.File{Id: id, FilePid: changeFileFolder.FilePid, LastUpdateTime: time.Now()}
		models.UpdatePidById(&file)
	}
	c.Data["json"] = response.Success(nil)
}

// CreateDownloadUrl 生成下载链接
func (c *FileController) CreateDownloadUrl() {
	defer c.ServeJSON()
	userId := c.GetSession(constants.SESSION_USER_ID_KEY)
	if userId == nil {
		c.Data["json"] = response.Fail(enums.CODE_901)
		return
	}
	fileId := c.Ctx.Input.Param(":fileId")
	if fileId == "" {
		c.Data["json"] = response.Fail(enums.CODE_600)
		return
	}
	c.Data["json"] = response.Success(fileId)
}

// Download 下载文件
func (c *FileController) Download() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":code"))
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_600)
		return
	}
	err, path := models.GetFilePathById(id)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_500)
		return
	}
	c.Ctx.Output.Download(utils.GetFilePath() + path)
}

// DeleteFile 删除文件
func (c *FileController) DeleteFile() {
	defer c.ServeJSON()
	userId := c.GetSession(constants.SESSION_USER_ID_KEY)
	if userId == nil {
		c.Data["json"] = response.Fail(enums.CODE_901)
		return
	}
	err, fileIds := bindAndValidate[dto.DeleteFile](c.Controller)
	if err != nil {
		c.Data["json"] = response.Fail(enums.CODE_600)
		return
	}
	ids := strings.Split(fileIds.FileIds, ",")
	for _, fileId := range ids {
		id, err := strconv.Atoi(fileId)
		if err != nil {
			c.Data["json"] = response.Fail(enums.CODE_600)
			return
		}
		file, err := models.GetFileById(id)
		if err != nil {
			c.Data["json"] = response.Fail(enums.CODE_500)
			return
		}
		models.DeleteFile(id)
		if file.FolderType == 0 {
			go func() {
				deleteFile(file.FilePath)
			}()
		}
	}
	c.Data["json"] = response.Success(nil)
}

// mergeFile 合并临时文件
func mergeFile(infoChan chan enums.Info, fileName string, chunks int, pid int, userId int, fileId int, fileMd5 string) {
	path := utils.GetFilePath()
	timePath := utils.GetTimePath()
	tempPath := utils.GetTempPath()
	utils.MakePathIfNotExist(path + timePath)
	file, _ := os.Create(path + timePath + "/" + fileName)
	defer file.Close()
	for i := 0; i < chunks; i++ {
		tempFilePath := tempPath + "/" + strconv.Itoa(i)
		tempFile, _ := os.Open(tempFilePath)
		buf := make([]byte, 1024)
		for {
			n, _ := tempFile.Read(buf)
			if n == 0 {
				break
			}
			file.Write(buf[:n])
		}
		tempFile.Close()
		os.Remove(tempFilePath)
	}
	os.Remove(tempPath)
	info, _ := os.Stat(path + timePath + "/" + fileName)
	fileCategory, fileType := enums.GetTypeAndCategoryByExt(filepath.Ext(fileName))
	fileInfo := models.File{
		FileName:       fileName,
		FilePid:        pid,
		FileMd5:        fileMd5,
		UserId:         userId,
		FolderType:     0,
		FileCategory:   fileCategory,
		FileType:       fileType,
		FileSize:       info.Size(),
		FilePath:       timePath + "/" + fileName,
		CreateTime:     time.Now(),
		LastUpdateTime: time.Now(),
		RecoveryTime:   time.Now(),
		Status:         2,
	}
	userSpace := models.UserSpace{UserId: userId}
	models.GetUserSpaceByUserId(&userSpace)
	userSpace.UseSpace += info.Size()
	err := orm.NewOrm().DoTx(func(c context.Context, o orm.TxOrmer) error {
		_, err := o.Insert(&fileInfo)
		if err != nil {
			return err
		}
		return models.UpdateUserSpaceWithOrm(&userSpace, o)
	})
	if err != nil {
		infoChan <- enums.FAIL_INFO
	} else {
		infoChan <- enums.SUCCESS_INFO
	}
}

// 删除文件
func deleteFile(filePath string) {
	os.Remove(utils.GetFilePath() + filePath)
}

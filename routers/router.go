// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"gopan/controllers"
)

func init() {

	//ns := beego.NewNamespace("/api",
	//
	//	beego.NSNamespace("/file",
	//		beego.NSInclude(
	//			&controllers.FileController{},
	//		),
	//	),
	//
	//	beego.NSNamespace("/share",
	//		beego.NSInclude(
	//			&controllers.ShareController{},
	//		),
	//	),
	//
	//	beego.NSNamespace("/",
	//		beego.NSInclude(
	//			&controllers.UserController{},
	//		),
	//	),
	//
	//	beego.NSNamespace("/user_space",
	//		beego.NSInclude(
	//			&controllers.UserSpaceController{},
	//		),
	//	),
	//)
	//beego.AddNamespace(ns)
	initUserRouters()
	initFileRouters()
	initShareRouters()

}

func initUserRouters() {
	beego.CtrlGet("/api/checkCode", (*controllers.UserController).CheckCode)
	beego.CtrlPost("/api/sendEmailCode", (*controllers.UserController).SendEmailCode)
	beego.CtrlPost("/api/register", (*controllers.UserController).Register)
	beego.CtrlPost("/api/login", (*controllers.UserController).Login)
	beego.CtrlPost("/api/resetPwd", (*controllers.UserController).ResetPwd)
	beego.CtrlGet("/api/getAvatar/:userId", (*controllers.UserController).GetAvatar)
	beego.CtrlPost("/api/getUseSpace", (*controllers.UserController).GetUseSpace)
	beego.CtrlPost("/api/logout", (*controllers.UserController).Logout)
	beego.CtrlPost("/api/updatePassword", (*controllers.UserController).UpdatePassword)
	beego.CtrlPost("/api/updateUserAvatar", (*controllers.UserController).UpdateAvatar)
}

func initFileRouters() {
	beego.CtrlPost("/api/file/loadDataList", (*controllers.FileController).LoadList)
	beego.CtrlPost("/api/file/uploadFile", (*controllers.FileController).UploadFile)
	beego.CtrlPost("/api/file/newFoloder", (*controllers.FileController).NewFolder)
	beego.CtrlPost("/api/file/getFolderInfo", (*controllers.FileController).GetFolderInfo)
	beego.CtrlPost("/api/file/rename", (*controllers.FileController).Rename)
	beego.CtrlPost("/api/file/loadAllFolder", (*controllers.FileController).LoadAllFolder)
	beego.CtrlPost("/api/file/changeFileFolder", (*controllers.FileController).ChangeFileFolder)
	beego.CtrlPost("/api/file/createDownloadUrl/:fileId", (*controllers.FileController).CreateDownloadUrl)
	beego.CtrlGet("/api/file/download/:code", (*controllers.FileController).Download)
	beego.CtrlPost("/api/file/delFile", (*controllers.FileController).DeleteFile)
}

func initShareRouters() {
	beego.CtrlPost("/api/share/loadShareList", (*controllers.ShareController).LoadShareList)
	beego.CtrlPost("/api/share/shareFile", (*controllers.ShareController).ShareFile)
	beego.CtrlPost("/api/share/cancelShare", (*controllers.ShareController).CancelShare)
}

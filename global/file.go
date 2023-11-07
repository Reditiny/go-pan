package global

import "github.com/beego/beego/v2/server/web"

// StoreFolder 存储文件夹
var StoreFolder, _ = web.AppConfig.String("storefolder")

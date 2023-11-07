package main

import (
	_ "gopan/routers"

	beeLogger "github.com/beego/bee/v2/logger"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
	"gopan/global"
)

func main() {
	sqlConn, err := beego.AppConfig.String("sqlconn")
	if err != nil {
		beeLogger.Log.Fatal("sqlconn error: ")
	}
	orm.RegisterDataBase("default", "mysql", sqlConn)
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	global.InitRedis()
	beego.Run()
}

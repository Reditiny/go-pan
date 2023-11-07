package controllers

import (
	"errors"
	beego "github.com/beego/beego/v2/server/web"
	"gopan/global"
)

func bindAndValidate[V any](controller beego.Controller) (error, *V) {
	var requestBody V
	err1 := controller.BindForm(&requestBody)
	_, err2 := global.Validator.Valid(&requestBody)
	if err1 != nil || err2 != nil {
		return errors.New("参数错误"), nil
	}
	return nil, &requestBody
}

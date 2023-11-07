package utils

import (
	"gopan/global"
	"math/rand"
	"os"
	"strings"
	"time"
)

func GetAvatarPath() string {
	return MakePathIfNotExist(global.StoreFolder + "/avatar")
}

func GetTempPath() string {
	return MakePathIfNotExist(global.StoreFolder + "/temp")
}

func GetFilePath() string {
	return MakePathIfNotExist(global.StoreFolder + "/file")
}

func GetDefaultAvatarPath() string {
	return "/default.jpg"
}

func GetTimePath() string {
	return "/" + time.Now().Format("2006/01/02")
}

func ToAvatarPath(fileName string) string {
	exists := PathExists(GetAvatarPath() + GetTimePath())
	if !exists {
		os.MkdirAll(GetAvatarPath()+GetTimePath(), os.ModePerm)
	}
	return GetAvatarPath() + "/" + fileName
}

func MakePathIfNotExist(path string) string {
	exists := PathExists(path)
	if !exists {
		os.MkdirAll(path, os.ModePerm)
	}
	return path
}

// PathExists 判断一个文件或文件夹是否存在
// 输入文件路径，根据返回的bool值来判断文件或文件夹是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

// 获取随机字母+数字组合字符串
func getRandstring(length int) string {
	if length < 1 {
		return ""
	}
	char := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charArr := strings.Split(char, "")
	charlen := len(charArr)
	ran := rand.New(rand.NewSource(time.Now().Unix()))
	var rchar string = ""
	for i := 1; i <= length; i++ {
		rchar = rchar + charArr[ran.Intn(charlen)]
	}
	return rchar
}

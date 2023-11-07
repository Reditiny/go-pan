package utils

import "github.com/jinzhu/copier"

/**
类型拷贝工具  TODO 错误处理
*/

func CopyOne[R any, T any](res R) T {
	var tar T
	copier.Copy(&tar, &res)
	return tar
}

func CopyList[R any, T any](res []*R) []*T {
	var tar []*T
	copier.Copy(&tar, &res)
	return tar
}

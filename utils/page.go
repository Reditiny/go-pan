package utils

import "gopan/vo"

func MakePage[T any](list []*T, pageSize int, pageNo int) vo.Page[T] {
	if pageSize == 0 {
		pageSize = 10
	}
	if pageNo == 0 {
		pageNo = 1
	}
	totalCount := len(list)
	pageTotal := 0
	if totalCount > 0 {
		pageTotal = totalCount/pageSize + 1
	}
	page := vo.Page[T]{
		TotalCount: totalCount,
		PageSize:   pageSize,
		PageNo:     pageNo,
		PageTotal:  pageTotal,
		List:       make([]*T, 0),
	}
	if pageNo <= pageTotal {
		start := (pageNo - 1) * pageSize
		end := pageNo * pageSize
		if end >= len(list) {
			end = len(list)
		}
		page.List = list[start:end]
	}
	return page
}

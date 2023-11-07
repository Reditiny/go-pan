package vo

type Page[T any] struct {
	TotalCount int  `json:"totalCount"`
	PageSize   int  `json:"pageSize"`
	PageNo     int  `json:"pageNo"`
	PageTotal  int  `json:"pageTotal"`
	List       []*T `json:"list"`
}

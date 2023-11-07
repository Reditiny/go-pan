package enums

type Info struct {
	Code    int
	Message string
}

var (
	SUCCESS_INFO = Info{0, "操作成功"}
	FAIL_INFO    = Info{1, "操作失败"}
)

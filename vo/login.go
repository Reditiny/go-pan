package vo

type Login struct {
	Id       int    `json:"userId"`
	NickName string `json:"nickName"`
	QqAvatar string `json:"avatar"`
	Admin    bool   `json:"admin"`
}

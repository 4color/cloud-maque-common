package model

type UserTokenUserId struct {
	Token       string
	UserId      string
	UserName    string
	Authorities []string `json:"authorities"` //授权角色
}

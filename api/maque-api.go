package api

import (
	"cloud-maque-common/domain"
	"cloud-maque-common/utils"
	"github.com/gin-gonic/gin"
)

type ApiMaque struct {
}

// 获取超管用户的角色
func (p *ApiMaque) GetSdpRolesByUserId(gc *gin.Context, maqueUrl string, roleCode string) (haveRole bool, err error) {

	token, err := domain.GetToken(gc, maqueUrl)
	if err != nil {
		return
	}

	if token.UserId == "" {
		err = &utils.Error{ErrMsg: "Token获取不正确，可能是麻雀版本号不匹配或会话已失效", ErrCode: 500}
		return
	}
	for i := range token.Authorities {
		if token.Authorities[i] == roleCode {
			haveRole = true
			return
		}
	}

	return
}

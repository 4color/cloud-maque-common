package domain

import (
	"cloud-maque-common/consts"
	"cloud-maque-common/model"
	"cloud-maque-common/utils"
	"cloud-maque-common/utils/netutils"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func GetToken(gc *gin.Context, maqueapi string) (token model.UserTokenUserId, err error) {

	tkid := gc.GetHeader("tkid")
	if tkid == "" {
		tkid = gc.DefaultQuery("tkid", "")
	}
	token1 := gc.GetHeader("maque-token")
	if token1 == "" {
		token1 = gc.DefaultQuery("tk", "")
	}

	if token1 != "" {
		return
	}

	if tkid != "" {

		header := "[{\"mq-token\":\"" + token1 + "\"},{\"mq-tokenid\":\"" + tkid + "\"}]"
		result, err2 := netutils.GetWebRequestGetWithHeader(maqueapi+consts.URL_MAQUE_TOKE, header, "")
		if err2 != nil {
			err = err2
			return
		}

		res := model.NewResponseBodyModel()
		err = json.Unmarshal([]byte(result), &res)
		if err != nil {
			return
		}

		if res.Status == 200 {

			mapv := res.Data.(map[string]interface{})
			bytes, err3 := json.Marshal(mapv)
			if err3 != nil {
				err = err2
				return
			}

			err = json.Unmarshal(bytes, &token)

		} else {
			err = utils.NewErrorDefault(res.Message)
		}
	}

	return
}

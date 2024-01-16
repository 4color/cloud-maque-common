package model

type ResponseBodyModel struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Detail  string      `json:"detail"`
	Data    interface{} `json:"data"`
}

func NewResponseBodyModel() *ResponseBodyModel {
	r := new(ResponseBodyModel)
	r.Status = 501
	r.Message = "未定义错误"
	return r
}

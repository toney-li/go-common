package model

type JsonResp struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

const (
	// resp code/msg
	SuccCode = "0"
	SuccMsg  = "操作成功"
	FailCode = "1"
	FailMsg  = "操作失败"
)

func DefaultNew() *JsonResp {
	var resp = &JsonResp{}
	resp.Code = SuccCode
	resp.Msg = SuccMsg
	return resp
}

func FailNew(j *JsonResp) {
	j.Code = FailCode
	j.Msg = FailMsg
}

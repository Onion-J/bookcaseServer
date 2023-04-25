package model

// Code2SessionResult 登录凭证校验后返回的JSON数据包模型
type Code2SessionResult struct {
	OpenId     string `json:"openId"`
	SessionKey string `json:"sessionkey"`
	UnionId    string `json:"unionId"`
	ErrCode    uint   `json:"errCode"`
	ErrMsg     string `json:"errMsg"`
}

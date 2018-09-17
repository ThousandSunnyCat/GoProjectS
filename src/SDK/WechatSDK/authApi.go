package WechatSDK

import (
	"encoding/json"
	"fmt"
)

type AppSessionResponse struct {
	Openid			string	`json:"openid"`
	Session_key		string	`json:"session_key"`
	Unionid			string	`json:"unionid"`
	Errcode			string	`json:"errcode"`
	ErrMsg			string	`json:"errMsg"`
}

func GetAppSessionKey(appId, secret, jsCode string) (*AppSessionResponse, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appId, secret, jsCode)
	res, err := httpGet(url)
	if err != nil {
		return nil, err
	}

	var mapRes AppSessionResponse
	json.Unmarshal(getBody(res), &mapRes)

	return &mapRes, nil
}
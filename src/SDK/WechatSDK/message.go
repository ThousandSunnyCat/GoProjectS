package WechatSDK

import (
	"encoding/json"
	"fmt"
)

var _baseUrl = "https://api.weixin.qq.com"	//api固定值

type MPTemplateMsgRequest struct {
	Touser			string					`json:"touser"`
	Template_id		string					`json:"template_id"`
	Url				string					`json:"url"`
	Miniprogram		*MiniProgramRequest		`json:"miniprogram"`
	Data			map[string]DataRequest	`json:"data"`
}

type MiniProgramRequest struct {
	Appid			string					`json:"appid"`
	Pagepath		string					`json:"pagepath"`
}

type DataRequest struct {
	Value			string					`json:"value"`
	Color			string					`json:"color"`
}

func SendAsync(request *MPTemplateMsgRequest, access_token string) (string, error) {
	url := fmt.Sprintf("%s/cgi-bin/message/template/send?access_token=%s", _baseUrl, access_token)
	j, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	res, err := httpPost(url, "application/json", string(j))
	if err != nil {
		return "", err
	}
	return string(getBody(res)), nil
}
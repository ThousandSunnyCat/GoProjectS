package WechatSDK

import (
	"strings"
	"fmt"
	"encoding/json"
	"time"
	"encoding/xml"
	"github.com/satori/go.uuid"
)

type UnifiedOrderBSPRequest struct {
	Appid				string  `xml:"appid" json:"appid"`
	Mch_id				string  `xml:"mch_id" json:"mch_id"`
	Openid				string  `xml:"openid" json:"openid"`
	Out_trade_no		string  `xml:"out_trade_no" json:"out_trade_no"`
	Total_fee			string  `xml:"total_fee" json:"total_fee"`
	Spbill_create_ip	string  `xml:"spbill_create_ip" json:"spbill_create_ip"`
	Notify_url			string  `xml:"notify_url" json:"notify_url"`
	Body				string	`xml:"body" json:"body"`

	Nonce_str			string	`xml:"nonce_str" json:"nonce_str"`
	Time_start			string	`xml:"time_start" json:"time_start"`
	Time_expire			string	`xml:"time_expire" json:"time_expire"`
	Trade_type			string	`xml:"trade_type" json:"trade_type"`
	Sign				string	`xml:"sign" json:"sign"`
}

type UnifiedOrderResponse struct {
	AppId				string
	TimeStamp			string
	NonceStr			string
	Package				string
	SignType			string
	PaySign				string
	PrepayId			string
}

var timeLayout = "20060102150405"	//golang时间模板固定值
var _mchbaseUrl = "https://api.mch.weixin.qq.com"	//api固定值
var _REPORTLEVENL = 1	// 上报等级：0: 不上报, 1: 仅上报失败, n: 全数上报

func UnifiedOrder(request *UnifiedOrderBSPRequest, signKey string) (res *UnifiedOrderResponse, err error) {

	defer func(){
		if p := recover(); p != nil {
			err = fmt.Errorf("UnifiedOrder: %v", p)
		}
	}()

	request.Nonce_str = uuid.Must(uuid.NewV4()).String()
	request.Time_start = time.Now().Format(timeLayout)
	request.Time_expire = time.Now().Add(30*time.Minute).Format(timeLayout)
	request.Trade_type = "JSAPI"

	// 下单请求
	maps := make(map[string]interface{})
	j, _ := json.Marshal(request)
	json.Unmarshal(j, &maps)
	request.Sign = makeMD5(maps, signKey)
	
	url := fmt.Sprintf("%v/pay/unifiedorder", _mchbaseUrl)
	start := time.Now()
	if res, err := httpPost(url, "text/xml", getXMLString(maps, "UnifiedOrderBSPRequest")); err != nil {
		return nil ,err
	} else {
		xml.Unmarshal(getBody(res), &maps)
		go reportCostTime(request.Appid, request.Mch_id, signKey, url, start, maps)
	}
	//下单成功

	// 验证返回参数
	if maps["prepay_id"] == nil || maps["prepay_id"] == "" {
		return nil, fmt.Errorf("prepay_id异常: %v", maps["return_msg"])
	}

	response := &UnifiedOrderResponse {
		AppId: request.Appid,
		TimeStamp: string(time.Now().Unix()),
		NonceStr: uuid.Must(uuid.NewV4()).String(),
		Package: fmt.Sprintf("prepay_id=%s", maps["prepay_id"]),
		SignType: "MD5",
		PrepayId: fmt.Sprintf("%v", maps["prepay_id"]),
	}

	j, _ = json.Marshal(response)
	json.Unmarshal(j, &maps)
	response.PaySign = makeMD5(maps, signKey)

	return response, nil
}

func NotifyOrder(request []byte, signKey string) (map[string]interface{}, error) {
	notifyData := make(map[string]interface{})
	if err := xml.Unmarshal(request, (*StringMap)(&notifyData)); err != nil {
		return nil, err
	}

	// 签名不通过
	if !CheckSign_MD5(notifyData, signKey) {
		return returnFail("签名错误")
	}
	// 支付失败
	if strings.ToUpper(fmt.Sprintf("%v", notifyData["return_code"])) != "SUCCESS" {
		return returnFail(notifyData["return_msg"])
	}
	// 若transaction_id不存在，则立即返回结果给微信支付后台
	if notifyData["transaction_id"] == nil {
		return returnFail("支付结果中微信订单号不存在")
	}

	// 处理回调
	appid := fmt.Sprintf("%v", notifyData["appid"])
	mch_id := fmt.Sprintf("%v", notifyData["mch_id"])
	transaction_id := fmt.Sprintf("%v", notifyData["transaction_id"])
	if !queryOrder(appid, mch_id, signKey, transaction_id) {
		return returnFail("订单查询失败")
	}

	return notifyData, nil;
}

// 
func makeMD5(maps map[string]interface{}, signKey string) string {
	if r, err := MakeSign_MD5(maps, signKey); err != nil {
		panic(err)
	} else {
		return r
	}
}

func getXMLString(request map[string]interface{}, reqType string) string {
	if xmlByteDate, err := xml.Marshal(request); err != nil {
		// error
		panic(err)
	} else {
		return strings.Replace(string(xmlByteDate), reqType,"xml",-1)
	}
}

func returnFail(msg interface{}) (map[string]interface{}, error) {
	resFail := make(map[string]interface{})
	resFail["return_code"] = "FAIL"
	resFail["return_msg"] = msg
	return resFail, nil
}

func queryOrder(appid, mch_id, signKey, transaction_id string) (isSuccess bool) {

	defer func(){
		if p := recover(); p != nil {
			isSuccess = false
		}
	}()

	reqmaps := make(map[string]interface{})
	reqmaps["appid"] = appid
	reqmaps["mch_id"] = mch_id
	reqmaps["transaction_id"] = transaction_id
	reqmaps["nonce_str"] = uuid.Must(uuid.NewV4()).String()
	reqmaps["sign"] = makeMD5(reqmaps, signKey)

	url := fmt.Sprintf("%s/pay/orderquery", _mchbaseUrl)
	maps := make(map[string]interface{})
	start := time.Now()
	if res, err := httpPost(url, "text/xml", getXMLString(StringMap(reqmaps), "StringMap")); err != nil {
		return false
	} else {
		xml.Unmarshal(getBody(res), &maps)
		go reportCostTime(appid, mch_id, signKey, url, start, maps)
	}

	if maps["return_code"] == "SUCCESS" && maps["result_code"] == "SUCCESS" {
		return true
	}

	return false
}

// 上报
func reportCostTime(appid, mch_id, signKey, interface_url string, start time.Time, resMaps map[string]interface{}) {
	if _REPORTLEVENL == 0 || (_REPORTLEVENL == 1 && resMaps["return_code"] == "SUCCESS" && resMaps["result_code"] == "SUCCESS") {
		return
	}

	timeCost := (int)(time.Now().Sub(start).Nanoseconds() / 1000000)

	reqmaps := make(map[string]interface{})
	reqmaps["interface_url"] = interface_url
	reqmaps["execute_time"] = timeCost

	if resMaps["return_code"] != nil {
		reqmaps["return_code"] = resMaps["return_code"]
	}
	if resMaps["return_msg"] != nil {
		reqmaps["return_msg"] = resMaps["return_msg"]
	}
	if resMaps["result_code"] != nil {
		reqmaps["result_code"] = resMaps["result_code"]
	}
	if resMaps["err_code"] != nil {
		reqmaps["err_code"] = resMaps["err_code"]
	}
	if resMaps["err_code_des"] != nil {
		reqmaps["err_code_des"] = resMaps["err_code_des"]
	}
	if resMaps["out_trade_no"] != nil {
		reqmaps["out_trade_no"] = resMaps["out_trade_no"]
	}
	if resMaps["device_info"] != nil {
		reqmaps["device_info"] = resMaps["device_info"]
	}

	defer func(){
		if p := recover(); p != nil {
			// log
		}
	}()

	reqmaps["appid"] = appid
	reqmaps["mch_id"] = mch_id
	reqmaps["time"] = time.Now().Format(timeLayout)
	reqmaps["nonce_str"] = uuid.Must(uuid.NewV4()).String()
	reqmaps["sign"] = makeMD5(reqmaps, signKey)

	url := fmt.Sprintf("%s/payitil/report", _mchbaseUrl)
	httpPost(url, "text/xml", getXMLString(StringMap(reqmaps), "StringMap"))
}

// 

type MicroPayRequest struct {
	// 必填字段
	Appid				string		`xml:"appid"`
	Mch_id				string		`xml:"mch_id"`
	Body				string		`xml:"body"`
	Out_trade_no		string		`xml:"out_trade_no"`
	Total_fee			int			`xml:"total_fee"`
	Spbill_create_ip	string		`xml:"spbill_create_ip"`
	Auth_code			string		`xml:"auth_code"`
	// 程序字段
	Nonce_str			string		`xml:"nonce_str"`
	Sign				string		`xml:"sign"`
	// 非必填字段
	Device_info			string		`xml:"device_info"`			//商户设备号
	Attach				string		`xml:"attach"`				//附加数据包，会原样返回
}

type MicroPayResponse struct {
	Openid				string		`xml:"openid"`
	Device_info			string		`xml:"device_info"`
	Trade_type			string		`xml:"trade_type"`
	Bank_type			string		`xml:"bank_type"`
	//Fee_type			string		`xml:"fee_type"`
	Total_fee			int			`xml:"total_fee"`
	Cash_fee			int			`xml:"cash_fee"`
	Transaction_id		string		`xml:"transaction_id"`
	Out_trade_no		string		`xml:"out_trade_no"`
	Attach				string		`xml:"attach"`
	Time_end			string		`xml:"time_end"`
}

func MicroPay(request *MicroPayRequest, signKey string) (response *MicroPayResponse, err error) {
	
	defer func(){
		if p := recover(); p != nil {
			err = fmt.Errorf("UnifiedOrder: %v", p)
		}
	}()

	url := fmt.Sprintf("%s/pay/micropay", _mchbaseUrl)
	request.Nonce_str = uuid.Must(uuid.NewV4()).String()

	// 下单请求
	maps := make(map[string]interface{})
	j, _ := json.Marshal(request)
	json.Unmarshal(j, &maps)
	request.Sign = makeMD5(maps, signKey)

	if res, err := httpPost(url, "text/xml", getXMLString(maps, "MicroPayRequest")); err != nil {
		return nil, err
	} else {
		xml.Unmarshal(getBody(res), &maps)
	}

	// 判断返回结果
	/*if res["return_code"] {

	}*/

	// 待续
	return nil, nil
}
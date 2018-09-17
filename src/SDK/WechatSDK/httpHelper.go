package WechatSDK

import (
	"crypto/tls"
	"regexp"
	"io/ioutil"
	"errors"
	"strings"
	"net/http"
	"log"
)

func HttpSend(url, method, params string, success func(body string), fair func(err error)) {
	defer func(){
		if err := recover(); err != nil {
			log.Printf("run time panic: %v", err)
		}
	}()

	//
	var resp *http.Response
	var err error

	switch strings.ToUpper(method) {
	case http.MethodGet:
		// 参数需要转换
		resp, err = httpGet(url)
	case http.MethodPost:
		resp, err = httpPost(url, "application/json; charset=UTF-8", params)
	default:
		err = errors.New("method is not found")
	}

	if err != nil {
        fair(err)
    }
  
    success(string(getBody(resp)))
}

func getBody(resp *http.Response) []byte {
	defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic("解码失败")
	}
	return body
}

var regHttp = regexp.MustCompile(`^[\S]+://`)

var localClient *http.Client

func getClient() (client *http.Client) {
	if localClient == nil {		
		tr := &http.Transport {
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		localClient = &http.Client{Transport: tr}
	}

	return localClient
}

func httpGet(url string) (resp *http.Response, err error) {
	switch strings.ToLower(regHttp.FindString(url)) {
	case "http":
		return http.Get(url)
	case "https":
		return getClient().Get(url)
	default:
		return nil, errors.New("Unrecognized the protocol")
	}
}

func httpPost(url string, contentType string, params string) (resp *http.Response, err error) {
	switch strings.ToLower(regHttp.FindString(url)) {
	case "http":
		return http.Post(url, contentType, strings.NewReader(params))
	case "https":
		return getClient().Post(url, contentType, strings.NewReader(params))
	default:
		return nil, errors.New("Unrecognized the protocol")
	}
}
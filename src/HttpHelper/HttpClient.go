package HttpHelper

import (
	"io/ioutil"
	"errors"
	"strings"
	"net/http"
	"log"
)

func Send(url string, method string, params string, success func(body string), fair func(err error)) {
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
 
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic("解码失败")
    }
 
    success(string(body))
}
package WechatSDK

import (
	"encoding/hex"
	"strings"
	"crypto/md5"
	"sort"
	"fmt"
)

func MakeSign_MD5(v map[string]interface{}, signKey string) (s string, e error) {
	defer func(){
		if err := recover(); err != nil {
			s = ""
			e = fmt.Errorf("MakeSign Panic: %v", err)
		}
	}()
	
	res := md5.Sum(toUrl(v, signKey))
	return hex.EncodeToString(res[:]), nil
}

func CheckSign_MD5(v map[string]interface{}, signKey string) bool {
	if v["Sign"] == nil || v["Sign"] == "" {
		// error
		return false
	}

	if s, e := MakeSign_MD5(v, signKey); e != nil {
		// error
		return false
	} else if s != v["Sign"] {
		// 验证不通过
		return false
	}

	return true
}

func toUrl(maps map[string]interface{}, signKey string) []byte {
	
	var keys []string
	for k := range maps {
		keys = append(keys, strings.ToLower(k))
	}
	sort.Strings(keys)

	var buff string
	for _, k := range keys {
		if k == "sign" {
			continue
		} 
		buff = fmt.Sprintf("%s%s=%v&", buff, k, maps[k])
	}

	return []byte(buff + "key=" + signKey)
}
package HttpHelper

import (
	"strings"
	"net/http"

)

func httpPost(url string, contentType string, params string) (resp *http.Response, err error) {
    return http.Post(url, contentType, strings.NewReader(params))
}
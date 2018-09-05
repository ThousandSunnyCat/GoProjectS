package HttpHelper

import (
	"net/http"
)

func httpGet(url string) (resp *http.Response, err error) {
    return http.Get(url)
}
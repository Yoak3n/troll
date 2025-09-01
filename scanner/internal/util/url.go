package util

import (
	"fmt"
	"strings"
)

func AppendParamsToUrl(url string, params map[string]string) string {
	url = url + "?"
	list := make([]string, 0)
	for k, v := range params {
		p := fmt.Sprintf("%s=%s", k, v)
		list = append(list, p)
	}

	return url + strings.Join(list, "&")
}

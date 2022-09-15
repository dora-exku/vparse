package utils

import (
	"net/http"
	"strings"
)

func GetCk(cks []*http.Cookie, name string) string {
	for _, item := range cks {
		if item.Name == name {
			return item.Value
		}
	}
	return ""
}

func SplitCks(cks string) (cookies []*http.Cookie) {

	arr := strings.Split(cks, ";")
	for _, item := range arr {
		tmpArr := strings.Split(strings.TrimSpace(item), "=")
		if len(tmpArr) >= 2 {
			cookies = append(cookies, &http.Cookie{
				Name:  tmpArr[0],
				Value: strings.Join(tmpArr[1:], ""),
			})
		}
	}

	return cookies
}

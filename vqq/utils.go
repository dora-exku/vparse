package vqq

import (
	"net/http"
	"strconv"
	"strings"
)

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

func time33(v string) string {
	i := 5381
	for _, item := range v {
		i += (i << 5) + int(item)
	}
	return strconv.Itoa(2147483647 & i)
}

func GetCk(cks []*http.Cookie, name string) string {
	for _, item := range cks {
		if item.Name == name {
			return item.Value
		}
	}
	return ""
}

func getVid(url string) string {

	hasQuestionMark := strings.Index(url, "?")
	if hasQuestionMark > 0 {
		url = url[0:hasQuestionMark]
	}
	url = strings.TrimSuffix(url, ".html")
	return url[strings.LastIndex(url, "/")+1:]
}

package vparse

import (
	"net/http"
	"net/url"
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
				Value: url.QueryEscape(strings.Trim(strings.Join(tmpArr[1:], ""), "\"")),
			})
		}
	}

	return cookies
}

func GetCk(cks []*http.Cookie, name string) string {
	for _, item := range cks {
		if item.Name == name {
			return item.Value
		}
	}
	return ""
}

func time33(v string) string {
	i := 5381
	for _, item := range v {
		i += (i << 5) + int(item)
	}
	return strconv.Itoa(2147483647 & i)
}

//func getVid(url string) string {
//	//https://m.v.qq.com/play.html?vid=k0025c8k9hr&cid=9p15mebx5gn4pz4
//	if strings.Contains(url, "m.v.qq.com/play.html") {
//		v, err := u.Parse(url)
//		if err != nil {
//			return ""
//		}
//		return v.Query().Get("vid")
//	}
//
//	hasQuestionMark := strings.Index(url, "?")
//	if hasQuestionMark > 0 {
//		url = url[0:hasQuestionMark]
//	}
//	url = strings.TrimSuffix(url, ".html")
//	return url[strings.LastIndex(url, "/")+1:]
//}

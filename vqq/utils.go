package vqq

import (
	"strconv"
	"strings"
)

func time33(v string) string {
	i := 5381
	for _, item := range v {
		i += (i << 5) + int(item)
	}
	return strconv.Itoa(2147483647 & i)
}


func getVid(url string) string {

	hasQuestionMark := strings.Index(url, "?")
	if hasQuestionMark > 0 {
		url = url[0:hasQuestionMark]
	}
	url = strings.TrimSuffix(url, ".html")
	return url[strings.LastIndex(url, "/")+1:]
}

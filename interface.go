package vparse

import "net/http"

type Parse interface {
	SetCall(name string, call CallFunc)
	SetCookies(cookies []*http.Cookie)
	Parse(url, definition string) (string, error)
	SetUserAgent(string)
}

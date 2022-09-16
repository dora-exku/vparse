package vparse

import "net/http"

type Parse interface {
	WithCall(name string, call CallFunc)
	WithCookies(cookies []*http.Cookie)
	Parse(url, definition string) (string, error)
	WithUserAgent(string)
}

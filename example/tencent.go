package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dora-exku/vparse"
	"github.com/go-resty/resty/v2"
	"os"
)

func main() {
	c, err := os.ReadFile("cookie_tencent.ck")
	if err != nil {
		fmt.Println(err)
		return
	}

	cks := vparse.SplitCks(string(c))

	video := vparse.New("tencent")
	video.SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36")
	video.SetCookies(cks)
	video.SetCall("ckey", func(args ...any) (string, error) {
		// url vid guid tm
		if len(args) != 4 {
			return "", errors.New("ckey params invalid")
		}
		// 获取ckey
		client := resty.New()
		ckeyResp, _ := client.R().SetQueryParams(map[string]string{
			"vid":      args[1].(string),
			"tm":       args[3].(string),
			"guid":     args[2].(string),
			"version":  "3.5.57",
			"platform": "10901",
			"url":      args[0].(string),
			"referer":  args[0].(string),
		}).Get("http://localhost:5050/tencent/ckey81")

		var result struct {
			Ckey string `json:"ckey"`
		}
		err = json.Unmarshal(ckeyResp.Body(), &result)
		if err != nil {
			return "", err
		}
		return result.Ckey, nil
	})

	m3u8, err := video.Parse(
		"https://v.qq.com/x/cover/mzc002000xg1sad/t0044iedn5j.html",
		"fhd",
	)
	fmt.Println(m3u8, err)
}

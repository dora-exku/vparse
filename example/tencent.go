package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dora-exku/vparse"
	"github.com/dora-exku/vparse/utils"
	"github.com/go-resty/resty/v2"
	"os"
)

func main() {
	c, err := os.ReadFile("v.ck")
	if err != nil {
		fmt.Println(err)
		return
	}

	cks := utils.SplitCks(string(c))

	v := vparse.TencentParse{
		UA: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
	}
	v.WithCookies(cks)
	v.WithCall("ckey", func(args ...any) (string, error) {
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

		var result struct{
			Ckey string `json:"ckey"`
		}
		err = json.Unmarshal(ckeyResp.Body(), &result)
		if err != nil {
			return "", err
		}
		return result.Ckey, nil
	})

	m3u8, err := v.Parse(
		"https://m.v.qq.com/play.html?vid=k0025c8k9hr&cid=9p15mebx5gn4pz4",
		"fhd",
	)
	fmt.Println(m3u8, err)
}

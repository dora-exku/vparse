package main

import (
	"errors"
	"fmt"
	"github.com/dora-exku/vparse"
	"github.com/go-resty/resty/v2"
	"os"
)

func main() {

	c, err := os.ReadFile("cookie_iqiyi.ck")
	if err != nil {
		fmt.Println(err)
		return
	}

	cks := vparse.SplitCks(string(c))

	video := vparse.New("iqiyi")

	video.SetUserAgent("Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1")
	video.SetCookies(cks)
	video.SetCall("authkey", func(args ...any) (string, error) {
		if len(args) != 2 {
			return "", errors.New("params invalid")
		}

		client := resty.New()
		resp, err := client.SetQueryParam("tm", args[0].(string)).SetQueryParam("vid", args[1].(string)).R().Get("http://127.0.0.1:5050/iqiyi/authkey")

		if err != nil {
			fmt.Println(err)
			return "", err
		}

		k := resp.Body()

		return string(k[12 : len(k)-2]), nil
	})

	video.SetCall("vf", func(args ...any) (string, error) {
		if len(args) != 1 {
			return "", errors.New("params invalid")
		}

		client := resty.New()
		resp, err := client.SetQueryParam("url", args[0].(string)).R().Get("http://127.0.0.1:5050/iqiyi/cmd5x")
		if err != nil {
			return "", err
		}

		k := resp.Body()

		return string(k[7 : len(k)-2]), nil
	})

	m3u8, err := video.Parse(
		//"https://www.iqiyi.com/v_19rrlzcmcg.html",
		"https://www.iqiyi.com/v_19rrlzcmcg.html",
		"500",
	)

	fmt.Println(m3u8, err)
}

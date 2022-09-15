package vqq

import (
	"encoding/json"
	"fmt"
	"github.com/dora-exku/v-analysis/utils"
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
)

type TokenInfo struct {
	UserID      int    `json:"vuserid"`
	Session     string `json:"vusession"`
	AccessToken string `json:"access_token"`
}

func AuthRefresh(cks string) (tokenInfo TokenInfo, ncks string, err error) {
	client := resty.New()
	//var cookies []*http.Cookie
	cookies := utils.SplitCks(cks)

	timeM := strconv.FormatInt(time.Now().UnixMilli(), 10)

	callBack := "jQuery19109216653952017793_" + timeM

	resp, err := client.SetHeaders(map[string]string{
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
		"referer":    "https://v.qq.com/",
	}).SetCookies(cookies).SetQueryParams(map[string]string{
		"vappid":   "11059694",
		"vsecret":  "fdf61a6be0aad57132bc5cdf78ac30145b6cd2c1470b0cfe",
		"type":     "qq",
		"g_tk":     "",
		"g_vstk":   time33(utils.GetCk(cookies, "vqq_vusession")),
		"g_actk":   time33(utils.GetCk(cookies, "vqq_access_token")),
		"callback": callBack,
		"_":        timeM,
	}).R().Get("https://access.video.qq.com/user/auth_refresh")

	fmt.Println(resp.String())

	if err != nil {
		return TokenInfo{}, "", err
	}

	// 合并cookie
	for _, item := range cookies {
		var n = utils.GetCk(resp.Cookies(), item.Name)
		if n != item.Value && n != "" {
			item.Value = n
		}
	}

	for _, item := range cookies {
		ncks += fmt.Sprintf("%s=%s;", item.Name, item.Value)
	}

	body := resp.Body()

	res := body[len(callBack)+1 : len(body)-2]

	var data TokenInfo

	err = json.Unmarshal(res, &data)
	if err != nil {
		return TokenInfo{}, "", err
	}

	return data, ncks, nil
}

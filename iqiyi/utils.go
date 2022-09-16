package iqiyi

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strconv"
	"strings"
)

func GetVid(url string, cks []*http.Cookie) (tvid, vid string, bid int) {
	client := resty.New()

	resp, err := client.SetHeaders(map[string]string{
		"referer":    "https://www.iqiyi.com/",
		"user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
	}).SetCookies(cks).R().Get(url)

	if err != nil {
		return "", "", 0
	}


	data := resp.String()

	start := strings.Index(data, "window.QiyiPlayerProphetData=")
	if start == -1 {
		return "", "", 0
	}

	data = data[start+29:]
	data = data[:strings.Index(data, "</script>")]

	var result struct {
		Tvid int64 `json:"tvid"`
		V    struct {
			Vidl []struct {
				Bid int    `json:"bid"`
				Vid string `json:"vid"`
			} `json:"vidl"`
		} `json:"v"`
	}

	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		return "", "", 0
	}

	if len(result.V.Vidl) == 0 {
		return "", "", 0
	}

	for _, item := range result.V.Vidl {
		if item.Bid > bid {
			bid = item.Bid
			vid = item.Vid
		}
	}

	return strconv.FormatInt(result.Tvid, 10), vid, bid
}

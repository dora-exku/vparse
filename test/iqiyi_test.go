package test

import (
	"fmt"
	iqiyi2 "github.com/dora-exku/v-analysis/iqiyi"
	"github.com/dora-exku/v-analysis/utils"
	"github.com/go-resty/resty/v2"
	"os"
	"testing"
)

func TestIqiyiAnalysis(t *testing.T) {
	c, err := os.ReadFile("../iqiyi.ck")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	cks := utils.SplitCks(string(c))

	i := iqiyi2.New()
	m3u8 := i.WithCookie(cks).Analysis(
		"https://www.iqiyi.com/v_22bfixztj38.html?vfrm=pcw_home&vfrmblk=712211_dianshiju&vfrmrst=712211_dianshiju_float_video_area2",
		"",
		func(tm, vid string) string {

			client := resty.New()
			resp, err := client.SetQueryParam("tm", tm).SetQueryParam("vid", vid).R().Get("http://127.0.0.1:5050/iqiyi/authkey")

			if err != nil {
				fmt.Println(err)
				return ""
			}

			k := resp.Body()

			return string(k[12 : len(k)-2])
		},
		func(url string) string {
			client := resty.New()
			resp, err := client.SetQueryParam("url", url).R().Get("http://127.0.0.1:5050/iqiyi/cmd5x")
			if err != nil {
				fmt.Println("cmd5x", err)
				return ""
			}

			k := resp.Body()

			return string(k[7 : len(k)-2])
		},
	)

	if m3u8 == "" {
		t.FailNow()
	}
}

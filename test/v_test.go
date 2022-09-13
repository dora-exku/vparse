package test

import (
	"encoding/json"
	"github.com/dora-exku/v-analysis/vqq"
	"github.com/go-resty/resty/v2"
	"os"
	"testing"
)

func TestAnalysis(t *testing.T) {

	cks, err := os.ReadFile("../v.ck")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	_, ncks, err := vqq.AuthRefresh(string(cks))

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	var ncksArr = vqq.SplitCks(ncks)
	url := "https://v.qq.com/x/cover/324olz7ilvo2j5f/t0035aw2v35.html"

	v := vqq.New()

	m3url := v.WithCookie(ncksArr).Analysis(url, "fhd", func(vqq vqq.Vqq, url, vid, timeStr string) string {
		// 获取ckey
		client := resty.New()
		ckeyResp, _ := client.R().SetQueryParams(map[string]string{
			"vid":      vid,
			"tm":       timeStr,
			"guid":     v.Guid,
			"version":  "3.5.57",
			"platform": "10901",
			"url":      url,
			"referer":  url,
		}).Get("http://localhost:5050/tencent/ckey81")

		var ckeyBody map[string]string
		json.Unmarshal(ckeyResp.Body(), &ckeyBody)

		ckey := ckeyBody["ckey"]
		return ckey
	})

	if m3url == "" {
		t.FailNow()
	}
}

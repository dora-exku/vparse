package test

import (
	"encoding/json"
	"github.com/dora-exku/v-analysis/vqq"
	"github.com/go-resty/resty/v2"
	"os"
	"testing"
)

func TestQqAnalysis(t *testing.T) {

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
	//url := "https://v.qq.com/x/cover/324olz7ilvo2j5f/t0035aw2v35.html"
	//url := "https://v.qq.com/x/cover/mzc00200lojsjys/s00446869zn.html"
	url := "http://m.v.qq.com/x/cover/x/mzc002000ry9s13/p0044gtj3t2.html?&url_from=share&second_share=0&share_from=copy"

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

	t.Log(m3url)
	t.FailNow()

	//if m3url == "" {
	//	t.FailNow()
	//}
}

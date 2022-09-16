package iqiyi

import (
	"encoding/json"
	"errors"
	"github.com/dora-exku/v-analysis/logger"
	"github.com/dora-exku/v-analysis/utils"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Viqiyi struct {
	Cookies []*http.Cookie
}

func New() *Viqiyi {
	return &Viqiyi{}
}

func (v *Viqiyi) WithCookie(cookie []*http.Cookie) *Viqiyi {
	v.Cookies = cookie
	return v
}

func (v Viqiyi) Analysis(url, bid string, authKeyCall func(tm, vid string) string, vfCall func(url string) string) string {

	rtvid, rvid, rbid := GetVid(url, v.Cookies)

	if rtvid == "" {
		return ""
	}

	if bid == "" {
		if rbid > 500 {
			bid = "500"
		} else {
			bid = strconv.Itoa(rbid)
		}
	}

	timeN := time.Now().UnixMilli()

	client := resty.New()
	resp, err := client.SetHeaders(map[string]string{
		//"origin":     "https://m.iqiyi.com",
		"referer":    "https://m.iqiyi.com/",
		"user-agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1",
	}).SetQueryParams(map[string]string{
		"tvid": rtvid,
		//"tvid":          "120250200",
		"bid": bid,
		"vid": rvid,
		//"vid":           "34669bbc2271499db22fdb80459f0dda",
		"src":           "02020031010000000000",
		"vt":            "0",
		"rs":            "0",
		"uid":           utils.GetCk(v.Cookies, "P00010"),
		"ori":           "h5",
		"ps":            "1",
		"k_uid":         utils.GetCk(v.Cookies, "QC005"),
		"pt":            "0",
		"d":             "0",
		"s":             "",
		"lid":           "",
		"cf":            "",
		"ct":            "",
		"authKey":       authKeyCall(strconv.FormatInt(timeN, 10), rtvid),
		"k_tag":         "1",
		"dfp":           "e1a6a88afc165c4c8481f9ac08eb2b552afa87a15f4880589b713fa6b574b4955f",
		"locale":        "zh_cn",
		"prio":          `{"ff":"m3u8","code":-1}`,
		"pck":           "6cF0krJFUMn3OCPm2WA45Bvm25eeI5lAbObz5lpvih2d7eVMYYSOqsa35kyHowK3WwAZ5f",
		"k_err_retries": "0",
		"up":            "",
		"qd_v":          "2",
		"tm":            strconv.FormatInt(timeN, 10),
		"qdy":           "i",
		"qds":           "0",
		"k_ft1":         "755914244096",
		//"k_ft4":         "1161084347621380",
		"k_ft5":    "262145",
		"bop":      `{"version":"10.0","dfp":"e1a6a88afc165c4c8481f9ac08eb2b552afa87a15f4880589b713fa6b574b4955f"}`,
		"callback": "Qb3626d3edb4502deea40c76a1b4839c9",
		"ut":       "1",
	}).OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {

		vf := vfCall("/jp/dash?" + client.QueryParam.Encode())
		if vf == "" {
			return errors.New("vf is empty")
		}
		client.SetQueryParam("vf", vf)
		return nil
	}).R().Get("https://cache.video.iqiyi.com/jp/dash")

	body := resp.String()

	//fmt.Println(body)

	if strings.HasPrefix(body, "try") {
		body = strings.TrimSuffix(body, "\n);}catch(e){};")
		body = strings.TrimPrefix(body, "try{Qb3626d3edb4502deea40c76a1b4839c9(")
	}

	//fmt.Println(body)

	if err != nil {
		logger.Error(err)
		return ""
	}

	var result struct {
		Data struct {
			Program struct {
				Videos []struct {
					M3u8Url string `json:"m3u8Url"`
					M3u8    string `json:"m3u8"`
					Bid     int    `json:"bid"`
					Url     string `json:"url"`
				} `json:"video"`
			} `json:"program"`
			DD string `json:"dd"`
		} `json:"data"`
		Code string `json:"code"`
	}

	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		logger.Error(err)
		return ""
	}

	if result.Code != "A00000" {
		logger.Error(errors.New("Code:" + result.Code))
		return ""
	}

	for _, item := range result.Data.Program.Videos {
		if strconv.Itoa(item.Bid) == bid {
			return item.Url
		}
	}
	return ""
}

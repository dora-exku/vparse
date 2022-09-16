package vparse

import (
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type IqiyiParse struct {
	Cookies     []*http.Cookie
	UA          string
	callFuncMap map[string]CallFunc
}

func (parse *IqiyiParse) WithCall(name string, call CallFunc) {
	if parse.callFuncMap == nil {
		parse.callFuncMap = make(map[string]CallFunc)
	}
	parse.callFuncMap[name] = call
}

func (parse *IqiyiParse) WithCookies(cookies []*http.Cookie) {
	parse.Cookies = cookies
}

func (parse *IqiyiParse) WithUserAgent(ua string) {
	parse.UA = ua
}

func (parse *IqiyiParse) Parse(url, definition string) (m3u8 string, err error) {
	rtvid, rvid, rbid, err := parse.getVid(url)
	if err != nil {
		return "", err
	}

	if rtvid == "" {
		return "", errors.New("tvid invalid")
	}

	if definition == "" {
		if rbid > 500 {
			definition = "500"
		} else {
			definition = strconv.Itoa(rbid)
		}
	}

	timeN := time.Now().UnixMilli()

	// authKey
	if parse.callFuncMap == nil {
		return "", errors.New("call invalid")
	}
	authKeyCall, ok := parse.callFuncMap["authkey"]
	if !ok {
		return "", errors.New("auth key call invalid")
	}

	authKey, err := authKeyCall(strconv.FormatInt(timeN, 10), rtvid)
	if err != nil {
		return "", err
	}

	// vf call
	vfCall, ok := parse.callFuncMap["vf"]
	if !ok {
		return "", errors.New("vf call invalid")
	}

	client := resty.New()

	resp, err := client.SetHeaders(map[string]string{
		"referer":    "https://m.iqiyi.com/",
		"user-agent": parse.UA,
	}).SetQueryParams(map[string]string{
		"tvid":          rtvid,
		"bid":           definition,
		"vid":           rvid,
		"src":           "02020031010000000000",
		"vt":            "0",
		"rs":            "0",
		"uid":           GetCk(parse.Cookies, "P00010"),
		"ori":           "h5",
		"ps":            "1",
		"k_uid":         GetCk(parse.Cookies, "QC005"),
		"pt":            "0",
		"d":             "0",
		"s":             "",
		"lid":           "",
		"cf":            "",
		"ct":            "",
		"authKey":       authKey,
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
		vf, err := vfCall("/jp/dash?" + client.QueryParam.Encode())
		if err != nil {
			return err
		}
		client.SetQueryParam("vf", vf)
		return nil
	}).R().Get("https://cache.video.iqiyi.com/jp/dash")

	body := resp.String()

	if strings.HasPrefix(body, "try") {
		body = strings.TrimSuffix(body, "\n);}catch(e){};")
		body = strings.TrimPrefix(body, "try{Qb3626d3edb4502deea40c76a1b4839c9(")
	}

	if err != nil {
		return "", err
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
		return "", err
	}

	if result.Code != "A00000" {
		return "", errors.New("error : " + body)
	}

	for _, item := range result.Data.Program.Videos {
		if item.Url != "" {
			return item.Url, nil
		}
	}

	return "", errors.New(body)
}

func (parse IqiyiParse) getVid(url string) (tvid, vid string, bid int, err error) {
	client := resty.New()

	resp, err := client.SetHeaders(map[string]string{
		"referer":    "https://www.iqiyi.com/",
		"user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/105.0.0.0 Safari/537.36",
	}).SetCookies(parse.Cookies).R().Get(url)

	if err != nil {
		return "", "", 0, err
	}

	data := resp.String()

	start := strings.Index(data, "window.QiyiPlayerProphetData=")
	if start == -1 {
		return "", "", 0, errors.New("resp invalid:" + data)
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
		return "", "", 0, err
	}

	if len(result.V.Vidl) == 0 {
		return "", "", 0, err
	}

	for _, item := range result.V.Vidl {
		if item.Bid > bid {
			bid = item.Bid
			vid = item.Vid
		}
	}

	return strconv.FormatInt(result.Tvid, 10), vid, bid, nil
}

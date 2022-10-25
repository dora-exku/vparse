package vparse

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"net/http"
	netUrl "net/url"
	"strconv"
	"strings"
	"time"
)

type TencentParse struct {
	accessToken string
	appID       string
	session     string
	openID      string
	userID      string
	guid        string
	mainLogin   string
	Cookies     []*http.Cookie
	UA          string
	callFuncMap map[string]CallFunc
}

func (parse *TencentParse) SetCall(name string, call CallFunc) {
	if parse.callFuncMap == nil {
		parse.callFuncMap = make(map[string]CallFunc)
	}
	parse.callFuncMap[name] = call

}

func (parse *TencentParse) SetCookies(cookies []*http.Cookie) {
	parse.Cookies = cookies
	parse.appID = GetCk(cookies, "vqq_appid")
	parse.openID = GetCk(cookies, "vqq_openid")
	parse.guid = GetCk(cookies, "video_guid")
	parse.accessToken = GetCk(cookies, "vqq_access_token")
	parse.session = GetCk(cookies, "vqq_vusession")
	parse.userID = GetCk(cookies, "vqq_vuserid")
	parse.mainLogin = GetCk(cookies, "main_login")
}

func (parse *TencentParse) SetUserAgent(ua string) {
	parse.UA = ua
}

func (parse *TencentParse) Parse(url, definition string) (m3u8 string, err error) {

	err = parse.authRefresh()
	if err != nil {
		// 刷新失败
		return "", err
	}

	var client *resty.Client
	var vid string
	var timeStr = strconv.FormatInt(time.Now().Unix(), 10)

	vid = parse.getVid(url)

	md5h := md5.New()

	md5h.Write([]byte(time.Now().String()))

	flowid := hex.EncodeToString(md5h.Sum(nil))

	if parse.callFuncMap == nil {
		return "", errors.New("ckey call is nil")
	}
	ckeyCall, ok := parse.callFuncMap["ckey"]
	if !ok {
		return "", errors.New("ckey call is nil")
	}

	ckey, err := ckeyCall(url, vid, parse.guid, timeStr)
	if err != nil {
		return "", err
	}
	
	client = resty.New()

	loginaccess, _ := json.Marshal(map[string]string{
		"access_token": parse.accessToken,
		"appid":        parse.appID,
		"vusession":    parse.session,
		"openid":       parse.openID,
		"vuserid":      parse.userID,
		"video_guid":   parse.guid,
		"main_login":   parse.mainLogin,
	})

	call := "getinfo_callback_" + strconv.Itoa(rand.Intn(999999-100000)+100000)
	resp, _ := client.SetHeaders(map[string]string{
		"user-agent": parse.UA,
	}).R().SetQueryParams(map[string]string{
		"charge":     "0",
		"otype":      "json",
		"defnpayver": "3",
		"spau":       "1",
		"spaudio":    "0",
		"spwm":       "1",
		"sphls":      "2",
		"host":       "v.qq.com",
		"refer":      url,
		"ehost":      url,
		"sphttps":    "1",
		"encryptVer": "8.1",
		"cKey":       ckey,
		"clip":       "4",
		"guid":       parse.guid,
		"flowid":     flowid,
		"platform":   "10901",
		"sdtfrom":    "v1010",
		"appVer":     "3.5.57",
		"unid":       "",
		"auth_from":  "",
		"auth_ext":   "",
		"vid":        vid,
		"defn":       definition, // 清晰度
		"fhdswitch":  "0",
		"dtype":      "3",
		"spsrt":      "2",
		"tm":         timeStr,
		"lang_code":  "0",
		"logintoken": string(loginaccess),
		"spvvpay":    "1",
		"spadseg":    "3",
		"hevclv":     "0",
		"spsfrhdr":   "0",
		"spvideo":    "0",
		"drm":        "40",
		"callback":   call,
	}).Get("https://h5vv6.video.qq.com/getvinfo")

	var data struct {
		Vl struct {
			Vi []struct {
				Ul struct {
					Ui []struct {
						Dt  int    `json:"dt"`
						Dtc int    `json:"dtc"`
						URL string `json:"url"`
					} `json:"ui"`
				} `json:"ul"`
			} `json:"vi"`
		} `json:"vl"`
	}



	var s = strings.TrimPrefix(resp.String(), call+"(")
	s = strings.TrimSuffix(s, ")")


	err = json.Unmarshal([]byte(s), &data)
	if err != nil {
		return "", err
	}

	//var ndata items
	//ndata = data.(items)
	if len(data.Vl.Vi) == 0 {

		return "", errors.New("Videos is empty" + s)
	}

	for _, item := range data.Vl.Vi[0].Ul.Ui {
		return item.URL, nil
	}

	return "", nil
}

func (parse *TencentParse) authRefresh() error {
	client := resty.New()
	timeM := strconv.FormatInt(time.Now().UnixMilli(), 10)
	callback := fmt.Sprintf("jQuery19109216653952017793_%s", timeM)

	resp, err := client.SetHeaders(map[string]string{
		"User-Agent": parse.UA,
		"referer":    "https://v.qq.com/",
	}).SetCookies(parse.Cookies).SetQueryParams(map[string]string{
		"vappid":   "11059694",
		"vsecret":  "fdf61a6be0aad57132bc5cdf78ac30145b6cd2c1470b0cfe",
		"type":     "qq",
		"g_tk":     "",
		"g_vstk":   time33(parse.session),
		"g_actk":   time33(parse.accessToken),
		"callback": callback,
		"_":        timeM,
	}).R().Get("https://access.video.qq.com/user/auth_refresh")
	if err != nil {
		// 发生错误
		return err
	}
	if resp.StatusCode() != 200 {
		// 刷新是返回了非200
		return errors.New("刷新Token失败")
	}
	// 处理返回值
	data := resp.String()
	data = strings.TrimPrefix(data, callback+"(")
	data = strings.TrimSuffix(data, ");")

	var result struct {
		ErrCode     int    `json:"errcode"`
		Session     string `json:"vusession"`
		AccessToken string `json:"access_token"`
	}


	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		// 解析失败
		return err
	}

	parse.accessToken = result.AccessToken
	parse.session = result.Session
	// 处理 token
	for _, item := range parse.Cookies {
		var n = GetCk(resp.Cookies(), item.Name)
		if n != item.Value && n != "" {
			item.Value = n
		}
	}

	return nil
}

func (parse TencentParse) getVid(url string) string {
	//https://m.v.qq.com/play.html?vid=k0025c8k9hr&cid=9p15mebx5gn4pz4 pc客户端分享
	//http://m.v.qq.com/x/cover/x/mzc00200fr1ry1o/m00441h6knj.html // 手机客户端分享
	//https://v.qq.com/x/cover/mzc00200fr1ry1o/m00441h6knj.html // web端url
	if strings.Contains(url, "m.v.qq.com/play.html") {
		v, err := netUrl.Parse(url)
		if err != nil {
			return ""
		}
		return v.Query().Get("vid")
	}

	hasQuestionMark := strings.Index(url, "?")
	if hasQuestionMark > 0 {
		url = url[0:hasQuestionMark]
	}
	url = strings.TrimSuffix(url, ".html")
	return url[strings.LastIndex(url, "/")+1:]
}

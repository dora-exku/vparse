package vqq

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/dora-exku/v-analysis/logger"
	"github.com/go-resty/resty/v2"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Vqq struct {
	AccessToken string
	AppID       string
	Session     string
	OpenID      string
	UserID      string
	Guid        string
	MainLogin   string
}

func New() *Vqq {
	return &Vqq{
		AccessToken: "",
		AppID:       "",
		Session:     "",
		OpenID:      "",
		UserID:      "",
		Guid:        "",
		MainLogin:   "",
	}
}

func (v *Vqq) WithCookie(ncks []*http.Cookie) *Vqq {
	v.AppID = GetCk(ncks, "vqq_appid")
	v.OpenID = GetCk(ncks, "vqq_openid")
	v.Guid = GetCk(ncks, "video_guid")
	v.AccessToken = GetCk(ncks, "vqq_access_token")
	v.Session = GetCk(ncks, "vqq_vusession")
	v.UserID = GetCk(ncks, "vqq_vuserid")

	return v
}

func (v *Vqq) WithPlatform(pf string) *Vqq {
	v.MainLogin = pf
	return v
}

// Analysis 解析 url pc端的地址  defn 清晰度 已知 sd：270P  hd：480P  shd：720P   fhd：1080P
func (v Vqq) Analysis(url string, defn string, ckcall func(vqq Vqq, url, vid, timeStr string) string) string {

	var client *resty.Client
	var vid string
	var timeStr = strconv.FormatInt(time.Now().Unix(), 10)

	vid = getVid(url)

	md5h := md5.New()

	md5h.Write([]byte(time.Now().String()))

	flowid := hex.EncodeToString(md5h.Sum(nil))

	client = resty.New()

	ckey := ckcall(v, url, vid, timeStr)

	client = resty.New()

	loginaccess, _ := json.Marshal(map[string]string{
		"access_token": v.AccessToken,
		"appid":        v.AppID,
		"vusession":    v.Session,
		"openid":       v.OpenID,
		"vuserid":      v.UserID,
		"video_guid":   v.Guid,
		"main_login":   v.MainLogin,
	})

	//?charge=0&otype=json&defnpayver=3&spau=1&spaudio=0&spwm=1&sphls=2&host=v.qq.com&refer=https://v.qq.com/x/cover/324olz7ilvo2j5f/v0044u7dz79.html&ehost=https://v.qq.com/x/cover/324olz7ilvo2j5f/v0044u7dz79.html&sphttps=1&encryptVer=8.1&cKey=229AAAF1FDEA5CF00ED2A945C0CFA0B21DE7677D5877D9652FC1D58427855DBAD5057CCAE859A4D719D4A25A4F14401E177A8BDF2F8C50361B38E62292878CF3B120510782BA6D5C018E5C6719B096D9EFDE7FA0E2A7B83BD872DD0A24F53A40D9454D4032543BF1C68FE54A8EAB17DFA96DDAADEA6B9F9B157BA531A5CA57E73620E87085F0C556F2C7F42153A45E8394ABACB926657E412F89AFED2AC67DDF916BEEA5C8466E105EBA80B61CDBA3072B34940F4C92A50F57476EFA9B4FBE389E6CB1B5B589CDB7F29E33AA58DCFE9A799A586CB41191F2B95E3E5B5A4523DB&clip=4&guid=c504bdf377891e043b5a8618b5a4fe27&flowid=eba2893875ef3fde64f9ad295a91bdc8&platform=10901&sdtfrom=v1010&appVer=3.5.57&unid=&auth_from=&auth_ext=&vid=v0044u7dz79&defn=fhd&fhdswitch=0&dtype=3&spsrt=2&tm=1662003290&lang_code=0&logintoken=%7B%22access_token%22%3A%229F78F31093E1865B6A43CF1664064E01%22%2C%22appid%22%3A%22101483052%22%2C%22vusession%22%3A%22AgEsMf2vw4qzX9dabVMd6g.N%22%2C%22openid%22%3A%229DC2DF0053768143217AC15D2E0E1465%22%2C%22vuserid%22%3A%22774829724%22%2C%22video_guid%22%3A%22c504bdf377891e043b5a8618b5a4fe27%22%2C%22main_login%22%3A%22qq%22%7D&spvvpay=1&spadseg=3&hevclv=0&spsfrhdr=0&spvideo=0&drm=40&callback=getinfo_callback_465911
	call := "getinfo_callback_" + strconv.Itoa(rand.Intn(999999-100000)+100000)
	resp, _ := client.R().SetQueryParams(map[string]string{
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
		"guid":       v.Guid,
		"flowid":     flowid,
		"platform":   "10901",
		"sdtfrom":    "v1010",
		"appVer":     "3.5.57",
		"unid":       "",
		"auth_from":  "",
		"auth_ext":   "",
		"vid":        vid,
		"defn":       defn, // 清晰度
		"fhdswitch":  "0",
		"dtype":      "3",
		"spsrt":      "2",
		"tm":         timeStr,
		"lang_code":  "0",
		//"logintoken": `{"access_token":"9F78F31093E1865B6A43CF1664064E01","appid":"101483052","vusession":"AgEsMf2vw4qzX9dabVMd6g.N","openid":"9DC2DF0053768143217AC15D2E0E1465","vuserid":"774829724","video_guid":"c504bdf377891e043b5a8618b5a4fe27","main_login":"qq"}`,
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

	//fmt.Println(s)

	err := json.Unmarshal([]byte(s), &data)
	if err != nil {
		return ""
	}

	//var ndata items
	//ndata = data.(items)
	if len(data.Vl.Vi) == 0 {
		logger.Debug("返回内容", s, err)
		return ""
	}

	for _, item := range data.Vl.Vi[0].Ul.Ui {
		return item.URL
	}

	return ""
}

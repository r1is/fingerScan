package core

import (
	"fingerScan/utils"
	"github.com/PuerkitoBio/goquery"
	"github.com/imroc/req/v3"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

type Resps struct {
	Url        string
	Body       string
	Header     map[string][]string
	Server     string
	Statuscode int
	Length     int
	Title      string
	Jsurl      []string
	Favhash    string
}

func rndua() string {
	ua := []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 YaBrowser/22.1.0.2517 Yowser/2.5 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0",
		"Mozilla/5.0 (X11; Linux x86_64; rv:96.0) Gecko/20100101 Firefox/96.0",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.1 Safari/605.1.15",
		"Mozilla/5.0 (X11; Linux x86_64; rv:95.0) Gecko/20100101 Firefox/95.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:96.0) Gecko/20100101 Firefox/96.0",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.71 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.3 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.110 YaBrowser/22.1.0.2517 Yowser/2.5 Safari/537.36"}
	n := rand.Intn(13) + 1
	return ua[n]
}

func gettitle(httpbody string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(httpbody))
	if err != nil {
		return "Not found"
	}
	title := doc.Find("title").Text()
	title = strings.Replace(title, "\n", "", -1)
	title = strings.Trim(title, " ")
	return title
}

func getfavicon(httpbody string, turl string) string {
	faviconpaths := utils.Xegexpjs(`href="(.*?favicon....)"`, httpbody)
	var faviconpath string
	u, err := url.Parse(turl)
	if err != nil {
		panic(err)
	}
	turl = u.Scheme + "://" + u.Host
	if len(faviconpaths) > 0 {
		fav := faviconpaths[0][1]
		if fav[:2] == "//" {
			faviconpath = "http:" + fav
		} else {
			if fav[:4] == "http" {
				faviconpath = fav
			} else {
				faviconpath = turl + "/" + fav
			}

		}
	} else {
		faviconpath = turl + "/favicon.ico"
	}
	return utils.Favicohash(faviconpath)
}

func Httprequest(url []string, proxy string) (*Resps, error) {
	// 判断是否设置代理，如果设置 加上
	client := req.C()
	//cookie := &http.Cookie{
	//	Name:  "rememberMe",
	//	Value: "me",
	//}
	//client.EnableInsecureSkipVerify().SetUserAgent(rndua()).SetTLSFingerprintChrome().SetTimeout(5 * time.Second).SetCommonCookies(cookie)
	client.EnableInsecureSkipVerify().SetUserAgent(rndua()).SetTLSFingerprintChrome().SetTimeout(5 * time.Second)
	if proxy != "" {
		client.SetProxyURL(proxy)
	}
	response, err := client.R().Get(url[0])
	if err != nil {
		return nil, err
	}
	httpbody, err := response.ToString()
	if err != nil {
		return nil, err
	}
	title := gettitle(httpbody)
	var server string
	capital := response.Header.Get("Server")
	if capital != "" {
		server = capital
	} else {
		powered := response.Header.Get("X-Powered-By")
		if powered != "" {
			server = powered
		} else {
			server = "None"
		}
	}
	var jsurl []string
	if url[1] == "0" {
		jsurl = utils.Jsjump(httpbody, url[0])
	} else {
		jsurl = []string{""}
	}
	favhash := getfavicon(httpbody, url[0])
	s := Resps{url[0], httpbody, response.Header, server, response.StatusCode, len(httpbody), title, jsurl, favhash}
	return &s, nil
}

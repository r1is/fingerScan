package core

import (
	"fingerScan/utils"
	"fmt"
	"github.com/gookit/color"
	"github.com/panjf2000/ants/v2"
	"os"
	"strings"
	"sync"
)

type Outrestul struct {
	Url        string `json:"url"`
	Cms        string `json:"cms"`
	Server     string `json:"server"`
	Statuscode int    `json:"statuscode"`
	Length     int    `json:"length"`
	Title      string `json:"title"`
}

type FinScan struct {
	UrlQueue    *Queue
	Ch          chan []string
	Wg          sync.WaitGroup
	Thread      int
	Output      string
	Proxy       string
	AllResult   []Outrestul
	FocusResult []Outrestul // 重点资产
	Finpx       *Packjson   // 指纹
}

/*
创建扫描
*/
func NewScan(urls []string, thread int, output string, proxy string) *FinScan {
	s := &FinScan{
		UrlQueue:    NewQueue(),
		Ch:          make(chan []string, thread),
		Wg:          sync.WaitGroup{},
		Thread:      thread,
		Output:      output,
		Proxy:       proxy,
		AllResult:   []Outrestul{},
		FocusResult: []Outrestul{},
	}
	err := LoadWebfingerprint()
	if err != nil {
		utils.PrintColorl("[error] fingerprint file error!!!", "237,64,35")
		os.Exit(1)
	}
	s.Finpx = GetWebfingerprint()
	for _, url := range urls {
		s.UrlQueue.Push([]string{url, "0"})
	}
	return s
}

/*
启动扫描
*/
func (s *FinScan) StartScan() {
	//for i := 0; i <= s.Thread; i++ {
	//	s.Wg.Add(1)
	//	go func() {
	//		defer s.Wg.Done()
	//		s.fingerScan()
	//	}()
	//}
	//s.Wg.Wait()

	defer ants.Release()
	pool, _ := ants.NewPool(s.Thread)
	task := func() {
		s.fingerScan()
		s.Wg.Done()
	}
	for i := 0; i < s.UrlQueue.Len(); i++ {
		s.Wg.Add(1)
		_ = pool.Submit(task)
	}
	s.Wg.Wait()

	utils.PrintColorl("\n重点资产：", "244,211,49")
	for _, aas := range s.FocusResult {
		fmt.Printf(fmt.Sprintf("[ %s | ", aas.Url))
		color.RGBStyleFromString("237,64,35").Printf(fmt.Sprintf("%s", aas.Cms))
		fmt.Printf(fmt.Sprintf(" | %s | %d | %d | %s ]\n", aas.Server, aas.Statuscode, aas.Length, aas.Title))
	}
	if s.Output != "" {
		SaveFile(s.Output, s.AllResult)
	}
}

/*
扫描
*/
func (s *FinScan) fingerScan() {
	for s.UrlQueue.Len() != 0 {
		dataface := s.UrlQueue.Pop()
		switch dataface.(type) {
		case []string:
			url := dataface.([]string)
			var data *Resps
			data, err := Httprequest(url, s.Proxy)
			if err != nil {
				url[0] = strings.ReplaceAll(url[0], "https://", "http://")
				data, err = Httprequest(url, s.Proxy)
				if err != nil {
					continue
				}
			}
			for _, jurl := range data.Jsurl {
				if jurl != "" {
					s.UrlQueue.Push([]string{jurl, "1"})
				}
			}
			headers := utils.MapToJson(data.Header)
			var cms []string
			for _, finp := range s.Finpx.Fingerprint {
				if finp.Location == "body" {
					if finp.Method == "keyword" {
						if utils.Iskeyword(data.Body, finp.Keyword) {
							cms = append(cms, finp.Cms)
						}
					}
					if finp.Method == "faviconhash" {
						if data.Favhash == finp.Keyword[0] {
							cms = append(cms, finp.Cms)
						}
					}
					if finp.Method == "regular" {
						if utils.Isregular(data.Body, finp.Keyword) {
							cms = append(cms, finp.Cms)
						}
					}
				}
				if finp.Location == "header" {
					if finp.Method == "keyword" {
						if utils.Iskeyword(headers, finp.Keyword) {
							cms = append(cms, finp.Cms)
						}
					}
					if finp.Method == "regular" {
						if utils.Isregular(headers, finp.Keyword) {
							cms = append(cms, finp.Cms)
						}
					}
				}
				if finp.Location == "title" {
					if finp.Method == "keyword" {
						if utils.Iskeyword(data.Title, finp.Keyword) {
							cms = append(cms, finp.Cms)
						}
					}
					if finp.Method == "regular" {
						if utils.Isregular(data.Title, finp.Keyword) {
							cms = append(cms, finp.Cms)
						}
					}
				}
			}
			cms = utils.RemoveDuplicatesAndEmpty(cms)
			cmss := strings.Join(cms, ",")
			out := Outrestul{data.Url, cmss, data.Server, data.Statuscode, data.Length, data.Title}
			s.AllResult = append(s.AllResult, out)
			if len(out.Cms) != 0 {
				outstr := fmt.Sprintf("[ %s | %s | %s | %d | %d | %s ]", out.Url, out.Cms, out.Server, out.Statuscode, out.Length, out.Title)
				utils.PrintColorl(outstr, "237,64,35")
				s.FocusResult = append(s.FocusResult, out)
			} else {
				outstr := fmt.Sprintf("[ %s | %s | %s | %d | %d | %s ]", out.Url, out.Cms, out.Server, out.Statuscode, out.Length, out.Title)
				fmt.Println(outstr)
			}
		default:
			continue
		}
	}
}

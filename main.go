package main

import (
	"fingerScan/core"
	"fingerScan/utils"
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	file   string
	url    string
	output string
	thread int
	proxy  string
)

func main() {
	flag.StringVar(&file, "f", "", "待识别的文件")
	flag.StringVar(&url, "u", "", "待识别的url")
	flag.StringVar(&output, "o", "", "保存的文件名")
	flag.IntVar(&thread, "t", 100, "扫描线程")
	flag.StringVar(&proxy, "p", "", "代理")
	flag.Parse()
	startTime := time.Now()
	if file != "" {
		localFileUrls := utils.RemoveRepeatedElement(core.LocalFile(file))
		s := core.NewScan(localFileUrls, thread, output, proxy)
		s.StartScan()
		fmt.Println("运行完成,耗时 ", time.Since(startTime))
		os.Exit(1)
	} else if url != "" {
		s := core.NewScan([]string{url}, thread, output, proxy)
		s.StartScan()
		fmt.Println("运行完成,耗时 ", time.Since(startTime))
		os.Exit(1)
	} else {
		flag.PrintDefaults()
	}
}

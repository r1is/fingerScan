package core

import (
	"bufio"
	"encoding/json"
	"fingerScan/utils"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func LocalFile(filename string) (urls []string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Local file read error:", err)
		utils.PrintColorl("[error] the input file is wrong!!!", "237,64,35")
		os.Exit(1)
	}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "http") {
			urls = append(urls, scanner.Text())
		} else {
			urls = append(urls, "https://"+scanner.Text())
		}
	}
	return
}

func SaveFile(filename string, allresult []Outrestul) {
	// 判断allresult是否为空
	if len(allresult) == 0 {
		fmt.Println("allresult is empty")
		return
	}

	//获取后缀名 .json .xlsx .html
	fileExt := filepath.Ext(filename)
	if fileExt == ".json" {
		buf, err := json.MarshalIndent(allresult, "", " ")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		outjson(filename, buf)
	}
	if fileExt == ".xlsx" {
		outxlsx(filename, allresult)
	}
	if fileExt == ".html" {
		outhtml(filename, allresult)
	}

}

func outjson(filename string, data []byte) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err.Error())
		return
	} else {
		defer f.Close()
		_, err = f.Write(data)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func outxlsx(filename string, msg []Outrestul) {
	xlsx := excelize.NewFile()
	xlsx.SetCellValue("Sheet1", "A1", "url")
	xlsx.SetCellValue("Sheet1", "B1", "cms")
	xlsx.SetCellValue("Sheet1", "C1", "server")
	xlsx.SetCellValue("Sheet1", "D1", "statuscode")
	xlsx.SetCellValue("Sheet1", "E1", "length")
	xlsx.SetCellValue("Sheet1", "F1", "title")
	for k, v := range msg {
		xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(k+2), v.Url)
		xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(k+2), v.Cms)
		xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(k+2), v.Server)
		xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(k+2), v.Statuscode)
		xlsx.SetCellValue("Sheet1", "E"+strconv.Itoa(k+2), v.Length)
		xlsx.SetCellValue("Sheet1", "F"+strconv.Itoa(k+2), v.Title)
	}
	err := xlsx.SaveAs(filename)
	if err != nil {
		fmt.Println(err)
	}
}

// 排序规则，将重点资产放到前面，方便在html中查看
type SortByCms []Outrestul

func (a SortByCms) Len() int           { return len(a) }
func (a SortByCms) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByCms) Less(i, j int) bool { return a[i].Cms != "" && a[j].Cms == "" }

func outhtml(filename string, msg []Outrestul) {
	// 创建HTML文件
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	sort.Sort(SortByCms(msg))
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Ehole</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f0f0f0;
            margin: 0;
            padding: 0;
        }
        .container {
            max-width: 80%;
            margin: 20px auto;
            background-color: #fff;
            padding: 20px;
            border-radius: 5px;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
        }
        h1 {
            color: #333;
        }
        table {
            /* width: 1100px; */
            border-collapse: collapse;
            margin-top: 20px;
			margin-left: auto;
			margin-right: auto;
        }
        th, td {
            padding: 10px;
            text-align: left;
            border-bottom: 1px solid #ddd;
			word-break: break-all;
        }
        th {
            background-color: #f0f0f0;
        }
        tr:hover {
            background-color: #f9f9f9;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Ehole 资产</h1>
        <table width="100%" cellpadding="0" cellspacing="0" style="table-layout:fixed">
		<tr>
		<th style="width: 30%;">URL</th>
		<th style="width: 20%;">CMS</th>
		<th style="width: 30%;">Title</th>
		<th style="width: 10%;">Server</th>
		<th style="width: 5%;">SC</th>
		<th style="width: 5%;">Len</th>
		
		</tr>
			{{range .}}

			<tr>
			<td><a href="{{.Url}}" target="_blank">{{.Url}}</a></td>
			<td>{{.Cms}}</td>
			<td>{{.Title}}</td>
			<td>{{.Server}}</td>
			<td>{{.Statuscode}}</td>
			<td>{{.Length}}</td>
			</tr>

			{{end}}
        </table>
    </div>
</body>
</html>
	`

	// 解析HTML模板
	t := template.Must(template.New("html").Parse(tmpl))
	// 将数据写入HTML文件
	err = t.Execute(file, msg)
	if err != nil {
		fmt.Println(err)
	}
}

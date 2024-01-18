package core

import (
	"bufio"
	"encoding/json"
	"fingerScan/utils"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"log"
	"os"
	"strconv"
	"strings"
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
	file := strings.Split(filename, ".")
	if len(file) == 2 {
		if file[1] == "json" {
			buf, err := json.MarshalIndent(allresult, "", " ")
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			outjson(filename, buf)
		}
		if file[1] == "xlsx" {
			outxlsx(filename, allresult)
		}
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

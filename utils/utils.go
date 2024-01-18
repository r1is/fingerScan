package utils

import (
	"encoding/json"
	"github.com/gookit/color"
	"log"
	"os"
	"path/filepath"
)

// 获取当前执行程序所在的绝对路径
func GetCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

func MapToJson(param map[string][]string) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}

func PrintColorl(msg string, colors string) {
	color.RGBStyleFromString(colors).Println(msg)
}

func PrintColorf(msg string, colors string) {
	color.RGBStyleFromString(colors).Printf(msg)
}

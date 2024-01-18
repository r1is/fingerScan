package core

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

type Packjson struct {
	Fingerprint []Fingerprint
}

type Fingerprint struct {
	Cms      string
	Method   string
	Location string
	Keyword  []string
}

//go:embed finger/finger.json
var eHoleFinger string

var (
	Webfingerprint *Packjson
)

func LoadWebfingerprint() error {

	var config Packjson
	err := json.Unmarshal([]byte(eHoleFinger), &config)
	if err != nil {
		fmt.Println(err)
		return err
	}
	Webfingerprint = &config
	return nil
}

func GetWebfingerprint() *Packjson {
	return Webfingerprint
}

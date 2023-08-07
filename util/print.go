package util

import (
	"encoding/json"
)

func PrettifyObject(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

package datautils

import "encoding/json"

func Stringify(data interface{}) (bt []byte, err error) {
	switch t := data.(type) {
	case string:
		bt = []byte(t)
	case []byte:
		bt = t
	default:
		bt, err = json.Marshal(data)
	}

	return
}

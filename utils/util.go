package utils

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

func JSONString(v interface{}) string {
	if v == nil || v == "" {
		return " object is nil  of to json"
	}
	if j, err := json.Marshal(v); err != nil {
		return err.Error()
	} else {
		return string(j)
	}
	return ""
}

func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// encode 序列化要保存的值
func Encode(val interface{}) (interface{}, error) {
	var value interface{}
	switch v := val.(type) {
	case string, int, uint, int8, int16, int32, int64, float32, float64, bool:
		value = v
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		value = string(b)
	}
	return value, nil
}

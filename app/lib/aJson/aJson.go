package aJson

import (
	"errors"
	"fmt"
	"math"

	gjson "github.com/tidwall/gjson"
	sjson "github.com/tidwall/sjson"
)

type Json struct {
	Data string
}

type Result struct {
	gjson.Result
}

func New() *Json {

	return &Json{
		Data: "",
	}
}

func (j *Json) Get(key string) *Result {

	stringValue := string(j.Data)

	value := gjson.Get(stringValue, key)

	result := Result{
		Result: value,
	}

	return &result
}

func (j *Json) Set(key string, value interface{}) {

	if ajson, ok := value.(*Json); ok {

		modifiedJSON, err := sjson.Set(j.Data, key, ajson.Data)
		if err != nil {
			fmt.Println("sjson Set Error:", err)
			return
		}

		//fmt.Println("Modified JSON:", modifiedJSON)

		j.Data = modifiedJSON

	} else {

		modifiedJSON, err := sjson.Set(j.Data, key, value)
		if err != nil {
			fmt.Println("sjson Set Error:", err)
			return
		}

		//fmt.Println("Modified JSON:", modifiedJSON)

		j.Data = modifiedJSON
	}

}

func (j *Json) Exists(key string) bool {

	value := j.Get(key)

	if !value.Exists() {
		return false
	}

	return true
}

func ParseByte(data []byte) (*Json, error) {

	// 将字节切片转换为字符串
	json := string(data)

	return &Json{
		Data: json,
	}, nil
}

func (r *Result) Get(key string) *Result {

	// stringValue := string(j.Data)

	// json, err := r.TryString()

	// if err != nil {
	// 	return nil
	// }

	value := r.Result.Get(key)

	result := Result{
		Result: value,
	}

	return &result
}

func (r *Result) TryMap() (map[string]gjson.Result, error) {

	result := r.Map()

	return result, nil
}

func (r *Result) TryJsonArray() ([]*Json, error) {

	result := r.Array()

	var jsonArray []*Json

	for _, r2 := range result {

		j := new(Json)
		j.Data = r2.String()
		jsonArray = append(jsonArray, j)
	}

	return jsonArray, nil
}

func (r *Result) TryUint64() (uint64, error) {

	if !r.Exists() {
		return 0, errors.New("value no exists")
	}

	result := r.Uint()

	return result, nil
}

func (r *Result) TryString() (string, error) {

	if !r.Exists() {
		return "", errors.New("string value no exists")
	}

	result := r.String()

	return result, nil
}

func (r *Result) TryInt() (int32, error) {

	if !r.Exists() {
		return 0, errors.New("int value no exists")
	}

	result := r.Int()

	// var int64Value int64 = math.MaxInt64
	// fmt.Println("Original int64 value:", int64Value)

	if result > math.MaxInt32 || result < math.MinInt32 {

		fmt.Println("Value is out of int32 range.")

		return 0, errors.New("Value is out of int32 range")

	} else {

		int32Value := int32(result)
		//fmt.Println("Converted int32 value:", int32Value)

		return int32Value, nil
	}
}

func (r *Result) TryInt64() (int64, error) {

	if !r.Exists() {
		return 0, errors.New("value no exists")
	}

	result := r.Int()

	return result, nil
}

func (r *Result) TryBool() (bool, error) {

	if !r.Exists() {
		return false, errors.New("value no exists")
	}

	result := r.Bool()

	return result, nil
}

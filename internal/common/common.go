package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
)

type Check struct {
	Id                  string `json:"id"`
	Name                string `json:"name"`
	Description         string `json:"description"`
	Criteria            string `json:"criteria"`
	RecommendedAction   string `json:"recommendedAction"`
	AdditionalResources string `json:"additionalResources"`
}

func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

func Call(funcName string, checksMap map[string]interface{}, method string, params ...interface{}) (result interface{}, err error) {
	f := reflect.ValueOf(checksMap[funcName]).MethodByName(method)
	if len(params) != f.Type().NumIn() {
		err = errors.New("the number of params is out of index")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	res := f.Call(in)
	result = res[0].Interface()
	return
}

func StringSliceContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
package aws

import (
	"errors"
	"reflect"

	"github.com/brittandeyoung/ckia/cmd/aws/cost"
	"github.com/brittandeyoung/ckia/cmd/aws/security"
)

type checkMapping map[string]interface{}

func buildChecksMap() map[string]interface{} {
	checksMap := checkMapping{
		cost.IdleDBInstanceCheckId:     cost.FindIdleDBInstances,
		security.RootAccountMFACheckId: security.FindRootAccountsMissingMFA,
	}
	return checksMap
}

func Call(funcName string, params ...interface{}) (result interface{}, err error) {
	checksMap := buildChecksMap()
	f := reflect.ValueOf(checksMap[funcName])
	if len(params) != f.Type().NumIn() {
		err = errors.New("The number of params is out of index.")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	var res []reflect.Value
	res = f.Call(in)
	result = res[0].Interface()
	return
}

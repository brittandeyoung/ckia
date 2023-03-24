package aws

import (
	"reflect"
	"testing"

	"github.com/brittandeyoung/ckia/internal/common"
)

func TestChecksMapStructContainsRunMethod(t *testing.T) {
	checksMap := BuildChecksMap()
	for k := range checksMap {
		st := reflect.TypeOf(checksMap[k])
		_, ok := st.MethodByName(common.MethodNameRun)
		if !ok {
			t.Fatalf("Check: (%s) is missing Run() method.", k)
		}
	}
}

func TestChecksMapStructContainsListMethod(t *testing.T) {
	checksMap := BuildChecksMap()
	for k := range checksMap {
		st := reflect.TypeOf(checksMap[k])
		_, ok := st.MethodByName(common.MethodNameList)
		if !ok {
			t.Fatalf("Check: (%s) is missing List() method.", k)
		}
	}
}

package security

import (
	"testing"

	"github.com/brittandeyoung/ckia/internal/create"
)

func TestExpandRootAccountMissingMFA_basic(t *testing.T) {

	summaryMap := make(map[string]int32)
	summaryMap["AccountMFAEnabled"] = 0
	accountNumber := "123456789011"

	account := expandRootAccountMissingMFA(summaryMap, accountNumber)

	if account == (RootAccountMissingMFA{}) {
		create.TestFailureEmptyStruct(t)
	}

	if account.AccountId != "123456789011" {
		create.TestFailureAttribute(t, "AccountId", "123456789011")
	}
}

func TestExpandRootAccountMissingMFA_enabled(t *testing.T) {

	summaryMap := make(map[string]int32)
	summaryMap["AccountMFAEnabled"] = 1
	accountNumber := "123456789011"

	account := expandRootAccountMissingMFA(summaryMap, accountNumber)

	if account != (RootAccountMissingMFA{}) {
		create.TestFailureNonEmptyStruct(t)
	}
}

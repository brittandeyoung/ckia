package create

import "testing"

func TestFailureAttribute(t *testing.T, attribute string, expectedValue string) {
	t.Fatalf(`expand function failed to set the (%s) attribute to the expected value: %s`, attribute, expectedValue)
}

func TestFailureEmptyStruct(t *testing.T) {
	t.Fatalf(`expand function returned an empty structure`)
}

func TestFailureNonEmptyStruct(t *testing.T) {
	t.Fatalf(`expand function returned a non empty structure`)
}

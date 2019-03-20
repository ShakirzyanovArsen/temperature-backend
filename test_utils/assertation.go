package test_utils

import (
	"math"
	"reflect"
	"testing"
)

func UnexpectedError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		t.FailNow()
	}
}
func AssertError(t *testing.T, err error) {
	if err == nil {
		t.Error("expected error")
	}
}

func AssertString(t *testing.T, expected string, actual string) {
	if expected != actual {
		t.Errorf("expected string: %s, actual: %s", expected, actual)
	}
}
func AssertFloat64(t *testing.T, expected float64, actual float64) {
	if math.Mod(expected, actual) > 10e-6 {
		t.Errorf("expected float64: %f , actual: %f (accuracy: 10e-6)", expected, actual)
	}
}
func AssertInt(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Errorf("expected int %d, actual: %d", expected, actual)
	}
}

func AssertNotNil(t *testing.T, ptr interface{}) {
	if ptr == nil {
		t.Errorf("expected not nil pointer")
	}
}

func AssertStruct(t *testing.T, expected interface{}, actual interface{}) {
	if reflect.DeepEqual(expected, actual) {
		t.Errorf("expected struct: %v, actual: %v", expected, actual)
	}
}

package deepequal

import (
	"fmt"
	"reflect"
)

func Compare(expected interface{}, actual interface{}) (bool, error) {
	expectedValue := reflect.ValueOf(expected)
	actualValue := reflect.ValueOf(actual)

	if expectedValue.Type() != actualValue.Type() {
		return false, fmt.Errorf("type mismatch: expected <%s> to be of type <%s>",
			actualValue.Type(), expectedValue.Type())
	}

	switch actualValue.Kind() {
	case reflect.Slice:
		return Slice(expectedValue, actualValue)

	case reflect.Map:
		return Map(expectedValue, actualValue)

	default:
		return deepEqualPrimitive(expected, actual)
	}
}

func deepEqualPrimitive(expectedPrimitive interface{}, actualPrimitive interface{}) (bool, error) {
	if !reflect.DeepEqual(expectedPrimitive, actualPrimitive) {
		return false, fmt.Errorf("value mismatch: expected <%T> %+v to equal <%T> %+v",
			actualPrimitive, actualPrimitive, expectedPrimitive, expectedPrimitive)
	}

	return true, nil
}

type sliceError struct {
	index int
	err   error
}

func (s sliceError) Error() string {
	return fmt.Sprintf("error at slice index %d: %s", s.index, s.err)
}

type mapError struct {
	key interface{}
	err error
}

func (m mapError) Error() string {
	return fmt.Sprintf("error at map key \"%+v\": %s", m.key, m.err)
}

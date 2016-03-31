package gomegamatchers

import (
	"fmt"
	"reflect"
)

func DeepEqual(expected, actual interface{}) error {
	expectedValue := reflect.ValueOf(expected)
	actualValue := reflect.ValueOf(actual)

	if expectedValue.Type() != actualValue.Type() {
		return fmt.Errorf("type mismatch: expected <%s> to be of type <%s>", actualValue.Type(), expectedValue.Type())
	}

	switch actualValue.Kind() {
	case reflect.Slice:
		for i := 0; i < actualValue.Len(); i++ {
			err := DeepEqual(expectedValue.Index(i).Interface(), actualValue.Index(i).Interface())
			if err != nil {
				return sliceError{
					index: i,
					err:   err,
				}
			}
		}

	case reflect.Map:
		for _, key := range actualValue.MapKeys() {
			err := DeepEqual(expectedValue.MapIndex(key).Interface(), actualValue.MapIndex(key).Interface())
			if err != nil {
				return mapError{
					key: key.Interface(),
					err: err,
				}
			}
		}

	default:
		if !reflect.DeepEqual(expected, actual) {
			return fmt.Errorf("value mismatch: expected <%T> %+v to equal <%T> %+v", actual, actual, expected, expected)
		}

	}

	return nil
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

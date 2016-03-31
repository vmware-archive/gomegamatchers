package gomegamatchers

import (
	"fmt"
	"reflect"
	"strings"
)

func DeepEqual(expected interface{}, actual interface{}) (bool, error) {
	expectedValue := reflect.ValueOf(expected)
	actualValue := reflect.ValueOf(actual)

	if expectedValue.Type() != actualValue.Type() {
		return false, fmt.Errorf("type mismatch: expected <%s> to be of type <%s>",
			actualValue.Type(), expectedValue.Type())
	}

	switch actualValue.Kind() {
	case reflect.Slice:
		return deepEqualSlice(expectedValue, actualValue)

	case reflect.Map:
		return deepEqualMap(expectedValue, actualValue)

	default:
		return deepEqualPrimitive(expected, actual)
	}
}

func deepEqualSlice(expectedSlice reflect.Value, actualSlice reflect.Value) (bool, error) {
	for i := 0; i < actualSlice.Len(); i++ {
		if i >= expectedSlice.Len() {
			return false, sliceError{
				index: i,
				err: fmt.Errorf("extra elements found: expected %s not to contain elements %s",
					prettyPrintValueSlice(actualSlice),
					prettyPrintValueSlice(actualSlice.Slice(i, actualSlice.Len())),
				),
			}
		}

		equal, err := DeepEqual(expectedSlice.Index(i).Interface(), actualSlice.Index(i).Interface())
		if !equal {
			return false, sliceError{
				index: i,
				err:   err,
			}
		}
	}

	if expectedSlice.Len() > actualSlice.Len() {
		return false, sliceError{
			index: actualSlice.Len(),
			err: fmt.Errorf("missing elements: expected %s to contain elements %s",
				prettyPrintValueSlice(actualSlice),
				prettyPrintValueSlice(expectedSlice.Slice(actualSlice.Len(), expectedSlice.Len())),
			),
		}
	}

	return true, nil
}

func deepEqualMap(expectedMap reflect.Value, actualMap reflect.Value) (bool, error) {
	for _, key := range actualMap.MapKeys() {
		if expectedMap.MapIndex(key).Kind() == reflect.Invalid {
			return false, mapError{
				key: key.Interface(),
				err: fmt.Errorf("extra key found: expected %s not to contain key <%T> %+v",
					prettyPrintSlice(actualMap.MapKeys()), key.Interface(),
					key,
				),
			}
		}

		equal, err := DeepEqual(expectedMap.MapIndex(key).Interface(), actualMap.MapIndex(key).Interface())
		if !equal {
			return false, mapError{
				key: key.Interface(),
				err: err,
			}
		}
	}

	for _, key := range expectedMap.MapKeys() {
		if actualMap.MapIndex(key).Kind() == reflect.Invalid {
			return false, mapError{
				key: key.Interface(),
				err: fmt.Errorf("missing key: expected %s to contain key <%T> %+v",
					prettyPrintSlice(actualMap.MapKeys()), key.Interface(),
					key,
				),
			}
		}
	}

	return true, nil
}

func deepEqualPrimitive(expectedPrimitive interface{}, actualPrimitive interface{}) (bool, error) {
	if !reflect.DeepEqual(expectedPrimitive, actualPrimitive) {
		return false, fmt.Errorf("value mismatch: expected <%T> %+v to equal <%T> %+v",
			actualPrimitive, actualPrimitive, expectedPrimitive, expectedPrimitive)
	}

	return true, nil
}

func prettyPrintValueSlice(values reflect.Value) string {
	var prettyPrintedValues []string

	for i := 0; i < values.Len(); i++ {
		prettyPrintedValues = append(prettyPrintedValues,
			fmt.Sprintf("<%T> %+v", values.Index(i).Interface(), values.Index(i)))
	}

	return "[" + strings.Join(prettyPrintedValues, ", ") + "]"
}

func prettyPrintSlice(values []reflect.Value) string {
	var prettyPrintedValues []string

	for _, value := range values {
		prettyPrintedValues = append(prettyPrintedValues, fmt.Sprintf("<%T> %+v", value.Interface(), value))
	}

	return "[" + strings.Join(prettyPrintedValues, ", ") + "]"
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

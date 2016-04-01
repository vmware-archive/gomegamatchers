package deepequal

import (
	"fmt"
	"reflect"

	"github.com/pivotal-cf-experimental/gomegamatchers/internal/prettyprint"
)

func Map(expectedMap reflect.Value, actualMap reflect.Value) (bool, error) {
	for _, key := range actualMap.MapKeys() {
		if expectedMap.MapIndex(key).Kind() == reflect.Invalid {
			return false, mapError{
				key: key.Interface(),
				err: fmt.Errorf("extra key found: expected %s not to contain key <%T> %+v",
					prettyprint.SliceOfValues(actualMap.MapKeys()), key.Interface(),
					key,
				),
			}
		}

		equal, err := Compare(expectedMap.MapIndex(key).Interface(), actualMap.MapIndex(key).Interface())
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
					prettyprint.SliceOfValues(actualMap.MapKeys()), key.Interface(),
					key,
				),
			}
		}
	}

	return true, nil
}

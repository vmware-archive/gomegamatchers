package deepequal

import (
	"fmt"
	"reflect"

	"github.com/pivotal-cf-experimental/gomegamatchers/internal/prettyprint"
)

func Slice(expectedSlice reflect.Value, actualSlice reflect.Value) (bool, error) {
	for i := 0; i < actualSlice.Len(); i++ {
		if i >= expectedSlice.Len() {
			return false, sliceError{
				index: i,
				err: fmt.Errorf("extra elements found: expected %s not to contain elements %s",
					prettyprint.SliceAsValue(actualSlice),
					prettyprint.SliceAsValue(actualSlice.Slice(i, actualSlice.Len())),
				),
			}
		}

		equal, err := Compare(expectedSlice.Index(i).Interface(), actualSlice.Index(i).Interface())
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
				prettyprint.SliceAsValue(actualSlice),
				prettyprint.SliceAsValue(expectedSlice.Slice(actualSlice.Len(), expectedSlice.Len())),
			),
		}
	}

	return true, nil
}

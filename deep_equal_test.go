package gomegamatchers_test

import (
	"fmt"

	"github.com/pivotal-cf-experimental/gomegamatchers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeepEqual", func() {
	types := map[string][]interface{}{
		"bool":       {true, false},
		"string":     {"a", "b"},
		"int":        {int(1), int(2)},
		"int8":       {int8(1), int8(2)},
		"int16":      {int16(1), int16(2)},
		"int32":      {int32(1), int32(2)},
		"int64":      {int64(1), int64(2)},
		"uint":       {uint(1), uint(2)},
		"uint8":      {uint8(1), uint8(2)},
		"uint16":     {uint16(1), uint16(2)},
		"uint32":     {uint32(1), uint32(2)},
		"uint64":     {uint64(1), uint64(2)},
		"uintptr":    {uintptr(1), uintptr(2)},
		"float32":    {float32(1.0), float32(2.0)},
		"float64":    {float64(1.0), float64(2.0)},
		"complex64":  {complex64(1i), complex64(2i)},
		"complex128": {complex128(1i), complex128(2i)},
	}

	Context("when the types are mismatched", func() {
		It("returns an error", func() {
			for expectedName, expectedValues := range types {
				for actualName, actualValues := range types {
					if expectedName != actualName {
						errorMessage := fmt.Sprintf("type mismatch: expected <%s> to be of type <%s>", actualName, expectedName)
						err := gomegamatchers.DeepEqual(expectedValues[0], actualValues[0])
						Expect(err).To(MatchError(errorMessage))
					}
				}
			}
		})
	})

	Context("when the values are mismatched", func() {
		It("returns an error", func() {
			for name, values := range types {
				errorMessage := fmt.Sprintf("value mismatch: expected <%s> %+v to equal <%s> %+v", name, values[1], name, values[0])
				err := gomegamatchers.DeepEqual(values[0], values[1])
				Expect(err).To(MatchError(errorMessage))
			}
		})
	})

	Context("when comparing slices", func() {
		It("references the index of the value that errored", func() {
			expected := []int{1, 2, 3, 4}
			actual := []int{1, 2, 0, 4}

			err := gomegamatchers.DeepEqual(expected, actual)
			Expect(err).To(MatchError("error at slice index 2: value mismatch: expected <int> 0 to equal <int> 3"))
		})
	})

	Context("when comparing maps", func() {
		It("references the key of the value that errored", func() {
			expected := map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			}
			actual := map[string]int{
				"a": 1,
				"b": 0,
				"c": 3,
			}

			err := gomegamatchers.DeepEqual(expected, actual)
			Expect(err).To(MatchError(`error at map key "b": value mismatch: expected <int> 0 to equal <int> 2`))
		})
	})

	Context("when comparing complex objects", func() {
		It("references a path indicating the location of the error", func() {
			expected := map[string]interface{}{
				"a": 1,
				"b": []int{1, 2, 3, 4},
				"c": 3,
			}
			actual := map[string]interface{}{
				"a": 1,
				"b": []int{1, 2, 0, 4},
				"c": 3,
			}

			err := gomegamatchers.DeepEqual(expected, actual)
			Expect(err).To(MatchError(`error at map key "b": error at slice index 2: value mismatch: expected <int> 0 to equal <int> 3`))
		})
	})
})

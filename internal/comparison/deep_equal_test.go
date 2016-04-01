package comparison_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/gomegamatchers/internal/comparison"
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
						equal, err := comparison.DeepEqual(expectedValues[0], actualValues[0])
						Expect(equal).To(BeFalse())
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
				equal, err := comparison.DeepEqual(values[0], values[1])
				Expect(equal).To(BeFalse())
				Expect(err).To(MatchError(errorMessage))
			}
		})
	})

	Context("when comparing slices", func() {
		It("returns an error when keys match but values do not", func() {
			expected := []int{1, 2, 3, 4}
			actual := []int{1, 2, 0, 4}

			equal, err := comparison.DeepEqual(expected, actual)
			Expect(equal).To(BeFalse())
			Expect(err).To(MatchError("error at slice index 2: value mismatch: expected <int> 0 to equal <int> 3"))
		})

		It("returns an error when the actual slice contains values that are not in the expected slice", func() {
			expected := []int{1, 2}
			actual := []int{1, 2, 3, 4}

			equal, err := comparison.DeepEqual(expected, actual)
			Expect(equal).To(BeFalse())
			Expect(err).To(MatchError("error at slice index 2: extra elements found: expected [<int> 1, <int> 2, <int> 3, <int> 4] not to contain elements [<int> 3, <int> 4]"))
		})

		It("returns an error when the expected slice contains values that are not in the actual slice", func() {
			expected := []int{1, 2, 3, 4}
			actual := []int{1, 2}

			equal, err := comparison.DeepEqual(expected, actual)
			Expect(equal).To(BeFalse())
			Expect(err).To(MatchError("error at slice index 2: missing elements: expected [<int> 1, <int> 2] to contain elements [<int> 3, <int> 4]"))
		})
	})

	Context("when comparing maps", func() {
		It("returns true when the keys and values match (regardless of the order of the keys)", func() {
			expected := map[string]int{"a": 1, "b": 2, "c": 3}
			actual := map[string]int{"c": 3, "b": 2, "a": 1}

			equal, _ := comparison.DeepEqual(expected, actual)
			Expect(equal).To(BeTrue())
		})

		It("returns an error when keys match but values do not", func() {
			expected := map[string]int{"a": 1, "b": 2, "c": 3}
			actual := map[string]int{"a": 1, "b": 0, "c": 3}

			equal, err := comparison.DeepEqual(expected, actual)
			Expect(equal).To(BeFalse())
			Expect(err).To(MatchError(`error at map key "b": value mismatch: expected <int> 0 to equal <int> 2`))
		})

		It("returns an error when the actual map contains keys that are not in the expected map", func() {
			expected := map[string]int{"a": 1}
			actual := map[string]int{"a": 1, "b": 2}

			equal, err := comparison.DeepEqual(expected, actual)
			Expect(equal).To(BeFalse())
			Expect(err).To(SatisfyAny(
				MatchError(`error at map key "b": extra key found: expected [<string> a, <string> b] not to contain key <string> b`),
				MatchError(`error at map key "b": extra key found: expected [<string> b, <string> a] not to contain key <string> b`),
			))
		})

		It("returns an error when the expected map contains keys that are not in the actual map", func() {
			expected := map[string]int{"a": 1, "b": 2}
			actual := map[string]int{"a": 1}

			equal, err := comparison.DeepEqual(expected, actual)
			Expect(equal).To(BeFalse())
			Expect(err).To(MatchError(`error at map key "b": missing key: expected [<string> a] to contain key <string> b`))
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

			equal, err := comparison.DeepEqual(expected, actual)
			Expect(equal).To(BeFalse())
			Expect(err).To(MatchError(`error at map key "b": error at slice index 2: value mismatch: expected <int> 0 to equal <int> 3`))
		})
	})
})

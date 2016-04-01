package deepequal_test

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/gomegamatchers/internal/deepequal"
)

var _ = Describe("Map", func() {
	It("returns true when the keys and values match (regardless of the order of the keys)", func() {
		expected := reflect.ValueOf(map[string]int{"a": 1, "b": 2, "c": 3})
		actual := reflect.ValueOf(map[string]int{"c": 3, "b": 2, "a": 1})

		equal, _ := deepequal.Map(expected, actual)
		Expect(equal).To(BeTrue())
	})

	It("returns an error when keys match but values do not", func() {
		expected := reflect.ValueOf(map[string]int{"a": 1, "b": 2, "c": 3})
		actual := reflect.ValueOf(map[string]int{"a": 1, "b": 0, "c": 3})

		equal, err := deepequal.Map(expected, actual)
		Expect(equal).To(BeFalse())
		Expect(err).To(MatchError(`error at map key "b": value mismatch: expected <int> 0 to equal <int> 2`))
	})

	It("returns an error when the actual map contains keys that are not in the expected map", func() {
		expected := reflect.ValueOf(map[string]int{"a": 1})
		actual := reflect.ValueOf(map[string]int{"a": 1, "b": 2})

		equal, err := deepequal.Map(expected, actual)
		Expect(equal).To(BeFalse())
		Expect(err).To(SatisfyAny(
			MatchError(`error at map key "b": extra key found: expected [<string> a, <string> b] not to contain key <string> b`),
			MatchError(`error at map key "b": extra key found: expected [<string> b, <string> a] not to contain key <string> b`),
		))
	})

	It("returns an error when the expected map contains keys that are not in the actual map", func() {
		expected := reflect.ValueOf(map[string]int{"a": 1, "b": 2})
		actual := reflect.ValueOf(map[string]int{"a": 1})

		equal, err := deepequal.Map(expected, actual)
		Expect(equal).To(BeFalse())
		Expect(err).To(MatchError(`error at map key "b": missing key: expected [<string> a] to contain key <string> b`))
	})
})

package deepequal_test

import (
	"reflect"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-cf-experimental/gomegamatchers/internal/deepequal"
)

var _ = Describe("Slice", func() {
	It("returns true when the lengths and values match", func() {
		slice := reflect.ValueOf([]int{1, 2, 3, 4})

		equal, err := deepequal.Slice(slice, slice)
		Expect(equal).To(BeTrue())
		Expect(err).NotTo(HaveOccurred())
	})

	It("returns an error when the lengths match but the values do not", func() {
		expected := reflect.ValueOf([]int{1, 2, 3, 4})
		actual := reflect.ValueOf([]int{1, 2, 0, 4})

		equal, err := deepequal.Slice(expected, actual)
		Expect(equal).To(BeFalse())
		Expect(err).To(MatchError("error at slice index 2: value mismatch: expected <int> 0 to equal <int> 3"))
	})

	It("returns an error when the actual slice contains values that are not in the expected slice", func() {
		expected := reflect.ValueOf([]int{1, 2})
		actual := reflect.ValueOf([]int{1, 2, 3, 4})

		equal, err := deepequal.Slice(expected, actual)
		Expect(equal).To(BeFalse())
		Expect(err).To(MatchError("error at slice index 2: extra elements found: expected [<int> 1, <int> 2, <int> 3, <int> 4] not to contain elements [<int> 3, <int> 4]"))
	})

	It("returns an error when the expected slice contains values that are not in the actual slice", func() {
		expected := reflect.ValueOf([]int{1, 2, 3, 4})
		actual := reflect.ValueOf([]int{1, 2})

		equal, err := deepequal.Slice(expected, actual)
		Expect(equal).To(BeFalse())
		Expect(err).To(MatchError("error at slice index 2: missing elements: expected [<int> 1, <int> 2] to contain elements [<int> 3, <int> 4]"))
	})
})

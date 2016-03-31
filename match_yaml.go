package gomegamatchers

import (
	"fmt"
	"reflect"

	"github.com/cloudfoundry-incubator/candiedyaml"
	"github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

func MatchYAML(expected interface{}) types.GomegaMatcher {
	return &MatchYAMLMatcher{expected}
}

type MatchYAMLMatcher struct {
	YAMLToMatch interface{}
}

func (matcher *MatchYAMLMatcher) Match(actual interface{}) (success bool, err error) {
	actualString, err := matcher.prettyPrint(actual)
	if err != nil {
		return false, err
	}

	expectedString, err := matcher.prettyPrint(matcher.YAMLToMatch)
	if err != nil {
		return false, err
	}

	var actualValue interface{}
	var expectedValue interface{}

	// this is guarded by prettyPrint
	candiedyaml.Unmarshal([]byte(actualString), &actualValue)
	candiedyaml.Unmarshal([]byte(expectedString), &expectedValue)

	return reflect.DeepEqual(actualValue, expectedValue), nil
}

func (matcher *MatchYAMLMatcher) FailureMessage(actual interface{}) (message string) {
	actualString, _ := matcher.prettyPrint(actual)
	expectedString, _ := matcher.prettyPrint(matcher.YAMLToMatch)
	return format.Message(actualString, "to match YAML of", expectedString)
}

func (matcher *MatchYAMLMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	actualString, _ := matcher.prettyPrint(actual)
	expectedString, _ := matcher.prettyPrint(matcher.YAMLToMatch)
	return format.Message(actualString, "not to match YAML of", expectedString)
}

func (matcher *MatchYAMLMatcher) prettyPrint(input interface{}) (formatted string, err error) {
	inputString, ok := toString(input)
	if !ok {
		return "", fmt.Errorf("MatchYAMLMatcher matcher requires a string or stringer.  Got:\n%s", format.Object(input, 1))
	}

	var data interface{}
	if err := candiedyaml.Unmarshal([]byte(inputString), &data); err != nil {
		return "", err
	}
	buf, _ := candiedyaml.Marshal(data)

	return string(buf), nil
}

func toString(value interface{}) (string, bool) {
	valueString, isString := value.(string)
	if isString {
		return valueString, true
	}

	valueBytes, isBytes := value.([]byte)
	if isBytes {
		return string(valueBytes), true
	}

	valueStringer, isStringer := value.(fmt.Stringer)
	if isStringer {
		return valueStringer.String(), true
	}

	return "", false
}

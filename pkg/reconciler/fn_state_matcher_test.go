package reconciler

import (
	"fmt"
	"regexp"

	"github.com/onsi/gomega/types"
)

func equalStateFunction(sFn stateFn) types.GomegaMatcher {
	return &stateFnNameMatcher{
		expected: sFn,
	}
}

type stateFnNameMatcher struct {
	expected stateFn
}

func lastIndexOf(s string, of rune) int {
	result := -1
	for i, r := range s {
		if r == of {
			result = i
		}
	}
	return result
}

func (m *stateFnNameMatcher) Match(actual any) (success bool, err error) {
	actualFn, ok := actual.(stateFn)
	if !ok {
		return false, fmt.Errorf("stateFnNameMatcher expects string")
	}

	r := regexp.MustCompile(".func[0-9]*|.glob")

	actName := r.ReplaceAllString(actualFn.name(), "")
	expName := r.ReplaceAllString(m.expected.name(), "")

	aliof := lastIndexOf(actName, '.')
	if aliof != -1 {
		actName = actName[aliof:]
	}

	eliof := lastIndexOf(expName, '.')
	if eliof != -1 {
		expName = expName[eliof:]
	}

	return actName == expName, nil
}

func (m *stateFnNameMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%s\nto be equal to\n\t%s", actual, m.expected)
}

func (m *stateFnNameMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%s\nnot to be equal to\n\t%s", actual, m.expected)
}

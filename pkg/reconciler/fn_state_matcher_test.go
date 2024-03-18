package reconciler

import (
	"fmt"
	"regexp"

	"github.com/onsi/gomega/types"
)

func equalStateFunction(sFn stateFn) types.GomegaMatcher {
	return &stateFnNameMatcher{
		Expected: sFn,
	}
}

type stateFnNameMatcher struct {
	Expected stateFn
	actName  string
	expName  string
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
		return false, fmt.Errorf("stateFnNameMatcher expects stateFn")
	}

	r := regexp.MustCompile(".func[0-9]*|.glob|.[0-9]*")

	m.actName = r.ReplaceAllString(actualFn.name(), "")
	m.expName = r.ReplaceAllString(m.Expected.name(), "")

	aliof := lastIndexOf(m.actName, '.')
	if aliof != -1 {
		m.actName = m.actName[aliof+1:]
	}

	eliof := lastIndexOf(m.expName, '.')
	if eliof != -1 {
		m.expName = m.expName[eliof+1:]
	}

	fmt.Println("match:", m.actName, m.expName, m.Expected.name())

	return m.actName == m.expName, nil
}

func (m *stateFnNameMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%s\nto be equal to\n\t%s", m.actName, m.expName)
}

func (m *stateFnNameMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%s\nnot to be equal to\n\t%s", m.actName, m.expName)
}

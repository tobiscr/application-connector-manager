package reconciler

import (
	"fmt"

	"github.com/kyma-project/application-connector-manager/api/v1alpha1"
	"github.com/onsi/gomega/types"
)

func haveExpectedDomainName(domainName string) types.GomegaMatcher {
	return &domainNameMatcher{
		expectedDomainName: domainName,
	}
}

type domainNameMatcher struct {
	expectedDomainName string
	actualDomainName   string
}

func (m *domainNameMatcher) Match(actual any) (success bool, err error) {
	appConn, ok := actual.(v1alpha1.ApplicationConnector)
	if !ok {
		return false, fmt.Errorf("stateFnNameMatcher expects unstructured.Unstructured")
	}
	m.actualDomainName = appConn.Spec.DomainName
	return m.actualDomainName == m.expectedDomainName, nil
}

func (m *domainNameMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%v\nto be equal\n\t%s", m.actualDomainName, m.expectedDomainName)
}

func (m *domainNameMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%v\nnot to be equal\n\t%s", m.actualDomainName, m.expectedDomainName)
}

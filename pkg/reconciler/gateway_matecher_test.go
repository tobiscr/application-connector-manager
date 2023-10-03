package reconciler

import (
	"fmt"

	"github.com/kyma-project/application-connector-manager/pkg/unstructured"
	"github.com/onsi/gomega/types"
	"golang.org/x/exp/slices"
	istio "istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/runtime"
)

type gatewayMatcher struct {
	expected    string
	actualHosts []string
}

func haveDomainNamePropagatedInGateway(v string) types.GomegaMatcher {
	return &gatewayMatcher{
		expected: v,
	}
}

func (m *gatewayMatcher) Match(actual any) (success bool, err error) {
	u, ok := actual.(unstructured.Unstructured)
	if !ok {
		return false, fmt.Errorf("stateFnNameMatcher expects unstructured.Unstructured")
	}

	var actualGateway istio.Gateway
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &actualGateway); err != nil {
		return false, fmt.Errorf("conversion error: %w", err)
	}

	if len(actualGateway.Spec.Servers) != 1 {
		return false, fmt.Errorf("invalid gateway servers cound: %d", len(actualGateway.Spec.Servers))
	}

	m.actualHosts = actualGateway.Spec.Servers[0].Hosts

	return slices.Contains(actualGateway.Spec.Servers[0].Hosts, m.expected), nil
}

func (m *gatewayMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%v\nto contain\n\t%s", m.actualHosts, m.expected)
}

func (m *gatewayMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%v\nnot to contain\n\t%s", m.actualHosts, m.expected)
}

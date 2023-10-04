package reconciler

import (
	"fmt"

	"github.com/kyma-project/application-connector-manager/pkg/unstructured"
	"github.com/onsi/gomega/types"
	"golang.org/x/exp/slices"
	istio "istio.io/client-go/pkg/apis/networking/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

type virtualServiceMatcher struct {
	expected    string
	actualHosts []string
}

func haveDomainNamePropagatedInVirtualService(v string) types.GomegaMatcher {
	return &virtualServiceMatcher{
		expected: v,
	}
}

func (m *virtualServiceMatcher) Match(actual any) (success bool, err error) {
	u, ok := actual.(unstructured.Unstructured)
	if !ok {
		return false, fmt.Errorf("stateFnNameMatcher expects unstructured.Unstructured")
	}

	var actualCirtualService istio.VirtualService
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &actualCirtualService); err != nil {
		return false, fmt.Errorf("conversion error: %w", err)
	}

	if len(actualCirtualService.Spec.Hosts) != 1 {
		return false, fmt.Errorf("Invalid virtual service host cound: %d", len(actualCirtualService.Spec.Hosts))
	}

	m.actualHosts = actualCirtualService.Spec.Hosts

	return slices.Contains(actualCirtualService.Spec.Hosts, m.expected), nil
}

func (m *virtualServiceMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%v\nto contain\n\t%s", m.actualHosts, m.expected)
}

func (m *virtualServiceMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%v\nnot to contain\n\t%s", m.actualHosts, m.expected)
}

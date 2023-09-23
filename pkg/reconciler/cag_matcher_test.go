package reconciler

import (
	"fmt"

	"github.com/kyma-project/application-connector-manager/api/v1alpha1"
	"github.com/kyma-project/application-connector-manager/pkg/unstructured"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"golang.org/x/exp/slices"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type cagMatcher struct {
	expected   v1alpha1.AppGatewaySpec
	actualArgs []string
	expectArgs []string
}

func haveAppGatewaySpec(v v1alpha1.AppGatewaySpec) types.GomegaMatcher {
	return &cagMatcher{expected: v}
}

func (m *cagMatcher) Match(actual any) (success bool, err error) {
	u, ok := actual.(unstructured.Unstructured)
	if !ok {
		return false, fmt.Errorf("stateFnNameMatcher expects unstructured.Unstructured")
	}

	var actualDeployment appsv1.Deployment
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &actualDeployment); err != nil {
		return false, fmt.Errorf("conversion error: %w", err)
	}

	index := slices.IndexFunc(actualDeployment.Spec.Template.Spec.Containers, func(c corev1.Container) bool {
		return c.Name == "central-application-gateway"
	})
	if index == -1 {
		return false, fmt.Errorf("central-application-gateway container not found")
	}
	// create expected command argumets
	m.expectArgs = []string{
		fmt.Sprintf("%s=%.0f", v1alpha1.ArgCentralAppGatewayRequestTimeout, m.expected.RequestTimeout.Seconds()),
		fmt.Sprintf("%s=%.0f", v1alpha1.ArgCentralAppGatewayProxyTimeout, m.expected.ProxyTimeout.Seconds()),
	}
	m.actualArgs = actualDeployment.Spec.Template.Spec.Containers[index].Args
	containsArgs := gomega.ContainElements(m.expectArgs)
	return containsArgs.Match(m.actualArgs)
}

func (m *cagMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%s\nto be equal to\n\t%s", m.actualArgs, m.expectArgs)
}

func (m *cagMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%s\nnot to be equal to\n\t%s", m.actualArgs, m.expectArgs)
}

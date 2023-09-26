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

//nolint:unused // remove on phase2: compass-runtime-agent in module
type craMatcher struct {
	expected   v1alpha1.RuntimeAgentSpec
	actualVars []corev1.EnvVar
	expectVars []corev1.EnvVar
}

//nolint:unused // remove on phase2: compass-runtime-agent in module
func haveRuntimeAgentSpec(rtAgentSpec v1alpha1.RuntimeAgentSpec) types.GomegaMatcher {
	return &craMatcher{expected: rtAgentSpec}
}

//nolint:unused // remove on phase2: compass-runtime-agent in module
func (m *craMatcher) Match(actual any) (success bool, err error) {
	u, ok := actual.(unstructured.Unstructured)
	if !ok {
		return false, fmt.Errorf("stateFnNameMatcher expects unstructured.Unstructured")
	}

	var actualDeployment appsv1.Deployment
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &actualDeployment); err != nil {
		return false, fmt.Errorf("conversion error: %w", err)
	}

	index := slices.IndexFunc(actualDeployment.Spec.Template.Spec.Containers, func(c corev1.Container) bool {
		return c.Name == "compass-runtime-agent"
	})
	if index == -1 {
		return false, fmt.Errorf("compass-runtime-agent container not found")
	}
	// create expected environment variables base on spec
	m.expectVars = []corev1.EnvVar{
		{
			Name:  v1alpha1.EnvRuntimeAgentControllerSyncPeriod,
			Value: m.expected.ControllerSyncPeriod.Duration.String(),
		},
		{
			Name:  v1alpha1.EnvRuntimeAgentCertValidityRenevalThreshold,
			Value: m.expected.CertValidityRenewalThreshold,
		},
		{
			Name:  v1alpha1.EnvRuntimeAgentMinimalCompassSyncTime,
			Value: m.expected.MinConfigSyncTime.Duration.String(),
		},
	}
	m.actualVars = actualDeployment.Spec.Template.Spec.Containers[index].Env
	containsEnvs := gomega.ContainElements(m.expectVars)
	return containsEnvs.Match(m.actualVars)
}

//nolint:unused // remove on phase2: compass-runtime-agent in module
func (m *craMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%s\nto be equal to\n\t%s", m.actualVars, m.expectVars)
}

//nolint:unused // remove on phase2: compass-runtime-agent in module
func (m *craMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%s\nnot to be equal to\n\t%s", m.actualVars, m.expectVars)
}

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

type appConnValidatorMatcher struct {
	expected   v1alpha1.AppConnValidatorSpec
	actualVars []corev1.EnvVar
	expectVars []corev1.EnvVar
}

func haveAppConnValidatorSpec(v v1alpha1.AppConnValidatorSpec) types.GomegaMatcher {
	return &appConnValidatorMatcher{expected: v}
}

func (m *appConnValidatorMatcher) Match(actual any) (success bool, err error) {
	u, ok := actual.(unstructured.Unstructured)
	if !ok {
		return false, fmt.Errorf("stateFnNameMatcher expects unstructured.Unstructured")
	}

	var actualDeployment appsv1.Deployment
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &actualDeployment); err != nil {
		return false, fmt.Errorf("conversion error: %w", err)
	}

	index := slices.IndexFunc(actualDeployment.Spec.Template.Spec.Containers, func(c corev1.Container) bool {
		return c.Name == "central-application-connectivity-validator"
	})
	if index == -1 {
		return false, fmt.Errorf("central-application-connectivity-validator container not found")
	}
	// create expected envs
	m.expectVars = []corev1.EnvVar{
		{
			Name:  v1alpha1.EnvAppConnValidatorLogLevel,
			Value: string(m.expected.LogLevel),
		},
		{
			Name:  v1alpha1.EnvAppConnValidatorLogFormat,
			Value: string(m.expected.LogFormat),
		},
	}
	m.actualVars = actualDeployment.Spec.Template.Spec.Containers[index].Env
	containsEnvs := gomega.ContainElements(m.expectVars)
	return containsEnvs.Match(m.actualVars)
}

func (m *appConnValidatorMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%s\nto be contain\n\t%s", m.actualVars, m.expectVars)
}

func (m *appConnValidatorMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%s\nnot to contain\n\t%s", m.actualVars, m.expectVars)
}

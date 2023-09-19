package reconciler

import (
	"github.com/kyma-project/application-connector-manager/api/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// the state of controlled system (k8s cluster)
type systemState struct {
	instance v1alpha1.ApplicationConnector
	// the state of module component parts on cluster used detect
	// module readiness
	objs []unstructured.Unstructured

	snapshot v1alpha1.Status
}

func (s *systemState) saveAppConStatus() {
	result := s.instance.Status.DeepCopy()
	if result == nil {
		result = &v1alpha1.Status{}
	}
	s.snapshot = *result
}

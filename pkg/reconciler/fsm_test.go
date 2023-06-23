package reconciler

import (
	"errors"
	"testing"

	"github.com/kyma-project/application-connector-manager/pkg/unstructured"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	apirt "k8s.io/apimachinery/pkg/runtime"
)

const (
	operatorName = "application-connector-manager"
)

func Test_updateObj_convert_errors(t *testing.T) {
	var errTest = errors.New("test error")

	type args struct {
		toUnstructed   func(interface{}) (map[string]interface{}, error)
		fromUnstructed func(map[string]interface{}, interface{}) error
	}

	u := unstructured.Unstructured{}
	u.SetName(operatorName)
	u.SetAPIVersion("apps/v1")
	u.SetKind("Deployment")

	tests := []struct {
		name          string
		args          args
		expectedError error
	}{
		{
			name: "from unstructed fail",
			args: args{
				fromUnstructed: func(u map[string]interface{}, obj interface{}) error {
					return errTest
				},
				toUnstructed: apirt.DefaultUnstructuredConverter.ToUnstructured,
			},
			expectedError: errTest,
		},
		{
			name: "to unstructed fail",
			args: args{
				toUnstructed: func(obj interface{}) (map[string]interface{}, error) {
					return nil, errTest
				},
				fromUnstructed: apirt.DefaultUnstructuredConverter.FromUnstructured,
			},
			expectedError: errTest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toUnstructed = tt.args.toUnstructed
			fromUnstructured = tt.args.fromUnstructed

			err := unstructured.Update(&u, nil, func(*appsv1.Deployment, interface{}) error {
				t.Log("deployment updated")
				return nil
			})

			g := NewWithT(t)

			g.Expect(err).Should(HaveOccurred())
			g.Expect(err).Should(Equal(tt.expectedError))
		})
	}

}

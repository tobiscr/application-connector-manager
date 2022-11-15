package controllers

/*import (
	"context"
	"fmt"
	"github.com/kyma-project/application-connector-manager/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

type testHelper struct {
	ctx           context.Context
	namespaceName string
}

func shouldCreateApplicationConnector(h testHelper, appConnecotSpec v1alpha1.ApplicationConnectorSpec) {
	// act
	h.createApplicationConnector(kedaName, kedaSpec)

	// we have to update deployment status manually
	//h.updateDeploymentStatus(metricsDeploymentName)
	//h.updateDeploymentStatus(kedaDeploymentName)

	// assert
	Eventually(h.createGetKedaStateFunc(kedaName)).
		WithPolling(time.Second * 2).
		WithTimeout(time.Second * 20).
		Should(Equal(rtypes.StateReady))
}

func (h *testHelper) createApplicationConnector(appConnecti string, spec v1alpha1.ApplicationConnectorSpec) {
	By(fmt.Sprintf("Creating crd: %s", kedaName))
	keda := v1alpha1.Keda{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kedaName,
			Namespace: h.namespaceName,
			Labels: map[string]string{
				"operator.kyma-project.io/kyma-name": "test",
			},
		},
		Spec: spec,
	}
	Expect(k8sClient.Create(h.ctx, &keda)).To(Succeed())
	By(fmt.Sprintf("Crd created: %s", kedaName))
}*/

package controllers

import (
	"context"
	"fmt"
	"github.com/kyma-project/application-connector-manager/api/v1alpha1"
	rtypes "github.com/kyma-project/module-manager/operator/pkg/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

type testHelper struct {
	ctx           context.Context
	namespaceName string
}

var _ = Describe("ApplicationConnector controller", func() {
	Context("When creating fresh instance", func() {
		const (
			namespaceName                                  = "kyma-system"
			appConnectorName                               = "test"
			applicationGatewayDeploymentName               = "central-application-gateway"
			applicationConnectivityValidatorDeploymentName = "central-application-connectivity-validator"
		)

		var (
			appConnectorSpec = v1alpha1.ApplicationConnectorSpec{
				DisableLegacyConnectivity: true,
			}

			//appConnectorUpdateSpec = v1alpha1.ApplicationConnectorSpec{
			//	DisableLegacyConnectivity: false,
			//}
		)

		It("The status should be Success", func() {
			h := testHelper{
				ctx:           context.Background(),
				namespaceName: namespaceName,
			}

			h.createNamespace()

			// operations like C(R)UD can be tested in separated tests,
			// but we have time-consuming flow and decided do it in one test
			shouldCreateApplicationConnector(h, appConnectorName, applicationGatewayDeploymentName, applicationConnectivityValidatorDeploymentName, appConnectorSpec)

			//shouldPropagateAppConnectorCrdSpecProperties(h, applicationGatewayDeploymentName, applicationConnectivityValidatorDeploymentName, appConnectorSpec)

			//TODO: disabled because of bug in operator (https://github.com/kyma-project/module-manager/issues/94)
			//shouldUpdateAppConnector(h, appConnectorName, applicationGatewayDeploymentName)
			//shouldDeleteAppConnector(h, appConnectorName)
		})
	})
})

func shouldCreateApplicationConnector(h testHelper, appConnectorName, appGatewayDeploymentName, appConValidatorDeploymentName string, appConnecorSpec v1alpha1.ApplicationConnectorSpec) {
	// act
	h.createApplicationConnector(appConnectorName, appConnecorSpec)

	// we have to update deployment status manually
	h.updateDeploymentStatus(appGatewayDeploymentName)
	h.updateDeploymentStatus(appConValidatorDeploymentName)

	// assert
	Eventually(h.createGetApplicationConnectorStateFunc(appConnectorName)).
		WithPolling(time.Second * 2).
		WithTimeout(time.Second * 20).
		Should(Equal(rtypes.StateProcessing))
}

func (h *testHelper) createGetApplicationConnectorStateFunc(appConnectorName string) func() (rtypes.State, error) {
	return func() (rtypes.State, error) {
		return h.getApplicationConnectorState(appConnectorName)
	}
}

func (h *testHelper) createGetNStatusFunc(namespaceName string) func() (rtypes.State, error) {
	return func() (rtypes.State, error) {
		return h.getNamespaceStatus(namespaceName)
	}
}

func (h *testHelper) getNamespaceStatus(namespaceName string) (rtypes.State, error) {
	return rtypes.StateReady, nil
}

func (h *testHelper) getApplicationConnectorState(appConnName string) (rtypes.State, error) {
	var emptyState = rtypes.State("")
	var connector v1alpha1.ApplicationConnector
	key := types.NamespacedName{
		Name:      appConnName,
		Namespace: h.namespaceName,
	}
	err := k8sClient.Get(h.ctx, key, &connector)
	if err != nil {
		return emptyState, err
	}
	return rtypes.State(connector.Status.State), nil
}

func (h *testHelper) createApplicationConnector(appConnectorName string, spec v1alpha1.ApplicationConnectorSpec) {
	By(fmt.Sprintf("Creating crd: %s", appConnectorName))
	appconnector := v1alpha1.ApplicationConnector{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appConnectorName,
			Namespace: h.namespaceName,
			Labels: map[string]string{
				"operator.kyma-project.io/kyma-name": "test",
			},
		},
		Spec: spec,
	}
	Expect(k8sClient.Create(h.ctx, &appconnector)).To(Succeed())
	By(fmt.Sprintf("Crd created: %s", appConnectorName))
}

func (h *testHelper) updateDeploymentStatus(deploymentName string) {
	By(fmt.Sprintf("Updating deployment status: %s", deploymentName))
	var deployment appsv1.Deployment
	Eventually(h.createGetKubernetesObjectFunc(deploymentName, &deployment)).
		WithPolling(time.Second * 2).
		WithTimeout(time.Second * 10).
		Should(BeTrue())

	deployment.Status.Conditions = append(deployment.Status.Conditions, appsv1.DeploymentCondition{
		Type:    appsv1.DeploymentAvailable,
		Status:  corev1.ConditionTrue,
		Reason:  "test-reason",
		Message: "test-message",
	})
	deployment.Status.Replicas = 1
	Expect(k8sClient.Status().Update(h.ctx, &deployment)).To(Succeed())

	replicaSetName := h.createReplicaSetForDeployment(deployment)

	var replicaSet appsv1.ReplicaSet
	Eventually(h.createGetKubernetesObjectFunc(replicaSetName, &replicaSet)).
		WithPolling(time.Second * 2).
		WithTimeout(time.Second * 10).
		Should(BeTrue())

	replicaSet.Status.ReadyReplicas = 1
	replicaSet.Status.Replicas = 1
	Expect(k8sClient.Status().Update(h.ctx, &replicaSet)).To(Succeed())

	By(fmt.Sprintf("Deployment status updated: %s", deploymentName))
}

func (h *testHelper) createGetKubernetesObjectFunc(serviceAccountName string, obj client.Object) func() (bool, error) {
	return func() (bool, error) {
		key := types.NamespacedName{
			Name:      serviceAccountName,
			Namespace: h.namespaceName,
		}
		err := k8sClient.Get(h.ctx, key, obj)
		if err != nil {
			return false, err
		}
		return true, err
	}
}

func (h *testHelper) createReplicaSetForDeployment(deployment appsv1.Deployment) string {
	replicaSetName := fmt.Sprintf("%s-replica-set", deployment.Name)
	By(fmt.Sprintf("Creating replica set (for deployment): %s", replicaSetName))
	var (
		trueValue = true
		one       = int32(1)
	)
	replicaSet := appsv1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      replicaSetName,
			Namespace: h.namespaceName,
			Labels: map[string]string{
				"app": deployment.Name,
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: "apps/v1",
					Kind:       "Deployment",
					Name:       deployment.Name,
					UID:        deployment.GetUID(),
					Controller: &trueValue,
				},
			},
		},
		// dummy values
		Spec: appsv1.ReplicaSetSpec{
			Replicas: &one,
			Selector: deployment.Spec.Selector,
			Template: deployment.Spec.Template,
		},
	}
	Expect(k8sClient.Create(h.ctx, &replicaSet)).To(Succeed())
	By(fmt.Sprintf("Replica set (for deployment) created: %s", replicaSetName))
	return replicaSetName
}

func (h *testHelper) createNamespace() {
	By(fmt.Sprintf("Creating namespace: %s", h.namespaceName))
	namespace := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: h.namespaceName,
		},
	}
	Expect(k8sClient.Create(h.ctx, &namespace)).To(Succeed())
	By(fmt.Sprintf("Namespace created: %s", h.namespaceName))
}

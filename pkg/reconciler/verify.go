package reconciler

import (
	"context"
	"fmt"
	"time"

	"github.com/kyma-project/application-connector-manager/api/v1alpha1"
	"github.com/kyma-project/application-connector-manager/pkg/unstructured"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apirt "k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	msgVerificationInProgress = "verification in progress"
)

var (
	fromUnstructured     = apirt.DefaultUnstructuredConverter.FromUnstructured
	defaultRequeDuration = time.Minute * 2
)

func validateDeployment(obj unstructured.Unstructured) (bool, error) {
	var deployment appsv1.Deployment
	if err := fromUnstructured(obj.Object, &deployment); err != nil {
		return false, err
	}

	for _, cond := range deployment.Status.Conditions {
		if cond.Type == appsv1.DeploymentAvailable && cond.Status == corev1.ConditionTrue {
			return true, nil
		}
	}
	return false, nil
}

func validateApiXtV1Beta1CRD(obj unstructured.Unstructured) (bool, error) {
	var crd apiextv1.CustomResourceDefinition
	if err := fromUnstructured(obj.Object, &crd); err != nil {
		return false, err
	}

	for _, cond := range crd.Status.Conditions {

		isEstablished := cond.Type == apiextv1.Established
		isConditionTrue := cond.Status == apiextv1.ConditionTrue

		if isEstablished && isConditionTrue {
			return true, nil
		}

		isNamesAccepted := cond.Type == apiextv1.NamesAccepted
		isConditionFalse := cond.Status == apiextv1.ConditionFalse

		if isNamesAccepted && isConditionFalse {
			return true, nil
		}
	}
	return false, nil
}

func validateService(obj unstructured.Unstructured) (bool, error) {
	var service corev1.Service
	if err := fromUnstructured(obj.Object, &service); err != nil {
		return false, err
	}

	// service does not have cluster IP address
	if service.Spec.ClusterIP == "" {
		return false, nil
	}

	if service.Spec.Type != corev1.ServiceTypeLoadBalancer {
		return true, nil
	}

	// service ready - some of external IPs are set
	if len(service.Spec.ExternalIPs) > 0 {
		return true, nil
	}

	// service does not have load balancer ingress IP address
	if service.Status.LoadBalancer.Ingress == nil {
		return false, nil
	}

	return true, nil
}

type validate = func(unstructured.Unstructured) (bool, error)

func sFnVerify(_ context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	inventory := inventory(map[string]bool{})

	for _, obj := range s.objs {
		var validate validate

		if unstructured.IsDeploymentKind(obj) {
			validate = validateDeployment
		}

		if unstructured.IsServiceKind(obj) {
			validate = validateService
		}

		if unstructured.IsApiXtV1Beta1CRDKind(obj) {
			validate = validateApiXtV1Beta1CRD
		}

		if validate == nil {
			// omit validation
			continue
		}

		valid, err := validate(obj)
		if err != nil {
			s.instance.UpdateStateFromErr(
				v1alpha1.ConditionTypeInstalled,
				v1alpha1.ConditionReasonVerificationErr,
				err,
			)
			return stopWithErrorAndNoRequeue(err)
		}

		key := fmt.Sprintf("%s/%s/%s", obj.GetNamespace(), obj.GetKind(), obj.GetName())
		inventory[key] = valid
	}

	if !inventory.ready() {
		s.instance.UpdateStateProcessing(
			v1alpha1.ConditionTypeInstalled,
			v1alpha1.ConditionReasonVerification,
			msgVerificationInProgress,
		)

		ready, total := inventory.count()
		m.log.Infof("resources not ready: [%d/%d]", ready, total)
		return stopWithNoRequeue()
	}

	if s.instance.Status.State == "Ready" {
		return stopWithRequeueAfter(defaultRequeDuration)
	}

	s.instance.UpdateStateReady(
		v1alpha1.ConditionTypeInstalled,
		v1alpha1.ConditionReasonVerified,
		"application-connector-manager ready",
	)
	return stopWithRequeueAfter(defaultRequeDuration)
}

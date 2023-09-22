package reconciler

import (
	"context"
	"fmt"
	"time"

	"github.com/kyma-project/application-connector-manager/api/v1alpha1"
	"github.com/kyma-project/application-connector-manager/pkg/unstructured"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
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
		if cond.Type == appsv1.DeploymentAvailable && cond.Status == v1.ConditionTrue {
			return true, nil
		}
	}
	return false, nil
}

type validate = func(unstructured.Unstructured) (bool, error)

func sFnVerify(_ context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	inventory := inventory(map[string]bool{})

	for _, obj := range s.objs {
		var validate validate

		if unstructured.IsDeploymentKind(obj) {
			validate = validateDeployment
		}

		if validate == nil {
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
		m.log.Infof("deployments not ready: [%d/%d]", ready, total)
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

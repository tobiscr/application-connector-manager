package reconciler

import (
	"context"
	"fmt"

	"github.com/kyma-project/application-connector-manager/api/v1alpha1"
	"github.com/kyma-project/application-connector-manager/pkg/unstructured"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	msgVerificationInProgress = "verification in progress"
)

func sFnVerify(_ context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	deployments := inventory(map[string]bool{})

loopObj:
	for _, obj := range s.objs {
		if !unstructured.IsDeploymentKind(obj) {
			continue
		}

		var deployment appsv1.Deployment
		if err := fromUnstructured(obj.Object, &deployment); err != nil {
			s.instance.UpdateStateFromErr(
				v1alpha1.ConditionTypeInstalled,
				v1alpha1.ConditionReasonVerificationErr,
				err,
			)
			return stopWithErrorAndNoRequeue(err)
		}

		key := fmt.Sprintf("%s/%s", deployment.GetNamespace(), deployment.GetName())
		deployments[key] = false

		for _, cond := range deployment.Status.Conditions {
			if cond.Type == appsv1.DeploymentAvailable && cond.Status == v1.ConditionTrue {
				deployments[key] = true
				continue loopObj
			}
		}
	}

	if !deployments.ready() {
		s.instance.UpdateStateProcessing(
			v1alpha1.ConditionTypeInstalled,
			v1alpha1.ConditionReasonVerification,
			msgVerificationInProgress,
		)

		ready, total := deployments.count()
		m.log.Infof("deployments not ready: [%d/%d]", ready, total)
		return stopWithNoRequeue()
	}

	if s.instance.Status.State == "Ready" {
		return nil, nil, nil
	}

	s.instance.UpdateStateReady(
		v1alpha1.ConditionTypeInstalled,
		v1alpha1.ConditionReasonVerified,
		"application-connector-manager ready",
	)
	return stopWithNoRequeue()
}

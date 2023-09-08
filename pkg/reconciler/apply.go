package reconciler

import (
	"context"
	"errors"
	"github.com/kyma-project/application-connector-manager/api/v1alpha1"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	InstallationErr = errors.New("installation error")
)

func sFnApply(ctx context.Context, r *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	var isError bool
	for _, obj := range r.Objs {
		r.log.
			With("gvk", obj.GetObjectKind().GroupVersionKind()).
			With("name", obj.GetName()).
			With("ns", obj.GetNamespace()).
			Debug("applying")

		err := r.Patch(ctx, &obj, client.Apply, &client.PatchOptions{
			Force:        ptr.To[bool](true),
			FieldManager: "application-connector-manager",
		})

		if err != nil {
			r.log.With("err", err).Error("apply error")
			isError = true
		}

		s.objs = append(s.objs, obj)
	}
	// no errors
	if !isError {
		return switchState(sFnVerify)
	}

	s.instance.UpdateStateFromErr(
		v1alpha1.ConditionTypeInstalled,
		v1alpha1.ConditionReasonApplyObjError,
		InstallationErr,
	)
	return stopWithNoRequeue()
}

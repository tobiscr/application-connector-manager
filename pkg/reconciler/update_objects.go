package reconciler

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnUpdate(ctx context.Context, r *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	return switchState(sFnApply)
}

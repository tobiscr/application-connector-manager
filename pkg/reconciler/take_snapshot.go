package reconciler

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
)

func sFnTakeSnapshot(_ context.Context, _ *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	s.saveAppConStatus()
	return sFnCheckDependencies, nil, nil
}

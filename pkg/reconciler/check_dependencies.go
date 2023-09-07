package reconciler

import (
	"context"
	"errors"

	"golang.org/x/exp/slices"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	_ "sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	ErrIstioNotFound = errors.New("ISTIO not found")
	istioGKS         = []schema.GroupKind{
		{
			Group: "networking.istio.io",
			Kind:  "VirtualService",
		},
		{
			Group: "networking.istio.io",
			Kind:  "Gateway",
		},
	}
)

func checkDeps(crds []v1.CustomResourceDefinition, gks ...schema.GroupKind) error {
	var ackCount int
	for _, crd := range crds {
		isGK := func(gk schema.GroupKind) bool {
			return gk.Group == crd.Spec.Group && gk.Kind == crd.Spec.Names.Kind
		}

		if isOneOfGKS := slices.ContainsFunc(gks, isGK); isOneOfGKS {
			ackCount++
		}
	}

	dependencyCount := len(istioGKS)
	if ackCount != dependencyCount {
		return ErrIstioNotFound
	}
	return nil
}

func sFnCheckDependencies(ctx context.Context, r *fsm, _ *systemState) (stateFn, *ctrl.Result, error) {
	if *r.dependencyACK {
		return switchState(sFnInitialize)
	}

	var crds v1.CustomResourceDefinitionList
	if err := r.List(ctx, &crds); err != nil {
		return stopWithErrorAndNoRequeue(err)
	}

	if err := checkDeps(crds.Items, istioGKS...); err != nil {
		return stopWithErrorAndNoRequeue(err)
	}

	return switchState(sFnRegisterDependencyWatch)
}

func sFnRegisterDependencyWatch(ctx context.Context, r *fsm, _ *systemState) (stateFn, *ctrl.Result, error) {

	*r.dependencyACK = true
	r.log.Info("dependency ack")

	return switchState(sFnInitialize)
}

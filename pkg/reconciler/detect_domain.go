package reconciler

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

var gardenerCM types.NamespacedName = types.NamespacedName{
	Name:      "shoot-info",
	Namespace: "kube-system",
}

func sFnDetectDomain(ctx context.Context, r *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	if s.instance.Spec.DomainName != "" {
		return switchState(sFnUpdate)
	}

	// check if domain name was already fetched
	if s.domainName != "" {
		s.instance.Spec.DomainName = s.domainName
		return switchState(sFnUpdate)
	}

	// try to fetch domain name from config map
	var cm v1.ConfigMap
	if err := r.Get(ctx, gardenerCM, &cm); err != nil {
		r.log.Warn("unable to domain name: %s", err)
		return switchState(sFnUpdate)
	}

	domainName, found := cm.Data["domain"]
	if !found {
		r.log.Warn("unable to domain name: not found")
		return switchState(sFnUpdate)
	}

	s.instance.Spec.DomainName = domainName
	return switchState(sFnUpdate)
}

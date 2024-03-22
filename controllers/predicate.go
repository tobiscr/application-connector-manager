package controllers

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type predicateCompassRtAgentSecret struct {
	objectName string
	namespace  string
	predicate.ResourceVersionChangedPredicate
	log *zap.SugaredLogger
}

func (p *predicateCompassRtAgentSecret) Update(e event.UpdateEvent) bool {
	return false
}

func (p *predicateCompassRtAgentSecret) Delete(e event.DeleteEvent) bool {
	return e.Object.GetNamespace() == p.namespace &&
		e.Object.GetName() == p.objectName &&
		p.ResourceVersionChangedPredicate.Delete(e)
}

func (p *predicateCompassRtAgentSecret) Create(e event.CreateEvent) bool {
	return e.Object.GetNamespace() == p.namespace &&
		e.Object.GetName() == p.objectName &&
		p.ResourceVersionChangedPredicate.Create(e)
}

func (p *predicateCompassRtAgentSecret) Generic(e event.GenericEvent) bool {
	return false
}

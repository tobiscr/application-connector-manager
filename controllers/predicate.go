package controllers

import (
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

//nolint:unused // remove on phase2: compass-runtime-agent in module
type predicateCompassRtAgentGenChange struct {
	objectName string
	namespace  string
	predicate.GenerationChangedPredicate
	log *zap.SugaredLogger
}

//nolint:unused // remove on phase2: compass-runtime-agent in module
func (p *predicateCompassRtAgentGenChange) Update(e event.UpdateEvent) bool {
	return e.ObjectNew.GetNamespace() == p.namespace &&
		e.ObjectNew.GetName() == p.objectName &&
		p.GenerationChangedPredicate.Update(e)
}

//nolint:unused // remove on phase2: compass-runtime-agent in module
func (p *predicateCompassRtAgentGenChange) Delete(e event.DeleteEvent) bool {
	return e.Object.GetNamespace() == p.namespace &&
		e.Object.GetName() == p.objectName &&
		p.GenerationChangedPredicate.Delete(e)
}

//nolint:unused // remove on phase2: compass-runtime-agent in module
func (p *predicateCompassRtAgentGenChange) Create(e event.CreateEvent) bool {
	return e.Object.GetNamespace() == p.namespace &&
		e.Object.GetName() == p.objectName &&
		p.GenerationChangedPredicate.Create(e)
}

//nolint:unused // remove on phase2: compass-runtime-agent in module
func (p *predicateCompassRtAgentGenChange) Generic(e event.GenericEvent) bool {
	return e.Object.GetNamespace() == p.namespace &&
		e.Object.GetName() == p.objectName &&
		p.GenerationChangedPredicate.Generic(e)
}

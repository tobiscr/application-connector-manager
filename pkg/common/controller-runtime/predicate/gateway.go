package predicate

import (
	"reflect"

	"go.uber.org/zap"
	istio "istio.io/client-go/pkg/apis/networking/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type gatewayPredicate struct {
	predicate.ResourceVersionChangedPredicate
	log *zap.SugaredLogger
}

func NewGatewayPredicate(log *zap.SugaredLogger) predicate.Predicate {
	return &gatewayPredicate{
		log: log,
	}
}

func (p gatewayPredicate) Update(e event.UpdateEvent) bool {
	var oldVS istio.Gateway
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(e.ObjectOld.(*unstructured.Unstructured).Object, &oldVS); err != nil {
		p.log.Warnf("unable to convert old gateway: %w", err)
		return true
	}

	var newVS istio.VirtualService
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(e.ObjectNew.(*unstructured.Unstructured).Object, &newVS); err != nil {
		p.log.Warnf("unable to convert new gateway: %w", err)
		return true
	}

	// check if status changed
	if statusEqual := reflect.DeepEqual(&oldVS.Status, &newVS.Status); !statusEqual {
		return true
	}

	// check if spec changed
	if specEqual := reflect.DeepEqual(&oldVS.Spec, &newVS.Spec); !specEqual {
		return true
	}

	// check if labels changed
	if labelsEqual := reflect.DeepEqual(oldVS.GetLabels(), newVS.GetLabels()); !labelsEqual {
		return true
	}

	// check if annotations changed
	if annotationsEqual := reflect.DeepEqual(oldVS.GetAnnotations(), newVS.GetAnnotations()); !annotationsEqual {
		return true
	}

	// check if namespace changed
	if oldVS.GetNamespace() != newVS.GetNamespace() {
		return true
	}

	return false
}

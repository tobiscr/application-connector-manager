package reconciler

import (
	"context"
	"errors"
	"fmt"

	"github.com/kyma-project/application-connector-manager/api/v1alpha1"
	"github.com/kyma-project/application-connector-manager/pkg/common/types"
	"github.com/kyma-project/application-connector-manager/pkg/unstructured"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	defaultDeletionStrategy = cascadeDeletionStrategy
)

func updateFinalizers(deps []unstructured.Unstructured, finalizer string) (result []unstructured.Unstructured) {
	for _, obj := range deps {
		newObj := obj.DeepCopy()
		if !controllerutil.RemoveFinalizer(newObj, finalizer) {
			// omit object
			continue
		}

		result = append(result, *newObj)
	}
	return
}

var ErrDeletionFailed = errors.New("installation error")

type list = func(context.Context, client.ObjectList, ...client.ListOption) error

func listUnstruct(ctx context.Context, gvk schema.GroupVersionKind, list list) ([]unstructured.Unstructured, error) {
	var u unstructured.UnstructuredList
	u.SetGroupVersionKind(gvk)

	err := list(ctx, &u, &client.ListOptions{
		Namespace: "kyma-system",
	})

	return u.Items, err
}

func listDependencies(ctx context.Context, list list) ([]unstructured.Unstructured, error) {
	var out []unstructured.Unstructured
	for _, gvk := range []schema.GroupVersionKind{
		types.VirtualService,
		types.Gateway,
	} {
		result, err := listUnstruct(ctx, gvk, list)
		if err != nil {
			return nil, err
		}
		out = append(out, result...)
	}
	return out, nil
}

func sFnDeleteResources(ctx context.Context, m *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	if !isDeleting(s) {
		s.instance.UpdateStateDeletion(
			v1alpha1.ConditionTypeInstalled,
			v1alpha1.ConditionReasonDeletion,
			"deletion in progress",
		)

		return stopWithRequeue()
	}

	// to remove finalizers operate directly on k8s objects
	deps, err := listDependencies(ctx, m.List)
	if err != nil {
		s.instance.UpdateStateFromErr(
			v1alpha1.ConditionTypeInstalled,
			v1alpha1.ConditionReasonDeletionErr,
			ErrDeletionFailed,
		)
		return stopWithErrorAndNoRequeue(err)
	}

	updated := updateFinalizers(deps, "application-connector-manager.kyma-project.io/deletion-hook")
	if len(updated) == 0 {
		return switchState(deletionStrategyBuilder(defaultDeletionStrategy))
	}

	var isError bool
	for _, obj := range updated {
		err := m.Update(ctx, &obj)
		if err != nil {
			m.log.With("err", err).Error("unable to remove finalizer")
			isError = true
		}
	}
	// no errors
	if !isError {
		return stopWithRequeue()
	}

	s.instance.UpdateStateFromErr(
		v1alpha1.ConditionTypeInstalled,
		v1alpha1.ConditionReasonDeletionErr,
		ErrDeletionFailed,
	)

	return stopWithErrorAndNoRequeue(fmt.Errorf("%w: unable to remove dependency finalizer[s]", ErrDeletionFailed))
}

type deletionStrategy string

const (
	cascadeDeletionStrategy  deletionStrategy = "cascadeDeletionStrategy"
	safeDeletionStrategy     deletionStrategy = "safeDeletionStrategy"
	upstreamDeletionStrategy deletionStrategy = "upstreamDeletionStrategy"
)

func deletionStrategyBuilder(strategy deletionStrategy) stateFn {
	switch strategy {
	case cascadeDeletionStrategy:
		return sFnCascadeDeletionState
	case upstreamDeletionStrategy:
		return sFnUpstreamDeletionState
	case safeDeletionStrategy:
		return sFnSafeDeletionState
	default:
		return deletionStrategyBuilder(safeDeletionStrategy)
	}
}

func sFnCascadeDeletionState(ctx context.Context, r *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	return deleteResourcesWithFilter(ctx, r, s)
}

func sFnUpstreamDeletionState(ctx context.Context, r *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	return deleteResourcesWithFilter(ctx, r, s, withoutCRDFilter)
}

func sFnSafeDeletionState(ctx context.Context, r *fsm, s *systemState) (stateFn, *ctrl.Result, error) {
	if err := checkCRDOrphanResources(ctx, r); err != nil {
		s.instance.UpdateStateFromErr(
			v1alpha1.ConditionTypeInstalled,
			v1alpha1.ConditionReasonDeletionErr,
			err,
		)

		// stop state machine with an error and requeue reconciliation
		return stopWithErrorAndNoRequeue(err)
	}

	return deleteResourcesWithFilter(ctx, r, s)
}

func withoutCRDFilter(u unstructured.Unstructured) bool {
	return !isCRD(u)
}

type filterFunc func(unstructured.Unstructured) bool

func deleteResourcesWithFilter(ctx context.Context, r *fsm, s *systemState, filterFunc ...filterFunc) (stateFn, *ctrl.Result, error) {
	var err error
	for _, obj := range append(r.Objs, r.Deps...) {
		if !fitToFilters(obj, filterFunc...) {
			r.log.
				With("objName", obj.GetName()).
				With("gvk", obj.GroupVersionKind()).
				Debug("skipped")
			continue
		}

		r.log.
			With("objName", obj.GetName()).
			With("gvk", obj.GroupVersionKind()).
			Debug("deleting")

		err = r.Delete(ctx, &obj)
		err = client.IgnoreNotFound(err)

		if err != nil {
			r.log.With("deleting resource").Error(err)
		}
	}

	if err != nil {
		s.instance.UpdateStateFromErr(
			v1alpha1.ConditionTypeInstalled,
			v1alpha1.ConditionReasonDeletionErr,
			ErrDeletionFailed,
		)
		// stop state machine with an error and requeue reconciliation
		return stopWithErrorAndNoRequeue(err)
	}
	return switchState(sFnRemoveFinalizer)
}

func fitToFilters(u unstructured.Unstructured, filterFunc ...filterFunc) bool {
	for _, fn := range filterFunc {
		if !fn(u) {
			return false
		}
	}

	return true
}

func checkCRDOrphanResources(ctx context.Context, r *fsm) error {
	for _, obj := range append(r.Objs, r.Deps...) {
		if !isCRD(obj) {
			continue
		}

		crdList, err := buildResourceListFromCRD(obj)
		if err != nil {
			return err
		}

		err = r.List(ctx, &crdList)
		if err != nil {
			return err
		}

		if len(crdList.Items) > 0 {
			return fmt.Errorf("found %d items with VersionKind %s", len(crdList.Items), crdList.GetAPIVersion())
		}
	}

	return nil
}

func isCRD(u unstructured.Unstructured) bool {
	return u.GroupVersionKind().GroupKind() == apiextensionsv1.Kind("CustomResourceDefinition")
}

func buildResourceListFromCRD(u unstructured.Unstructured) (unstructured.UnstructuredList, error) {
	crd := apiextensionsv1.CustomResourceDefinition{}
	crdList := unstructured.UnstructuredList{}

	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.Object, &crd)
	if err != nil {
		return crdList, err
	}

	crdList.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   crd.Spec.Group,
		Version: getCRDStoredVersion(crd),
		Kind:    crd.Spec.Names.Kind,
	})

	return crdList, nil
}

func getCRDStoredVersion(crd apiextensionsv1.CustomResourceDefinition) string {
	for _, version := range crd.Spec.Versions {
		if version.Storage {
			return version.Name
		}
	}

	return ""
}

func isDeleting(s *systemState) bool {
	condition := meta.FindStatusCondition(s.instance.Status.Conditions, string(v1alpha1.ConditionTypeInstalled))
	if condition == nil {
		return false
	}

	if condition.Reason != string(v1alpha1.ConditionReasonDeletion) &&
		condition.Reason != string(v1alpha1.ConditionReasonDeletionErr) {
		return false
	}

	return true
}

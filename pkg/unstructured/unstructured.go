package unstructured

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	apirt "k8s.io/apimachinery/pkg/runtime"
)

var (
	fromUnstructured = apirt.DefaultUnstructuredConverter.FromUnstructured
	toUnstructed     = apirt.DefaultUnstructuredConverter.ToUnstructured
)

type Unstructured = unstructured.Unstructured

// updates given object by applying provided function with given data
func Update[T any, R any](u *Unstructured, data R, update func(T, R) error) error {
	var result T
	err := fromUnstructured(u.Object, &result)
	if err != nil {
		return err
	}

	err = update(result, data)
	if err != nil {
		return err
	}

	u.Object, err = toUnstructed(&result)
	return err
}

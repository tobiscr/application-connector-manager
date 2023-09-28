package unstructured

import (
	"fmt"

	"github.com/kyma-project/application-connector-manager/pkg/common/types"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var ErrNotFound = fmt.Errorf("not found")

type Predicate func(Unstructured) bool

func (p Predicate) All(u []Unstructured) ([]*Unstructured, error) {
	var out []*Unstructured
	for i := range u {
		if !p(u[i]) {
			continue
		}
		out = append(out, &u[i])
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("%w: no object for given predicate", ErrNotFound)
	}
	return out, nil
}

func (p Predicate) First(u []Unstructured) (*Unstructured, error) {
	for i := range u {
		if !p(u[i]) {
			continue
		}
		return &u[i], nil
	}
	return nil, fmt.Errorf("%w: no object for given predicate", ErrNotFound)
}

func hasName(u Unstructured, name string) bool {
	return u.GetName() == name
}

func IsDeploymentKind(u Unstructured) bool {
	return u.GetKind() == "Deployment" && u.GetAPIVersion() == "apps/v1"
}

func IsServiceKind(u Unstructured) bool {
	return u.GetKind() == "Service" && u.GetAPIVersion() == "v1"
}

func IsApiXtV1Beta1CRDKind(u Unstructured) bool {
	return u.GetKind() == "CustomResourceDefinition" && u.GetAPIVersion() == "apiextensions.k8s.io"
}

func isGroupVersionKind(u Unstructured, gvk schema.GroupVersionKind) bool {
	return u.GroupVersionKind() == gvk
}

func IsNamedGroupVersionKind(name string, gvk schema.GroupVersionKind) Predicate {
	return func(u Unstructured) bool {
		return isGroupVersionKind(u, gvk) && hasName(u, name)
	}
}

func IsGatewayKind() Predicate {
	return func(u Unstructured) bool {
		return u.GetKind() == types.Gateway.Kind &&
			u.GetAPIVersion() == fmt.Sprintf("%s/%s", types.Gateway.Group, types.Gateway.Version)
	}
}

func isVirtualServiceKind(u Unstructured) bool {
	return u.GetKind() == types.VirtualService.Kind &&
		u.GetAPIVersion() == fmt.Sprintf("%s/%s", types.VirtualService.Group, types.VirtualService.Version)
}

func IsVirtualService() Predicate {
	return func(u Unstructured) bool {
		return isVirtualServiceKind(u)
	}
}

func IsNamedVirtualService(name string) Predicate {
	return func(u Unstructured) bool {
		return isVirtualServiceKind(u) && hasName(u, name)
	}
}

func IsDeployment(name string) Predicate {
	return func(u Unstructured) bool {
		return IsDeploymentKind(u) && hasName(u, name)
	}
}

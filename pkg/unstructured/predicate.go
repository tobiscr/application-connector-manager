package unstructured

import (
	"fmt"
)

var ErrNotFound = fmt.Errorf("not found")

type Predicate func(Unstructured) bool

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

func isDeploymentKind(u Unstructured) bool {
	return u.GetKind() == "Deployment" &&
		u.GetAPIVersion() == "apps/v1"
}

func IsDeployment(name string) Predicate {
	return func(u Unstructured) bool {
		return isDeploymentKind(u) && hasName(u, name)
	}

}

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

func (p Predicate) All(u []Unstructured) ([]Unstructured, error) {
	result := make([]Unstructured, 0)
	for _, u := range u {
		if p(u) {
			result = append(result, u)
		}
	}
	return result, nil
}

func hasName(u Unstructured, name string) bool {
	return u.GetName() == name
}

func IsDeploymentKind(u Unstructured) bool {
	return u.GetKind() == "Deployment" && u.GetAPIVersion() == "apps/v1"
}

func IsDeployment(name string) Predicate {
	return func(u Unstructured) bool {
		return IsDeploymentKind(u) && hasName(u, name)
	}
}

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Status struct {
	State      string             `json:"state"`
	Served     string             `json:"served"`
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

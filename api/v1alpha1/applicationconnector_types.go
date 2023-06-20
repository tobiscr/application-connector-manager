/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConditionReason string

type ConditionType string

const (
	StateReady      = "Ready"
	StateError      = "Error"
	StateProcessing = "Processing"
	StateDeleting   = "Deleting"

	ServedTrue  = "True"
	ServedFalse = "False"

	ConditionReasonDeploymentUpdateErr = ConditionReason("KedaDeploymentUpdateErr")
	ConditionReasonVerificationErr     = ConditionReason("VerificationErr")
	ConditionReasonVerified            = ConditionReason("Verified")
	ConditionReasonApplyObjError       = ConditionReason("ApplyObjError")
	ConditionReasonVerification        = ConditionReason("Verification")
	ConditionReasonInitialized         = ConditionReason("Initialized")
	ConditionReasonKedaDuplicated      = ConditionReason("KedaDuplicated")
	ConditionReasonDeletion            = ConditionReason("Deletion")
	ConditionReasonDeletionErr         = ConditionReason("DeletionErr")
	ConditionReasonDeleted             = ConditionReason("Deleted")

	ConditionTypeInstalled = ConditionType("Installed")
	ConditionTypeDeleted   = ConditionType("Deleted")

	Finalizer = "application-connector-manager.kyma-project.io/deletion-hook"
)

// ApplicationConnectorSpec defines the desired state of ApplicationConnector
type ApplicationConnectorSpec struct {
	DisableLegacyConnectivity bool `json:"disableLegacyConnectivity"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ApplicationConnector is the Schema for the applicationconnectors API
type ApplicationConnector struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationConnectorSpec `json:"spec,omitempty"`
	Status Status                   `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ApplicationConnectorList contains a list of ApplicationConnector
type ApplicationConnectorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApplicationConnector `json:"items"`
}

func (k *ApplicationConnector) UpdateStateProcessing(c ConditionType, r ConditionReason, msg string) {
	k.Status.State = StateProcessing
	condition := metav1.Condition{
		Type:               string(c),
		Status:             "Unknown",
		LastTransitionTime: metav1.Now(),
		Reason:             string(r),
		Message:            msg,
	}
	meta.SetStatusCondition(&k.Status.Conditions, condition)
}

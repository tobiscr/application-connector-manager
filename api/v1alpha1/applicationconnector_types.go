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
	"k8s.io/apimachinery/pkg/api/meta"
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

	ConditionReasonVerificationErr = ConditionReason("VerificationErr")
	ConditionReasonVerified        = ConditionReason("Verified")
	ConditionReasonApplyObjError   = ConditionReason("ApplyObjError")
	ConditionReasonVerification    = ConditionReason("Verification")
	ConditionReasonInitialized     = ConditionReason("Initialized")
	ConditionReasonDeletion        = ConditionReason("Deletion")
	ConditionReasonDeletionErr     = ConditionReason("DeletionErr")
	ConditionReasonDeleted         = ConditionReason("Deleted")

	ConditionTypeInstalled = ConditionType("Installed")
	ConditionTypeDeleted   = ConditionType("Deleted")

	Finalizer = "application-connector-manager.kyma-project.io/deletion-hook"
)

// +kubebuilder:validation:Enum=debug;info;warn;error
type LogLevel string

// +kubebuilder:validation:Enum=json;console
type LogFormat string

type AppGatewaySpec struct {
	ProxyTimeout   string   `json:"proxyTimeout"`
	RequestTimeout string   `json:"requestTimeout"`
	LogLevel       LogLevel `json:"logLevel"`
}

type AppConnValidatorSpec struct {
	LogLevel  LogLevel  `json:"logLevel"`
	LogFormat LogFormat `json:"logFormat"`
}

type RuntimeAgentSpec struct {
	ControllerSyncPeriod         metav1.Duration `json:"controllerSyncPeriod"`
	MinConfigSyncTime            string          `json:"minimalConfigSyncTime"`
	CertValidityRenewalThreshold string          `json:"certValidityRenewalThreshold"`
}

// ApplicationConnectorSpec contains configuration of ApplicationConnector module and its state

type ApplicationConnectorSpec struct {
	// +optional
	// +kubebuilder:default:={ proxyTimeout: "10", requestTimeout: "10", logLevel: "info" }
	ApplicationGatewaySpec AppGatewaySpec `json:"appGateway"`
	// +optional
	// +kubebuilder:default:={ logLevel: "info", logFormat: "json" }
	AppConValidatorSpec AppConnValidatorSpec `json:"appConnValidator"`
	// +optional
	// +kubebuilder:default:={ controllerSyncPeriod: "90s", minimalConfigSyncTime: "15s", certValidityRenewalThreshold: "0.3" }
	RuntimeAgentSpec RuntimeAgentSpec `json:"runtimeAgent"`
	DomainName       string           `json:"domainName"`
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

func (k *ApplicationConnector) UpdateStateFromErr(c ConditionType, r ConditionReason, err error) {
	k.Status.State = StateError
	condition := metav1.Condition{
		Type:               string(c),
		Status:             "False",
		LastTransitionTime: metav1.Now(),
		Reason:             string(r),
		Message:            err.Error(),
	}
	meta.SetStatusCondition(&k.Status.Conditions, condition)
}

func (k *ApplicationConnector) UpdateStateReady(c ConditionType, r ConditionReason, msg string) {
	k.Status.State = StateReady
	condition := metav1.Condition{
		Type:               string(c),
		Status:             "True",
		LastTransitionTime: metav1.Now(),
		Reason:             string(r),
		Message:            msg,
	}
	meta.SetStatusCondition(&k.Status.Conditions, condition)
}

func (k *ApplicationConnector) UpdateStateDeletion(c ConditionType, r ConditionReason, msg string) {
	k.Status.State = StateDeleting
	condition := metav1.Condition{
		Type:               string(c),
		Status:             "Unknown",
		LastTransitionTime: metav1.Now(),
		Reason:             string(r),
		Message:            msg,
	}
	meta.SetStatusCondition(&k.Status.Conditions, condition)
}

func init() {
	SchemeBuilder.Register(&ApplicationConnector{}, &ApplicationConnectorList{})
}

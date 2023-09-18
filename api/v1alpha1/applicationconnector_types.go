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

type HPASpec struct {
	IsEnabled       bool   `json:"enable"`
	CpuUsagePercent string `json:"cpuUsagePercent"`
	MinReplicas     string `json:"minReplicas"`
	MaxReplicas     string `json:"maxReplicas"`
}

type AppGatewaySpec struct {
	ProxyPort          string  `json:"proxyPort"`
	ProxyPortCompass   string  `json:"proxyPortCompass"`
	ProxyCacheTTL      string  `json:"proxyCacheTTL"`
	ExternalAPIPort    string  `json:"externalAPIPort"`
	ProxyTimeout       string  `json:"proxyTimeout"`
	RequestTimeout     string  `json:"requestTimeout"`
	AppSecretNamespace string  `json:"appSecretNamespace"`
	LogLevel           string  `json:"logLevel"`
	AutoscalingSpec    HPASpec `json:"autoscaling"`
}

type AppConnValidatorSpec struct {
	ProxyPort                string  `json:"proxyPort"`
	ExternalAPIPort          string  `json:"externalAPIPort"`
	EventingPathPrefixV1     string  `json:"eventingPathPrefixV1"`
	EventingPathPrefixV2     string  `json:"eventingPathPrefixV2"`
	EventingPathPrefixEvents string  `json:"eventingPathPrefixEvents"`
	EventingPublisherHost    string  `json:"eventingPublisherHost"` // namespaced name
	EventingDestinationPath  string  `json:"eventingDestinationPath"`
	LogLevel                 string  `json:"logLevel"`
	LogFormat                string  `json:"logFormat"`
	AutoscalingSpec          HPASpec `json:"autoscaling"`
}

type RuntimeAgentSpec struct {
	SkipCompassTLSVerify       bool   `json:"skipCompassTLSVerify"`
	SkipAppsTLSVerify          bool   `json:"skipAppsTLSVerify"`
	SkipDirectorProxyTLSVerify bool   `json:"skipDirectorProxyTLSVerify"`
	QueryLogging               bool   `json:"queryLogging"`
	MetricsLoggingTimeInterval string `json:"metricsLoggingTimeInterval"`
	ControllerSyncPeriod       string `json:"controllerSyncPeriod"`
	MinConfigSyncTime          string `json:"minimalConfigSyncTime"`
	ValidityRenewalThreshold   string `json:"validityRenewalThreshold"`
	ConfigSecretName           string `json:"configSecretName"` // namespaced name
	ClientSecretName           string `json:"clientSecretName"` // namespaced name
	CASecretName               string `json:"CASecretName"`     // namespaced name
	GatewayPort                string `json:"gatewayPort"`
	CentralGatewayServiceUrl   string `json:"centralGatewayServiceUrl"`
	UploadServiceUrl           string `json:"uploadServiceUrl"`
}

// ApplicationConnectorSpec defines the desired state of ApplicationConnector
type ApplicationConnectorSpec struct {
	ApplicationGatewaySpec AppGatewaySpec       `json:"appGateway"`
	AppConValidatorSpec    AppConnValidatorSpec `json:"appConnValidator"`
	RuntimeAgentSpec       RuntimeAgentSpec     `json:"runtimeAgent"`
	DomainName             string               `json:"domainName"`
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

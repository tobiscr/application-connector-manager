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

type HorizontalScalingSpec struct {
	Enabled         bool   `json:"enable,omitempty"`
	CpuUsagePercent string `json:"cpuUsagePercent,omitempty"`
	MinReplicas     string `json:"minReplicas,omitempty"`
	MaxReplicas     string `json:"maxReplicas,omitempty"`
}

type AppGatewaySpec struct {
	ProxyPort          string                `json:"proxyPort,omitempty"`
	ProxyPortCompass   string                `json:"proxyPortCompass,omitempty"`
	ProxyCacheTTL      string                `json:"proxyCacheTTL,omitempty"`
	ExternalAPIPort    string                `json:"externalAPIPort,omitempty"`
	ProxyTimeout       string                `json:"proxyTimeout,omitempty"`
	RequestTimeout     string                `json:"requestTimeout,omitempty"`
	AppSecretNamespace string                `json:"appSecretNamespace,omitempty"`
	LogLevel           string                `json:"logLevel,omitempty"`
	AutoscalingSpec    HorizontalScalingSpec `json:"autoscaling,omitempty"`
}

type AppConnValidatorSpec struct {
	ProxyPort                string                `json:"proxyPort,omitempty"`
	ExternalAPIPort          string                `json:"externalAPIPort,omitempty"`
	AppNamePlaceholder       string                `json:"appNamePlaceholder,omitempty"`
	EventingPathPrefixV1     string                `json:"eventingPathPrefixV1,omitempty"`
	EventingPathPrefixV2     string                `json:"eventingPathPrefixV2,omitempty"`
	EventingPathPrefixEvents string                `json:"eventingPathPrefixEvents,omitempty"`
	EventingPublisherHost    string                `json:"eventingPublisherHost,omitempty"` // namespaced name
	EventingDestinationPath  string                `json:"eventingDestinationPath,omitempty"`
	LogLevel                 string                `json:"logLevel,omitempty"`
	LogFormat                string                `json:"logFormat,omitempty"`
	AutoscalingSpec          HorizontalScalingSpec `json:"autoscaling,omitempty"`
}

type RuntimeAgentSpec struct {
	ConfigSecretName             string          `json:"configSecretName,omitempty"` // namespaced name
	ClientSecretName             string          `json:"clientSecretName,omitempty"` // namespaced name
	CASecretName                 string          `json:"CASecretName,omitempty"`     // namespaced name
	ControllerSyncPeriod         metav1.Duration `json:"controllerSyncPeriod,omitempty"`
	MinConfigSyncTime            string          `json:"minimalConfigSyncTime,omitempty"`
	MetricsLoggingTimeInterval   metav1.Duration `json:"metricsLoggingTimeInterval,omitempty"`
	CertValidityRenewalThreshold string          `json:"certValidityRenewalThreshold,omitempty"`
	GatewayPort                  string          `json:"gatewayPort,omitempty"` // int
	CentralGatewayServiceUrl     string          `json:"centralGatewayServiceUrl,omitempty"`
	RuntimeEventsURL             string          `json:"runtimeEventsURL,omitempty"`
	RuntimeConsoleURL            string          `json:"runtimeConsoleUrl,omitempty"`
	DirectorProxyPort            string          `json:"directorProxyPort,omitempty"` // int
	HealthcheckPort              string          `json:"healthcheckPort,omitempty"`   // int
	SkipCompassTLSVerify         bool            `json:"skipCompassTLSVerify,omitempty"`
	SkipAppsTLSVerify            bool            `json:"skipAppsTLSVerify,omitempty"`
	SkipDirectorProxyTLSVerify   bool            `json:"skipDirectorProxyTLSVerify,omitempty"`
	QueryLogging                 bool            `json:"queryLogging,omitempty"`
}

// ApplicationConnectorSpec contains configuration of ApplicationConnector module and its state

type ApplicationConnectorSpec struct {
	ApplicationGatewaySpec AppGatewaySpec       `json:"appGateway,omitempty"`
	AppConValidatorSpec    AppConnValidatorSpec `json:"appConnValidator,omitempty"`
	RuntimeAgentSpec       RuntimeAgentSpec     `json:"runtimeAgent,omitempty"`
	DomainName             string               `json:"domainName"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ApplicationConnector is the Schema for the applicationconnectors API
type ApplicationConnector struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationConnectorSpec   `json:"spec,omitempty"`
	Status ApplicationConnectorStatus `json:"status,omitempty"`
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

package istio

import (
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GatewaySpec struct {
	Servers  []*Server         `json:"servers,omitempty"`
	Selector map[string]string `json:"selector,omitempty"`
}

type Server struct {
	Port            *Port              `json:"port,omitempty"`
	Bind            string             `json:"bind,omitempty"`
	Hosts           []string           `json:"hosts,omitempty"`
	Tls             *ServerTLSSettings `json:"tls,omitempty"`
	DefaultEndpoint string             `json:"default_endpoint,omitempty"`
	Name            string             `json:"name,omitempty"`
}

type Port struct {
	Number     uint32 `json:"number,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
	Name       string `json:"name,omitempty"`
	TargetPort uint32 `json:"target_port,omitempty"`
}

type Gateway struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              GatewaySpec `json:"spec,omitempty"`
	Status            IstioStatus `json:"status"`
}

type IstioStatus struct {
	Conditions         []*IstioCondition      `json:"conditions,omitempty"`
	ValidationMessages []*AnalysisMessageBase `json:"validation_messages,omitempty"`
	ObservedGeneration int64                  `json:"observed_generation,omitempty"`
}

type AnalysisMessageBase_Level int32

type AnalysisMessageBase struct {
	Type             *AnalysisMessageBase_Type `json:"type,omitempty"`
	Level            AnalysisMessageBase_Level `json:"level,omitempty"`
	DocumentationUrl string                    `json:"documentation_url,omitempty"`
}

type AnalysisMessageBase_Type struct {
	Name string `json:"name,omitempty"`
	Code string `json:"code,omitempty"`
}

type IstioCondition struct {
	Type               string               `json:"type,omitempty"`
	Status             string               `json:"status,omitempty"`
	LastProbeTime      *timestamp.Timestamp `json:"last_probe_time,omitempty"`
	LastTransitionTime *timestamp.Timestamp `json:"last_transition_time,omitempty"`
	Reason             string               `json:"reason,omitempty"`
	Message            string               `json:"message,omitempty"`
}

type ServerTLSSettings_TLSmode uint32

type ServerTLSSettings_TLSProtocol uint32

type ServerTLSSettings struct {
	HttpsRedirect         bool                          `json:"https_redirect,omitempty"`
	Mode                  string                        `json:"mode,omitempty"`
	ServerCertificate     string                        `json:"server_certificate,omitempty"`
	PrivateKey            string                        `json:"private_key,omitempty"`
	CaCertificates        string                        `json:"ca_certificates,omitempty"`
	CredentialName        string                        `json:"credential_name,omitempty"`
	SubjectAltNames       []string                      `json:"subject_alt_names,omitempty"`
	VerifyCertificateSpki []string                      `json:"verify_certificate_spki,omitempty"`
	VerifyCertificateHash []string                      `json:"verify_certificate_hash,omitempty"`
	MinProtocolVersion    ServerTLSSettings_TLSProtocol `json:"min_protocol_version,omitempty"`
	MaxProtocolVersion    ServerTLSSettings_TLSProtocol `json:"max_protocol_version,omitempty"`
	CipherSuites          []string                      `json:"cipher_suites,omitempty"`
}

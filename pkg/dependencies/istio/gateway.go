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
	DefaultEndpoint string             `json:"defaultEndpoint,omitempty"`
	Name            string             `json:"name,omitempty"`
}

type Port struct {
	Number     uint32 `json:"number,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
	Name       string `json:"name,omitempty"`
	TargetPort uint32 `json:"targetPort,omitempty"`
}

type Gateway struct {
	TypeMeta          `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              GatewaySpec `json:"spec,omitempty"`
	Status            IstioStatus `json:"status"`
}

type IstioStatus struct {
	Conditions         []*IstioCondition      `json:"conditions,omitempty"`
	ValidationMessages []*AnalysisMessageBase `json:"validationMessages,omitempty"`
	ObservedGeneration int64                  `json:"observedGeneration,omitempty"`
}

type AnalysisMessageBase_Level int32

type AnalysisMessageBase struct {
	Type             *AnalysisMessageBase_Type `json:"type,omitempty"`
	Level            AnalysisMessageBase_Level `json:"level,omitempty"`
	DocumentationUrl string                    `json:"documentationUrl,omitempty"`
}

type AnalysisMessageBase_Type struct {
	Name string `json:"name,omitempty"`
	Code string `json:"code,omitempty"`
}

type IstioCondition struct {
	Type               string               `json:"type,omitempty"`
	Status             string               `json:"status,omitempty"`
	LastProbeTime      *timestamp.Timestamp `json:"lastProbeTime,omitempty"`
	LastTransitionTime *timestamp.Timestamp `json:"lastTransitionTime,omitempty"`
	Reason             string               `json:"reason,omitempty"`
	Message            string               `json:"message,omitempty"`
}

type ServerTLSSettings_TLSmode uint32

type ServerTLSSettings struct {
	HttpsRedirect         bool     `json:"httpsRedirect,omitempty"`
	Mode                  string   `json:"mode,omitempty"`
	ServerCertificate     string   `json:"serverCertificate,omitempty"`
	PrivateKey            string   `json:"privateKey,omitempty"`
	CaCertificates        string   `json:"caCertificates,omitempty"`
	CredentialName        string   `json:"credentialName,omitempty"`
	SubjectAltNames       []string `json:"subjectAltNames,omitempty"`
	VerifyCertificateSpki []string `json:"verifyCertificateSpki,omitempty"`
	VerifyCertificateHash []string `json:"verifyCertificateHash,omitempty"`
	MinProtocolVersion    string   `json:"minProtocolVersion,omitempty"`
	MaxProtocolVersion    string   `json:"maxProtocolVersion,omitempty"`
	CipherSuites          []string `json:"cipherSuites,omitempty"`
}

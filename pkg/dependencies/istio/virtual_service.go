package istio

import (
	duration "github.com/golang/protobuf/ptypes/duration"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type VirtualService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualServiceSpec `json:"spec,omitempty"`
	Status IstioStatus        `json:"status"`
}

type VirtualServiceSpec struct {
	Hosts    []string     `json:"hosts,omitempty"`
	Gateways []string     `json:"gateways,omitempty"`
	Http     []*HTTPRoute `json:"http,omitempty"`
	Tls      []*TLSRoute  `json:"tls,omitempty"`
	Tcp      []*TCPRoute  `json:"tcp,omitempty"`
	ExportTo []string     `json:"export_to,omitempty"`
}

type L4MatchAttributes struct {
	DestinationSubnets []string `json:"destination_subnets,omitempty"`
	Port               uint32   `json:"port,omitempty"`
	SourceSubnet       string   `json:"source_subnet,omitempty"`
	Gateways           []string `json:"gateways,omitempty"`
	SourceNamespace    string   `json:"source_namespace,omitempty"`
}

type TCPRoute struct {
	Match []*L4MatchAttributes `json:"match,omitempty"`
	Route []*RouteDestination  `json:"route,omitempty"`
}

type TLSRoute struct {
	Match []*TLSMatchAttributes `json:"match,omitempty"`
	Route []*RouteDestination   `json:"route,omitempty"`
}

type RouteDestination struct {
	Destination *Destination `json:"destination,omitempty"`
	Weight      int32        `json:"weight,omitempty"`
}

type TLSMatchAttributes struct {
	SniHosts           []string `json:"sni_hosts,omitempty"`
	DestinationSubnets []string `json:"destination_subnets,omitempty"`
	Port               uint32   `json:"port,omitempty"`
	Gateways           []string `json:"gateways,omitempty"`
	SourceNamespace    string   `json:"source_namespace,omitempty"`
}

type HTTPRoute struct {
	Name             string                  `json:"name,omitempty"`
	Match            []*HTTPMatchRequest     `json:"match,omitempty"`
	Route            []*HTTPRouteDestination `json:"route,omitempty"`
	Redirect         *HTTPRedirect           `json:"redirect,omitempty"`
	DirectResponse   *HTTPDirectResponse     `json:"direct_response,omitempty"`
	Delegate         *Delegate               `json:"delegate,omitempty"`
	Rewrite          *HTTPRewrite            `json:"rewrite,omitempty"`
	Timeout          *duration.Duration      `json:"timeout,omitempty"`
	Retries          *HTTPRetry              `json:"retries,omitempty"`
	Fault            *HTTPFaultInjection     `json:"fault,omitempty"`
	Mirror           *Destination            `json:"mirror,omitempty"`
	Mirrors          []*HTTPMirrorPolicy     `json:"mirrors,omitempty"`
	MirrorPercent    *wrappers.UInt32Value   `json:"mirror_percent,omitempty"`
	MirrorPercentage *Percent                `json:"mirror_percentage,omitempty"`
	CorsPolicy       *CorsPolicy             `json:"cors_policy,omitempty"`
	Headers          *Headers                `json:"headers,omitempty"`
}

type CorsPolicy struct {
	AllowOrigin      []string            `json:"allow_origin,omitempty"`
	AllowOrigins     []*string           `json:"allow_origins,omitempty"`
	AllowMethods     []string            `json:"allow_methods,omitempty"`
	AllowHeaders     []string            `json:"allow_headers,omitempty"`
	ExposeHeaders    []string            `json:"expose_headers,omitempty"`
	MaxAge           *duration.Duration  `json:"max_age,omitempty"`
	AllowCredentials *wrappers.BoolValue `json:"allow_credentials,omitempty"`
}

type HTTPMirrorPolicy struct {
	Destination *Destination `json:"destination,omitempty"`
	Percentage  *Percent     `json:"percentage,omitempty"`
}

type Percent struct {
	Value float64 `json:"value,omitempty"`
}

type HTTPFaultInjection_Delay struct {
	Percent    int32    `json:"percent,omitempty"`
	Percentage *Percent `json:"percentage,omitempty"`
}

type HTTPFaultInjection_Abort struct {
	Percentage *Percent `json:"percentage,omitempty"`
}

type HTTPFaultInjection struct {
	Delay *HTTPFaultInjection_Delay `json:"delay,omitempty"`
	Abort *HTTPFaultInjection_Abort `json:"abort,omitempty"`
}

type HTTPRetry struct {
	Attempts      int32              `json:"attempts,omitempty"`
	PerTryTimeout *duration.Duration `json:"per_try_timeout,omitempty"`
	RetryOn       string             `json:"retry_on,omitempty"`
}

type HTTPRewrite struct {
	Uri             *StringMatch  `json:"uri,omitempty"`
	Authority       *StringMatch  `json:"authority,omitempty"`
	UriRegexRewrite *RegexRewrite `json:"uri_regex_rewrite,omitempty"`
}

type StringMatch struct {
	Exact  string `json:"exact,omitempty"`
	Prefix string `json:"prefix,omitempty"`
	Regex  string `json:"regex,omitempty"`
}

type RegexRewrite struct {
	Match   string `json:"match,omitempty"`
	Rewrite string `json:"rewrite,omitempty"`
}

type Delegate struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type HTTPDirectResponse struct {
	Status uint32  `json:"status,omitempty"`
	Body   *string `json:"body,omitempty"`
}

type HTTPRedirect struct {
	Uri          string `json:"uri,omitempty"`
	Authority    string `json:"authority,omitempty"`
	Scheme       string `json:"scheme,omitempty"`
	RedirectCode uint32 `json:"redirect_code,omitempty"`
}

type HTTPMatchRequest struct {
	Name            string       `json:"name,omitempty"`
	Uri             *StringMatch `json:"uri,omitempty"`
	Scheme          *StringMatch `json:"scheme,omitempty"`
	Method          *StringMatch `json:"method,omitempty"`
	Authority       *StringMatch `json:"authority,omitempty"`
	Port            uint32       `json:"port,omitempty"`
	Gateways        []string     `json:"gateways,omitempty"`
	IgnoreUriCase   bool         `json:"ignore_uri_case,omitempty"`
	SourceNamespace string       `json:"source_namespace,omitempty"`
	StatPrefix      string       `json:"stat_prefix,omitempty"`
}

type HTTPRouteDestination struct {
	Destination *Destination `json:"destination,omitempty"`
	Weight      int32        `json:"weight,omitempty"`
	Headers     *Headers     `json:"headers,omitempty"`
}

type Destination struct {
	Host   string        `json:"host,omitempty"`
	Subset string        `json:"subset,omitempty"`
	Port   *PortSelector `json:"port,omitempty"`
}

type PortSelector struct {
	Number uint32 `json:"number,omitempty"`
}

type Headers struct {
	Request  *Headers_HeaderOperations `json:"request,omitempty"`
	Response *Headers_HeaderOperations `json:"response,omitempty"`
}

type Headers_HeaderOperations struct {
	Remove []string `json:"remove,omitempty"`
}

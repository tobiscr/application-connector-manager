package types

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	VirtualService = schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "VirtualService",
	}
	Gateway = schema.GroupVersionKind{
		Group:   "networking.istio.io",
		Version: "v1beta1",
		Kind:    "Gateway",
	}
	Dependencies = []schema.GroupVersionKind{
		VirtualService,
		Gateway,
	}
)

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	compassv1alpha1 "github.com/kyma-project/kyma/components/compass-runtime-agent/pkg/apis/compass/v1alpha1"
	versioned "github.com/kyma-project/kyma/components/compass-runtime-agent/pkg/client/clientset/versioned"
	internalinterfaces "github.com/kyma-project/kyma/components/compass-runtime-agent/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/kyma-project/kyma/components/compass-runtime-agent/pkg/client/listers/compass/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// CompassConnectionInformer provides access to a shared informer and lister for
// CompassConnections.
type CompassConnectionInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.CompassConnectionLister
}

type compassConnectionInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewCompassConnectionInformer constructs a new informer for CompassConnection type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewCompassConnectionInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredCompassConnectionInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredCompassConnectionInformer constructs a new informer for CompassConnection type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredCompassConnectionInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CompassV1alpha1().CompassConnections().List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CompassV1alpha1().CompassConnections().Watch(context.TODO(), options)
			},
		},
		&compassv1alpha1.CompassConnection{},
		resyncPeriod,
		indexers,
	)
}

func (f *compassConnectionInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredCompassConnectionInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *compassConnectionInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&compassv1alpha1.CompassConnection{}, f.defaultInformer)
}

func (f *compassConnectionInformer) Lister() v1alpha1.CompassConnectionLister {
	return v1alpha1.NewCompassConnectionLister(f.Informer().GetIndexer())
}
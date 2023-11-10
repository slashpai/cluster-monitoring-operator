// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	"context"
	time "time"

	configv1 "github.com/openshift/api/config/v1"
	versioned "github.com/openshift/client-go/config/clientset/versioned"
	internalinterfaces "github.com/openshift/client-go/config/informers/externalversions/internalinterfaces"
	v1 "github.com/openshift/client-go/config/listers/config/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// IngressInformer provides access to a shared informer and lister for
// Ingresses.
type IngressInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.IngressLister
}

type ingressInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// NewIngressInformer constructs a new informer for Ingress type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewIngressInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredIngressInformer(client, resyncPeriod, indexers, nil)
}

// NewFilteredIngressInformer constructs a new informer for Ingress type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredIngressInformer(client versioned.Interface, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ConfigV1().Ingresses().List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ConfigV1().Ingresses().Watch(context.TODO(), options)
			},
		},
		&configv1.Ingress{},
		resyncPeriod,
		indexers,
	)
}

func (f *ingressInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredIngressInformer(client, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *ingressInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&configv1.Ingress{}, f.defaultInformer)
}

func (f *ingressInformer) Lister() v1.IngressLister {
	return v1.NewIngressLister(f.Informer().GetIndexer())
}

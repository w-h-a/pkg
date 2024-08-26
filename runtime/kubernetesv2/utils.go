package kubernetesv2

import (
	"context"

	actionsv1alpha1 "github.com/w-h-a/crd/pkg/apis/eventing/v1alpha1"
	scheme "github.com/w-h-a/crd/pkg/client/clientset/versioned"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func DeploymentsIndexInformer(client kubernetes.Interface, namespace string) cache.SharedIndexInformer {
	deploymentsClient := client.AppsV1().Deployments(namespace)

	deploymentsInformer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return deploymentsClient.List(context.Background(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return deploymentsClient.Watch(context.Background(), options)
			},
		},
		&appsv1.Deployment{},
		0,
		cache.Indexers{},
	)

	return deploymentsInformer
}

func EventSourcesIndexInformer(client scheme.Interface, namespace string) cache.SharedIndexInformer {
	actionsClient := client.EventingV1alpha1().EventSources(namespace)

	eventSourcesInformer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return actionsClient.List(context.Background(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return actionsClient.Watch(context.Background(), options)
			},
		},
		&actionsv1alpha1.EventSource{},
		0,
		cache.Indexers{},
	)

	return eventSourcesInformer
}

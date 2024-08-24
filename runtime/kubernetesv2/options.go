package kubernetesv2

import (
	"context"

	scheme "github.com/w-h-a/crd/pkg/client/clientset/versioned"
	"github.com/w-h-a/pkg/client"
	"github.com/w-h-a/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type k8sConfigKey struct{}
type k8sClientKey struct{}
type actionClientKey struct{}
type grpcClientKey struct{}

func KubernetesRuntimeV2WithK8sConfigKey(cfg *rest.Config) runtime.RuntimeOption {
	return func(o *runtime.RuntimeOptions) {
		o.Context = context.WithValue(o.Context, k8sConfigKey{}, cfg)
	}
}

func GetK8sConfigFromContext(ctx context.Context) (*rest.Config, bool) {
	cfg, ok := ctx.Value(k8sConfigKey{}).(*rest.Config)
	return cfg, ok
}

func KubernetesRuntimeV2WithK8sClient(cs *kubernetes.Clientset) runtime.RuntimeOption {
	return func(o *runtime.RuntimeOptions) {
		o.Context = context.WithValue(o.Context, k8sClientKey{}, cs)
	}
}

func GetK8sClientFromContext(ctx context.Context) (*kubernetes.Clientset, bool) {
	cs, ok := ctx.Value(k8sClientKey{}).(*kubernetes.Clientset)
	return cs, ok
}

func KubernetesRuntimeV2WithActionsClient(ac *scheme.Clientset) runtime.RuntimeOption {
	return func(o *runtime.RuntimeOptions) {
		o.Context = context.WithValue(o.Context, actionClientKey{}, ac)
	}
}

func GetActionsClientFromContext(ctx context.Context) (*scheme.Clientset, bool) {
	ac, ok := ctx.Value(actionClientKey{}).(*scheme.Clientset)
	return ac, ok
}

func KubernetesRuntimeV2WithGrpcClient(c client.Client) runtime.RuntimeOption {
	return func(o *runtime.RuntimeOptions) {
		o.Context = context.WithValue(o.Context, grpcClientKey{}, c)
	}
}

func GetGrpcClientFromContext(ctx context.Context) (client.Client, bool) {
	c, ok := ctx.Value(grpcClientKey{}).(client.Client)
	return c, ok
}

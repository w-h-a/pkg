package kubernetes

import (
	"errors"
	"fmt"

	"github.com/w-h-a/pkg/runtime"
)

const (
	serviceAccountPath = "/var/run/secrets/kubernetes.io/serviceaccount"
)

var (
	ErrReadNamespace = errors.New("failed to read namespace from service account secret")
)

type kubernetesRuntime struct {
	options runtime.RuntimeOptions
}

func (k *kubernetesRuntime) Options() runtime.RuntimeOptions {
	return k.options
}

func (k *kubernetesRuntime) GetServices(opts ...runtime.GetServicesOption) ([]*runtime.Service, error) {
	options := runtime.NewGetOptions(opts...)

	labels := map[string]string{}

	if len(options.Namespace) > 0 {
		labels["namespace"] = options.Namespace
	}

	if len(options.Name) > 0 {
		labels["name"] = options.Name
	}

	if len(options.Version) > 0 {
		labels["version"] = options.Version
	}

	serviceList := &ServiceList{}

	serviceResource := &Resource{
		Kind:  "service",
		Value: serviceList,
	}

	if err := k.get(serviceResource, labels); err != nil {
		return nil, err
	}

	// doing a map first so that we can mutate in other loops
	serviceMap := map[string]*runtime.Service{}

	for _, k8Service := range serviceList.Items {
		namespace := k8Service.Metadata.Namespace

		name := k8Service.Metadata.Labels["name"]

		version := k8Service.Metadata.Labels["version"]

		address := fmt.Sprintf("%s:%d", k8Service.Spec.ClusterIP, k8Service.Spec.Ports[0].Port)

		// TODO: annotations/statuses from deployments/pods => metadata
		service := &runtime.Service{
			Namespace: namespace,
			Version:   version,
			Name:      name,
			Address:   address,
		}

		serviceMap[namespace+name+version] = service
	}

	services := []*runtime.Service{}

	for _, service := range serviceMap {
		services = append(services, service)
	}

	return services, nil
}

func (k *kubernetesRuntime) String() string {
	return "kubernetes"
}

func (k *kubernetesRuntime) get(resource *Resource, labels map[string]string) error {

	return nil
}

func NewRuntime(opts ...runtime.RuntimeOption) runtime.Runtime {
	options := runtime.NewRuntimeOptions(opts...)

	k := &kubernetesRuntime{
		options: options,
	}

	return k
}

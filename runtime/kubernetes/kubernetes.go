package kubernetes

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/w-h-a/pkg/runtime"
	"github.com/w-h-a/pkg/telemetry/log"
)

const (
	serviceAccountPath = "/var/run/secrets/kubernetes.io/serviceaccount"
)

var (
	ErrReadServiceAccount = errors.New("failed to read service account")
	ErrReadNamespace      = errors.New("failed to read namespace from service account")
	ErrResourceNotFound   = errors.New("resource not found")
	ErrDecode             = errors.New("failed to decode response")
	ErrUnknown            = errors.New("unknown error")
)

type kubernetesRuntime struct {
	options runtime.RuntimeOptions
}

func (k *kubernetesRuntime) Options() runtime.RuntimeOptions {
	return k.options
}

func (k *kubernetesRuntime) GetServices(namespace string, opts ...runtime.GetServicesOption) ([]*runtime.Service, error) {
	options := runtime.NewGetServicesOptions(opts...)

	labels := map[string]string{}

	labels["namespace"] = namespace

	if len(options.Name) > 0 {
		labels["name"] = options.Name
	}

	if len(options.Version) > 0 {
		labels["version"] = options.Version
	}

	serviceList := &ServiceList{}

	serviceResource := &Resource{
		Namespace: namespace,
		Kind:      "service",
		Value:     serviceList,
	}

	if err := k.get(serviceResource, labels); err != nil {
		return nil, err
	}

	// doing a map first so that we can mutate in other loops
	serviceMap := map[string]*runtime.Service{}

	for _, k8Service := range serviceList.Items {
		namespace := k8Service.Metadata.Labels["namespace"]

		name := k8Service.Metadata.Labels["name"]

		version := k8Service.Metadata.Labels["version"]

		address := fmt.Sprintf("%s:%d", k8Service.Spec.ClusterIP, k8Service.Spec.Ports[0].Port)

		// TODO: annotations/statuses from deployments/pods => metadata
		service := &runtime.Service{
			Namespace: namespace,
			Name:      name,
			Version:   version,
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

func (k *kubernetesRuntime) IsServicePresent(name, namespace string) bool {
	return false
}

func (k *kubernetesRuntime) CreateService(name, namespace string, labels map[string]string) error {
	return nil
}

func (k *kubernetesRuntime) UpdateDeployment(obj interface{}) {

}

func (k *kubernetesRuntime) Start() error {
	return nil
}

func (k *kubernetesRuntime) Stop() error {
	return nil
}

func (k *kubernetesRuntime) String() string {
	return "kubernetes"
}

func (k *kubernetesRuntime) get(resource *Resource, labels map[string]string) error {
	return newRequest(resource.Namespace, &k.options).
		get().
		setResource(resource.Kind).
		setParams(&params{labelSelector: labels}).
		do().
		decode(resource.Value)
}

func NewRuntime(opts ...runtime.RuntimeOption) runtime.Runtime {
	options := runtime.NewRuntimeOptions(opts...)

	if len(options.Host) == 0 {
		options.Host = "https://" + os.Getenv("KUBERNETES_SERVICE_HOST") + ":" + os.Getenv("KUBERNETES_SERVICE_PORT")
	}

	fileInfo, err := os.Stat(serviceAccountPath)
	if err != nil {
		log.Fatal(err)
	}

	if fileInfo == nil || !fileInfo.IsDir() {
		log.Fatal(ErrReadServiceAccount)
	}

	if options.BearerToken == "" {
		token, err := os.ReadFile(path.Join(serviceAccountPath, "token"))
		if err != nil {
			log.Fatal(err)
		}
		options.BearerToken = string(token)
	}

	cert, err := CertPoolFromFile(path.Join(serviceAccountPath, "ca.crt"))
	if err != nil {
		log.Fatal(err)
	}

	if options.Client == nil {
		c := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig:    &tls.Config{RootCAs: cert},
				DisableCompression: true,
			},
		}
		options.Client = c
	}

	k := &kubernetesRuntime{
		options: options,
	}

	return k
}

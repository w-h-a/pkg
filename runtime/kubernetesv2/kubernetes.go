package kubernetesv2

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"sync"

	scheme "github.com/w-h-a/crd/pkg/client/clientset/versioned"
	"github.com/w-h-a/pkg/client"
	"github.com/w-h-a/pkg/client/grpcclient"
	"github.com/w-h-a/pkg/runtime"
	"github.com/w-h-a/pkg/telemetry/log"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	actionSidecarContainerName    = "action"
	actionSidecarImage            = "wesha/action"
	actionsIDAnnotationKey        = "actions.xyz/id"
	actionsEnabledAnnotationKey   = "actions.xyz/enabled"
	actionsEnabledAnnotationValue = "true"
	actionsProtocolAnnotationKey  = "actions.xyz/protocol"
	HTTPProtocol                  = "http"
	GRPCProtocol                  = "grpc"
	actionSidecarHTTPPort         = 3500
	actionSidecarGRPCPort         = 50001
	actionSidecarHTTPPortName     = "actions-http"
	actionSidecarGRPCPortName     = "actions-grpc"
	apiAddress                    = "http://actions-api.default.svc.cluster.local:80"
	assignerAddress               = "actions-assigner.default.svc.cluster.local"
)

type kubernetesRuntimeV2 struct {
	options              runtime.RuntimeOptions
	kubern               kubernetes.Interface
	action               scheme.Interface
	client               client.Client
	deploymentsInformer  cache.SharedInformer
	eventSourcesInformer cache.SharedInformer
	running              bool
	closed               chan struct{}
	mainLock             sync.RWMutex
	deploymentLock       sync.RWMutex
	eventSourceLock      sync.RWMutex
}

func (k *kubernetesRuntimeV2) Options() runtime.RuntimeOptions {
	return k.options
}

func (k *kubernetesRuntimeV2) GetServices(namespace string, opts ...runtime.GetServicesOption) ([]*runtime.Service, error) {
	// TODO options

	services := []*runtime.Service{}

	ss, err := k.kubern.CoreV1().Services(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, s := range ss.Items {
		service := &runtime.Service{
			Namespace: s.GetNamespace(),
			Name:      s.GetName(),
			Version:   s.ObjectMeta.Labels["version"],
			Address:   fmt.Sprintf("%s:%d", s.Spec.ClusterIP, s.Spec.Ports[0].Port),
		}

		services = append(services, service)
	}

	return services, nil
}

func (k *kubernetesRuntimeV2) IsServicePresent(name, namespace string) bool {
	_, err := k.kubern.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
	return err == nil
}

func (k *kubernetesRuntimeV2) CreateService(name, namespace string, labels map[string]string) error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: map[string]string{actionsEnabledAnnotationKey: actionsEnabledAnnotationValue},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       80,
					TargetPort: intstr.FromInt(actionSidecarHTTPPort),
					Name:       actionSidecarHTTPPortName,
				},
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       int32(actionSidecarGRPCPort),
					TargetPort: intstr.FromInt(actionSidecarGRPCPort),
					Name:       actionSidecarGRPCPortName,
				},
			},
		},
	}

	if _, err := k.kubern.CoreV1().Services(namespace).Create(context.Background(), service, metav1.CreateOptions{}); err != nil {
		return err
	}

	return nil
}

func (k *kubernetesRuntimeV2) UpdateDeployment(obj interface{}) {
	k.deploymentLock.Lock()
	defer k.deploymentLock.Unlock()

	deployment, ok := obj.(*appsv1.Deployment)
	if !ok {
		log.Warn("did not receive an *appsv1.Deployment")
		return
	}

	name := deployment.GetName()

	actionEnabled := false
	for _, c := range deployment.Spec.Template.Spec.Containers {
		if c.Name == actionSidecarContainerName {
			actionEnabled = true
		}
	}

	annotated := false
	if deployment.ObjectMeta.Annotations != nil {
		if val, ok := deployment.ObjectMeta.Annotations[actionsEnabledAnnotationKey]; ok && val == actionsEnabledAnnotationValue {
			annotated = true
		}
	}

	if annotated && actionEnabled {
		log.Infof("notified of action annotated deployment %s", name)
		if err := k.enableActionOnDeployment(deployment); err != nil {
			log.Errorf("failed to enable action on deployment %s: %v", name, err)
		} else {
			log.Infof("action enabled on deployment %s", name)
		}
	} else if !annotated && actionEnabled {
		log.Infof("notified to remove action from deployment %s", name)
		if err := k.removeActionFromDeployment(deployment); err != nil {
			log.Errorf("failed to remove action from deployment %s: %v", name, err)
		} else {
			log.Infof("action removed from deployment %s", name)
		}
	}
}

// func (k *kubernetesRuntimeV2) GetEventSources(opts ...runtime.GetEventSourcesOption) ([]*source.EventSource, error) {
// 	// TODO: options

// 	eventSources := []*source.EventSource{}

// 	ess, err := k.action.EventingV1alpha1().EventSources(metav1.NamespaceAll).List(context.Background(), metav1.ListOptions{})
// 	if err != nil {
// 		log.Errorf("failed to get event sources: %v", err)
// 		return nil, err
// 	}

// 	for _, es := range ess.Items {
// 		eventSource := &source.EventSource{
// 			Name: es.GetName(),
// 			Spec: &source.EventSourceSpec{
// 				Type:           es.Spec.Type,
// 				ConnectionInfo: es.Spec.ConnectionInfo,
// 			},
// 		}

// 		eventSources = append(eventSources, eventSource)
// 	}

// 	return eventSources, nil
// }

// func (k *kubernetesRuntimeV2) UpdateEventSource(obj interface{}) {
// 	k.eventSourceLock.Lock()
// 	defer k.eventSourceLock.Unlock()

// 	eventSource, ok := obj.(*actions_v1alpha1.EventSource)
// 	if !ok {
// 		log.Warn("did not receive an *actions_v1alpha1.EventSource")
// 		return
// 	}

// 	name := eventSource.GetName()

// 	log.Infof("notified about event source %s", name)

// 	if err := k.publishEventSourceToActions(eventSource); err != nil {
// 		log.Errorf("failed to publish event source %s: %v", name, err)
// 	} else {
// 		log.Infof("published event source %s", name)
// 	}
// }

func (k *kubernetesRuntimeV2) Start() error {
	k.mainLock.Lock()
	defer k.mainLock.Unlock()

	if k.running {
		return nil
	}

	k.running = true
	k.closed = make(chan struct{})

	go func() {
		k.deploymentsInformer.Run(k.closed)
	}()

	go func() {
		k.eventSourcesInformer.Run(k.closed)
	}()

	return nil
}

func (k *kubernetesRuntimeV2) Stop() error {
	k.mainLock.Lock()
	defer k.mainLock.Unlock()

	if !k.running {
		return nil
	}

	select {
	case <-k.closed:
		return nil
	default:
		close(k.closed)
		k.running = false
	}

	return nil
}

func (k *kubernetesRuntimeV2) String() string {
	return "kubernetesV2"
}

func (k *kubernetesRuntimeV2) enableActionOnDeployment(deployment *appsv1.Deployment) error {
	appPort := k.getAppPort(deployment.Spec.Template.Spec.Containers)

	appProtocol := k.getAppProtocol(deployment)

	actionName := k.getActionName(deployment)

	sidecar := corev1.Container{
		Name:            actionSidecarContainerName,
		Image:           actionSidecarImage,
		ImagePullPolicy: corev1.PullIfNotPresent,
		Ports: []corev1.ContainerPort{
			{
				ContainerPort: int32(actionSidecarHTTPPort),
				Name:          actionSidecarHTTPPortName,
			},
			{
				ContainerPort: int32(actionSidecarGRPCPort),
				Name:          actionSidecarGRPCPortName,
			},
		},
		Args: []string{"action"},
		Env: []corev1.EnvVar{
			{
				Name:      "HOST_IP",
				ValueFrom: &corev1.EnvVarSource{FieldRef: &corev1.ObjectFieldSelector{FieldPath: "status.podIP"}},
			},
			{
				Name:  "MODE",
				Value: "kubernetes",
			},
			{
				Name:  "ACTION_HTTP_PORT",
				Value: fmt.Sprintf("%v", actionSidecarHTTPPort),
			},
			{
				Name:  "ACTION_GRPC_PORT",
				Value: fmt.Sprintf("%v", actionSidecarGRPCPort),
			},
			{
				Name:  "APP_PORT",
				Value: appPort,
			},
			{
				Name:  "APP_PROTOCOL",
				Value: appProtocol,
			},
			{
				Name:  "ACTION_ID",
				Value: actionName,
			},
			{
				Name:  "API_ADDRESS",
				Value: apiAddress,
			},
			{
				Name:  "ASSIGNER_ADDRESS",
				Value: assignerAddress,
			},
		},
	}

	deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers, sidecar)

	serviceName := fmt.Sprintf("%s-%s", actionName, "action")

	namespace := deployment.GetNamespace()

	isPresent := k.IsServicePresent(serviceName, namespace)
	if isPresent {
		log.Infof("service %s is already present", serviceName)
		return nil
	} else {
		if err := k.CreateService(serviceName, namespace, deployment.Spec.Selector.MatchLabels); err != nil {
			return err
		}
	}

	if err := k.updateDeployment(deployment); err != nil {
		return err
	}

	return nil
}

func (k *kubernetesRuntimeV2) getAppPort(containers []corev1.Container) string {
	for _, container := range containers {
		if (container.Name != actionSidecarHTTPPortName && container.Name != actionSidecarGRPCPortName) && len(container.Ports) > 0 {
			return fmt.Sprint(container.Ports[0].ContainerPort)
		}
	}

	return ""
}

func (k *kubernetesRuntimeV2) getAppProtocol(deployment *appsv1.Deployment) string {
	if val, ok := deployment.ObjectMeta.Annotations[actionsProtocolAnnotationKey]; ok && val != "" {
		if val != HTTPProtocol && val != GRPCProtocol {
			return HTTPProtocol
		}

		return val
	}

	return HTTPProtocol
}

func (k *kubernetesRuntimeV2) getActionName(deployment *appsv1.Deployment) string {
	annotations := deployment.ObjectMeta.GetAnnotations()

	if val, ok := annotations[actionsIDAnnotationKey]; ok && val != "" {
		return val
	}

	return deployment.GetName()
}

func (k *kubernetesRuntimeV2) removeActionFromDeployment(deployment *appsv1.Deployment) error {
	for i := len(deployment.Spec.Template.Spec.Containers) - 1; i >= 0; i-- {
		if deployment.Spec.Template.Spec.Containers[i].Name == actionSidecarContainerName {
			deployment.Spec.Template.Spec.Containers = append(deployment.Spec.Template.Spec.Containers[:i], deployment.Spec.Template.Spec.Containers[i+1:]...)
		}
	}

	if err := k.updateDeployment(deployment); err != nil {
		return err
	}

	return nil
}

func (k *kubernetesRuntimeV2) updateDeployment(deployment *appsv1.Deployment) error {
	if _, err := k.kubern.AppsV1().Deployments(deployment.GetNamespace()).Update(context.Background(), deployment, metav1.UpdateOptions{}); err != nil {
		return err
	}

	return nil
}

// func (k *kubernetesRuntimeV2) publishEventSourceToActions(eventSource *actions_v1alpha1.EventSource) error {
// 	payload := &pb.EventSource{
// 		Name: eventSource.GetName(),
// 		Spec: &pb.EventSourceSpec{
// 			Type: eventSource.Spec.Type,
// 		},
// 	}

// 	bytes, err := json.Marshal(eventSource.Spec.ConnectionInfo)
// 	if err != nil {
// 		return err
// 	}

// 	payload.Spec.ConnectionInfo = &any.Any{Value: bytes}

// 	services, err := k.kubern.CoreV1().Services(metav1.NamespaceAll).List(context.Background(), metav1.ListOptions{
// 		LabelSelector: labels.SelectorFromSet(map[string]string{actionsEnabledAnnotationKey: actionsEnabledAnnotationValue}).String(),
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	for _, service := range services.Items {
// 		name := service.GetName()

// 		log.Infof("updating action pod selected by service %s", name)

// 		endpoints, err := k.kubern.CoreV1().Endpoints(service.GetNamespace()).Get(context.Background(), name, metav1.GetOptions{})
// 		if err != nil {
// 			log.Errorf("failed to get endpoints for service %s: %v", name, err)
// 			continue
// 		}

// 		go k.publish(payload, endpoints)
// 	}

// 	return nil
// }

// func (k *kubernetesRuntimeV2) publish(eventSource *pb.EventSource, endpoints *corev1.Endpoints) {
// 	if endpoints == nil || len(endpoints.Subsets) == 0 {
// 		return
// 	}

// 	for _, addr := range endpoints.Subsets[0].Addresses {
// 		address := fmt.Sprintf("%s:%s", addr.IP, fmt.Sprintf("%v", actionSidecarGRPCPort))

// 		req := k.client.NewRequest(
// 			client.RequestWithMethod("Action.UpdateEventSource"),
// 			client.RequestWithUnmarshaledRequest(
// 				&pb.UpdateEventSourceRequest{
// 					EventSource: eventSource,
// 				},
// 			),
// 		)

// 		rsp := &pb.UpdateEventSourceResponse{}

// 		if err := k.client.Call(context.Background(), req, rsp, client.CallWithAddress(address)); err != nil {
// 			var internal *errorutils.Error
// 			if e, ok := status.FromError(err); ok {
// 				internal = errorutils.ParseError(e.Message())
// 			} else {
// 				internal = errorutils.ParseError(err.Error())
// 			}

// 			if internal.Code == 0 {
// 				internal.Code = 500
// 				internal.Id = "kubernetesV2"
// 				internal.Detail = "error during grpc call: " + internal.Detail
// 			}

// 			log.Errorf("failed to update event source at address %s: %v", address, internal)
// 		}
// 	}
// }

func NewRuntime(opts ...runtime.RuntimeOption) runtime.Runtime {
	options := runtime.NewRuntimeOptions(opts...)

	r := &kubernetesRuntimeV2{
		options:         options,
		mainLock:        sync.RWMutex{},
		deploymentLock:  sync.RWMutex{},
		eventSourceLock: sync.RWMutex{},
	}

	var config *rest.Config

	cfg, ok := GetK8sConfigFromContext(options.Context)
	if ok {
		config = cfg
	} else {
		var kubeconfig *string
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig")
		} else {
			kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig")
		}

		flag.Parse()

		var err error
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Errorf("failed to get the cluster config: %v", err)
			config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
			if err != nil {
				log.Fatalf("failed to build the cluster config: %v", err)
			}
		}
	}

	cs, ok := GetK8sClientFromContext(options.Context)
	if ok {
		r.kubern = cs
	} else {
		var err error
		r.kubern, err = kubernetes.NewForConfig(config)
		if err != nil {
			log.Fatalf("failed to get the k8s client: %v", err)
		}
		KubernetesRuntimeV2WithK8sClient(r.kubern.(*kubernetes.Clientset))
	}

	ac, ok := GetActionsClientFromContext(options.Context)
	if ok {
		r.action = ac
	} else {
		var err error
		r.action, err = scheme.NewForConfig(config)
		if err != nil {
			log.Fatalf("failed to get action client: %v", err)
		}
		KubernetesRuntimeV2WithActionsClient(r.action.(*scheme.Clientset))
	}

	r.deploymentsInformer = DeploymentsIndexInformer(r.kubern, metav1.NamespaceAll)

	r.deploymentsInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: r.UpdateDeployment,
		UpdateFunc: func(_, newObj interface{}) {
			r.UpdateDeployment(newObj)
		},
	})

	// r.eventSourcesInformer = EventSourcesIndexInformer(r.action, metav1.NamespaceAll)

	// r.eventSourcesInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
	// 	AddFunc: r.UpdateEventSource,
	// 	UpdateFunc: func(oldObj, newObj interface{}) {
	// 		r.UpdateEventSource(newObj)
	// 	},
	// })

	c, ok := GetGrpcClientFromContext(options.Context)
	if ok {
		r.client = c
	} else {
		r.client = grpcclient.NewClient()
	}

	return r
}

package kubernetes

type Resource struct {
	Name  string
	Kind  string
	Value interface{}
}

type ServiceList struct {
	Items []Service `json:"items"`
}

type Service struct {
	Metadata *Metadata    `json:"metadata"`
	Spec     *ServiceSpec `json:"spec,omitempty"`
}

type ServiceSpec struct {
	ClusterIP string            `json:"clusterIP"`
	Ports     []ServicePort     `json:"ports,omitempty"`
	Selector  map[string]string `json:"selector,omitempty"`
	Type      string            `json:"type,omitempty"`
}

type ServicePort struct {
	Port     int    `json:"port"`
	Name     string `json:"name,omitempty"`
	Protocol string `json:"protocol,omitempty"`
}

type DeploymentList struct {
	Items []Deployment `json:"items"`
}

type Deployment struct {
	Metadata *Metadata       `json:"metadata"`
	Spec     *DeploymentSpec `json:"spec,omitempty"`
}

type DeploymentSpec struct {
	Replicas int
	Selector *LabelSelector `json:"selector,omitempty"`
	Template *Template
}

type LabelSelector struct {
	MatchLabels map[string]string `json:"matchLabels,omitempty"`
}

type Template struct {
	Metadata *Metadata `json:"metadata,omitempty"`
	PodSpec  *PodSpec  `json:"spec,omitempty"`
}

type PodList struct {
	Items []Pod `json:"items"`
}

type Pod struct {
	Metadata *Metadata  `json:"metadata"`
	Spec     *PodSpec   `json:"spec,omitempty"`
	Status   *PodStatus `json:"status"`
}

type PodSpec struct {
	Containers []Container `json:"container"`
}

type PodStatus struct {
	Phase      string            `json:"phase"`
	Reason     string            `json:"reason"`
	Containers []ContainerStatus `json:"containerStatuses"`
	Conditions []PodCondition    `json:"conditions,omitempty"`
}

type PodCondition struct {
	Type    string `json:"type"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

type Container struct {
	Name    string          `json:"name"`
	Image   string          `json:"image"`
	Command []string        `json:"command,omitempty"`
	Args    []string        `json:"args,omitempty"`
	Env     []EnvVar        `json:"env,omitempty"`
	Ports   []ContainerPort `json:"ports,omitempty"`
}

type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value,omitempty"`
}

type ContainerPort struct {
	ContainerPort int    `json:"containerPort"`
	HostPort      int    `json:"hostPort,omitempty"`
	Name          string `json:"name,omitempty"`
	Protocol      string `json:"protocol,omitempty"`
}

type ContainerStatus struct {
	State ContainerState `json:"state"`
}

type ContainerState struct {
	Running    *Condition `json:"running"`
	Terminated *Condition `json:"terminated"`
	Waiting    *Condition `json:"waiting"`
}

type Condition struct {
	Started string `json:"startedAt,omitempty"`
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

type Metadata struct {
	Namespace   string            `json:"namespace,omitempty"`
	Name        string            `json:"name,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

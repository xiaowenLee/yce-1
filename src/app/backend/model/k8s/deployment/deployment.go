package deployment

// DeploymentList is the list of Deployment
// DeploymentList for Get @/apis/extensions/v1beta1/namespaces/{namespace}/deployments
type DeploymentList struct {
	Kind       string           `json:"kind"`
	ApiVersion string           `json:"apiVersion"`
	Metadata   MetadataListType `json:'metadata"`
	Items      []ItemsListType  `json:"items"`
}

// DeploymentList/Metadata
type MetadataListType struct {
	SelfLink        string `json:"selfLink"`
	ResourceVersion string `json:"resourceVersion"`
}

// DeploymentList/Items
type ItemsListType Deployment

// Deployment is the details of one Deployment
// Deployment for Post @/apis/extensions/v1beta1/namespaces/{namespace}/deployments
type Deployment struct {
	ApiVersion string       `json:"apiVersion"`
	Kind       string       `json:"kind"`
	Metadata   MetadataType `json:"metadata"`
	Spec       SpecType     `json:"spec"`
}

// Deployment/Metadata
type MetadataType struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	CreationTimestamp string            `json:"creationTimestamp,omitempty"`
	Labels            map[string]string `json:"labels"`
	Annotations       map[string]string `json:"annotations,omitempty"`
}

// Deployment/Spec
type SpecType struct {
	Replicas float64   `json:"replicas"`
	Template TemplateS `json:"template"`
}

// Deployment/Spec/Template
type TemplateS struct {
	Metadata MetadataTS `json:"metadata"`
	Spec     SpecTS     `json:"spec"`
}

// Deployment/Spec/Template/Metadata
type MetadataTS struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels"`
}

// Deployment/Spec/Template/Spec
type SpecTS struct {
	Volumes []VolumesSTS `json:"volumes,omitempty"`
	//Containers []ContainersSTSDL `json:"containers"`
	Containers []ContainerType `json:"containers"`
}

// Deployment/Spec/Template/Spec/Volumes
type VolumesSTS struct {
	Name                  string         `json:"name"`
	HostPath              *HostPathVSTS  `json:"hostPath,omitempty"`
	EmptyDir              *EmptyDirVSTS  `json:"emptyDir,omitempty"`
	PersistentVolumeClaim *PvClaimVSTS   `json:"persistentVolumeClaim,omitempty"`
	Rbd                   *RbdVSTS       `json:"rbd,omitempty"`
	ConfigMap             *ConfigMapVSTS `json:"configMap,omitempty"`
}

// Deployment/Spec/Template/Spec/Volumes/HostPath
type HostPathVSTS struct {
	Path string `json:"path"`
}

// Deployment/Spec/Template/Spec/Volumes/EmptyDir
type EmptyDirVSTS struct {
	Medium string `json:"medium"`
}

// Deployment/Spec/Template/Spec/Volumes/PersistentVolumeClaim
type PvClaimVSTS struct {
	ClaimName string `json:"claimName"`
	ReadOnly  bool   `json:"readOnly"`
}

// Deployment/Spec/Template/Spec/Volumes/Rbd
type RbdVSTS struct {
	Monitors  []string        `json:"monitors"`
	Image     string          `json:"image"`
	FsType    string          `json:"fsType"`
	Pool      string          `json:"pool"`
	User      string          `json:"user"`
	Keyring   string          `json:"keyring"`
	SecretRef *SecretRefRVSTS `json:"secretRef"`
	ReadOnly  bool            `json:"readOnly"`
}

// Deployment/Spec/Template/Spec/Volumes/Rbd/SecretRef
type SecretRefRVSTS struct {
	Name string `json: name"`
}

// Deployment/Spec/Template/Spec/Volumes/ConfigMap
type ConfigMapVSTS struct {
	Name  string               `json:"name"`
	Items []ItemsConfigMapVSTS `json:"items"`
}

// Deployment/Spec/Template/Spec/Volumes/ConfigMap/Items
type ItemsConfigMapVSTS struct {
	Key  string `json:"key"`
	Path string `json:"Path"`
}

// Deployment/Spec/Template/Spec/Containers
type ContainersSTS struct {
	Name           string              `json:"name"`
	Image          string              `json:"image"`
	Command        []string            `json:"command,omitempty"`
	Args           []string            `json:"args,omitempty"`
	Ports          []PortCSTS          `json:"ports"`
	Env            []EnvCSTS           `json:"env,omitempty"`
	Resources      *ResourcesCSTS      `json:"resources,omitempty"`
	VolumeMounts   []VolumeMountsCSTS  `json:"volumeMounts,omitempty"`
	LivenessProbe  *LivenessProbeCSTS  `json:"livenessProbe,omitempty"`
	ReadinessProbe *ReadinessProbeCSTS `json:"readinessProbe,omitempty"`
	Lifecycle      *LifecycleCSTS      `json:"lifecycle,omitempty"`
}

type ContainerType ContainersSTS

// Deployment/Spec/Template/Spec/Containers/Ports
type PortCSTS struct {
	Name          string  `json:"name"`
	HostPort      float64 `json:"hostPort"`
	ContainerPort float64 `json:"containerPort"`
	Protocol      string  `json:"protocol"`
	HostIP        string  `json:"hostIP"`
}

// Deployment/Spec/Template/Spec/Containers/Env
type EnvCSTS struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Deployment/Spec/Template/Spec/Containers/Resources
type ResourcesCSTS struct {
	Limits   map[string]string `json:"limits"`
	Requests map[string]string `json:"requests"`
}

// Deployment/Spec/Template/Spec/Containers/VolumeMounts
type VolumeMountsCSTS struct {
	Name      string `json:"name"`
	ReadOnly  bool   `json:"readOnly"`
	MountPath string `json:"mountPath"`
}

// Deployment/Spec/Template/Spec/Containers/LivenessProbe
type LivenessProbeCSTS struct {
	Exec                *ExecLCSTS      `json:"exec,omitempty"`
	HttpGet             *HttpGetLCSTS   `json:"httpGet,omitempty"`
	TcpSocket           *TcpSocketLCSTS `json:"tcpSocket,omitempty"`
	InitialDelaySeconds float64         `json:"initialDelaySeconds"`
	TimeoutSeconds      float64         `json:"timeoutSeconds"`
	PeriodSeconds       float64         `json:"periodSeconds"`
	SuccessThreshold    float64         `json:"successThreshold"`
	FailureThreshold    float64         `json:"failureThreshold"`
}

// Deployment/Spec/Template/Spec/Containers/LivenessProbe/Exec
type ExecLiveProbeCSTS struct {
	Command []string `json:"command"`
}

// Deployment/Spec/Template/Spec/Containers/LivenessProbe/HttpGet
type HttpGetLiveProbeCSTS struct {
	Path        string          `json:"path"`
	Port        float64         `json:"port"`
	Host        string          `json:"host"`
	Scheme      string          `json:"scheme"`
	HttpHeaders []HeadersHLCSTS `json:"httpHeaders"`
}

// Deployment/Spec/Template/Spec/Containers/LivenessProbe/HttpGet/HttpHeaders
type HeadersHLCSTS struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Deployment/Spec/Template/Spec/Containers/LivenessProbe/TcpSocket
type TcpSocketLCSTS struct {
	Port float64 `json:"port"`
}

//Deployment/Spec/Template/Spec/Containers/ReadinessProbe
type ReadinessProbeCSTS struct {
	HttpGet             *HttpGetRCSTS `json:"httpGet,omitempty"`
	InitialDelaySeconds float64       `json:"initialDelaySeconds"`
	TimeoutSeconds      float64       `json:"timeoutSeconds"`
	PeriodSeconds       float64       `json:"periodSeconds"`
	SuccessThreshold    float64       `json:"successThreshold"`
	FailureThreshold    float64       `json:"failureThreshold"`
}

//Deployment/Spec/Template/Spec/Containers/ReadnessProbe/HttpGet
type HttpGetRCSTS struct {
	Path        string          `json:"path"`
	Port        float64         `json:"port"`
	Host        string          `json:"host"`
	Scheme      string          `json:"scheme"`
	HttpHeaders []HeadersHRCSTS `json:"httpHeaders"`
}

//Deployment/Spec/Template/Spec/Containers/ReadnessProbe/HttpGet/HttpHeaders
type HeadersHRCSTS struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

//Deployment/Spec/Template/Spec/Containers/Livecycle
type LifecycleCSTS struct {
	PostStart *PostStartLCSTS `json:"postStart,omitempty"`
	PreStop   *PreStopLCSTS   `json:"preStop,omitempty"`
}

//Deployment/Spec/Template/Spec/Containers/Livecycle/PostStart
type PostStartLCSTS struct {
	Exec *ExecPostStartLCSTS `json:"exec,omitempty"`
}

//Deployment/Spec/Template/Spec/Containers/Livecycle/PostStart/Exec
type ExecPostStartLCSTS struct {
	Command []string `json:"command"`
}

//Deployment/Spec/Template/Spec/Containers/Livecycle/PreStop
type PreStopLCSTS struct {
	Exec *ExecPreStopLCSTS `json:"exec,omitempty"`
}

//Deployment/Spec/Template/Spec/Containers/Livecycle/PreStop/Exec
type ExecPreStopLCSTS struct {
	Command []string `json:"command"`
}

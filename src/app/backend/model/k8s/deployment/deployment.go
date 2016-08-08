package deployment

// DeploymentList for Get @/apis/extensions/v1beta1/namespaces/{namespace}/deployments
type DeploymentList struct {
	Kind       string           `json:"kind"`
	ApiVersion string           `json:"apiVersion"`
	Metadata   MetadataListType `json:'metadata"`
	Items      []ItemsListType  `json:"items"`
}

type MetadataListType struct { // DeploymentList/Metadata
	SelfLink        string `json:"selfLink"`
	ResourceVersion string `json:"resourceVersion"`
}

type ItemsListType Deployment // DeploymentList/Items

// Deployment for Post @/apis/extensions/v1beta1/namespaces/{namespace}/deployments
type Deployment struct {
	ApiVersion string       `json:"apiVersion"`
	Kind       string       `json:"kind"`
	Metadata   MetadataType `json:"metadata"`
	Spec       SpecType     `json:"spec"`
}

type MetadataType struct { // Deployment/Metadata
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	CreationTimestamp string            `json:"creationTimestamp,omitempty"`
	Labels            map[string]string `json:"labels"`
	Annotations       map[string]string `json:"annotations,omitempty"`
}

type SpecType struct { // Deployment/Spec
	Replicas float64      `json:"replicas"`
	Template TemplateSpec `json:"template"`
}

type TemplateSpec struct { // Deployment/Spec/Template
	Metadata MetadataTSDL `json:"metadata"`
	Spec     SpecTSDL     `json:"spec"`
}

type MetadataTSDL struct { // Deployment/Spec/Template/Metadata
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels"`
}

type SpecTSDL struct { // Deployment/Spec/Template/Spec
	Volumes []VolumesSTSDL `json:"volumes,omitempty"`
	//Containers []ContainersSTSDL `json:"containers"`
	Containers []ContainerType `json:"containers"`
}

type VolumesSTSDL struct { // Deployment/Spec/Template/Spec/Volumes
	Name                  string           `json:"name"`
	HostPath              *HostPathVSTSDL  `json:"hostPath,omitempty"`
	EmptyDir              *EmptyDirVSTSDL  `json:"emptyDir,omitempty"`
	PersistentVolumeClaim *PvClaimVSTSDL   `json:"persistentVolumeClaim,omitempty"`
	Rbd                   *RbdVSTSDL       `json:"rbd,omitempty"`
	ConfigMap             *ConfigMapVSTSDL `json:"configMap,omitempty"`
}

type HostPathVSTSDL struct { // Deployment/Spec/Template/Spec/HostPath
	Path string `json:"path"`
}

type EmptyDirVSTSDL struct { // Deployment/Spec/Template/Spec/EmptyDir
	Medium string `json:"medium"`
}

type PvClaimVSTSDL struct { // Deployment/Spec/Template/Spec/PersistentVolumeClaim
	ClaimName string `json:"claimName"`
	ReadOnly  bool   `json:"readOnly"`
}

type RbdVSTSDL struct { // Deployment/Spec/Template/Spec/Rbd
	Monitors  []string          `json:"monitors"`
	Image     string            `json:"image"`
	FsType    string            `json:"fsType"`
	Pool      string            `json:"pool"`
	User      string            `json:"user"`
	Keyring   string            `json:"keyring"`
	SecretRef *SecretRefRVSTSDL `json:"secretRef"`
	ReadOnly  bool              `json:"readOnly"`
}

type SecretRefRVSTSDL struct { // Deployment/Spec/Template/Spec/Rbd/SecretRef
	Name string `json: name"`
}

type ConfigMapVSTSDL struct { // Deployment/Spec/Template/Spec/ConfigMap
	Name  string           `json:"name"`
	Items []ItemsConfigMap `json:"items"`
}

type ItemsConfigMap struct { // Deployment/Spec/Template/Spec/ConfigMap/Items
	Key  string `json:"key"`
	Path string `json:"Path"`
}

type ContainersSTSDL struct { // Deployment/Spec/Template/Spec/Containers
	Name           string                   `json:"name"`
	Image          string                   `json:"image"`
	Command        []string                 `json:"command,omitempty"`
	Args           []string                 `json:"args,omitempty"`
	Ports          []PortContainer          `json:"ports"`
	Env            []EnvContainer           `json:"env,omitempty"`
	Resources      *ResourcesContainer      `json:"resources,omitempty"`
	VolumeMounts   []VolumeMountsContainer  `json:"volumeMounts,omitempty"`
	LivenessProbe  *LivenessProbeContainer  `json:"livenessProbe,omitempty"`
	ReadinessProbe *ReadinessProbeContainer `json:"readinessProbe,omitempty"`
	Lifecycle      *LifecycleContainer      `json:"lifecycle,omitempty"`
}

type ContainerType ContainersSTSDL

type PortContainer struct { // Deployment/Spec/Template/Spec/Containers/Ports
	Name          string  `json:"name"`
	HostPort      float64 `json:"hostPort"`
	ContainerPort float64 `json:"containerPort"`
	Protocol      string  `json:"protocol"`
	HostIP        string  `json:"hostIP"`
}

type EnvContainer struct { // Deployment/Spec/Template/Spec/Containers/Env
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ResourcesContainer struct { // Deployment/Spec/Template/Spec/Containers/Resources
	Limits   map[string]string `json:"limits"`
	Requests map[string]string `json:"requests"`
}

type VolumeMountsContainer struct { // Deployment/Spec/Template/Spec/Containers/VolumeMounts
	Name      string `json:"name"`
	ReadOnly  bool   `json:"readOnly"`
	MountPath string `json:"mountPath"`
}

type LivenessProbeContainer struct { // Deployment/Spec/Template/Spec/Containers/LivenessProbe
	Exec                *ExecLiveProbeType    `json:"exec,omitempty"`
	HttpGet             *HttpGetLiveProbeType `json:"httpGet,omitempty"`
	TcpSocket           *TcpLiveProbeType     `json:"tcpSocket,omitempty"`
	InitialDelaySeconds float64               `json:"initialDelaySeconds"`
	TimeoutSeconds      float64               `json:"timeoutSeconds"`
	PeriodSeconds       float64               `json:"periodSeconds"`
	SuccessThreshold    float64               `json:"successThreshold"`
	FailureThreshold    float64               `json:"failureThreshold"`
}

type ExecLiveProbeType struct { // Deployment/Spec/Template/Spec/Containers/LivenessProbe/Exec
	Command []string `json:"command"`
}

type HttpGetLiveProbeType struct { // Deployment/Spec/Template/Spec/Containers/LivenessProbe/HttpGet
	Path        string           `json:"path"`
	Port        float64          `json:"port"`
	Host        string           `json:"host"`
	Scheme      string           `json:"scheme"`
	HttpHeaders []HeadersGLPType `json:"httpHeaders"`
}

type HeadersGLPType struct { // Deployment/Spec/Template/Spec/Containers/LivenessProbe/HttpGet/HttpHeaders
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TcpLiveProbeType struct { // Deployment/Spec/Template/Spec/Containers/LivenessProbe/TcpSocket
	Port float64 `json:"port"`
}

type ReadinessProbeContainer struct { //Deployment/Spec/Template/Spec/Containers/ReadnessProbe
	HttpGet             *HttpGetReadProbeType `json:"httpGet,omitempty"`
	InitialDelaySeconds float64               `json:"initialDelaySeconds"`
	TimeoutSeconds      float64               `json:"timeoutSeconds"`
	PeriodSeconds       float64               `json:"periodSeconds"`
	SuccessThreshold    float64               `json:"successThreshold"`
	FailureThreshold    float64               `json:"failureThreshold"`
}

type HttpGetReadProbeType struct { //Deployment/Spec/Template/Spec/Containers/ReadnessProbe/HttpGet
	Path        string           `json:"path"`
	Port        float64          `json:"port"`
	Host        string           `json:"host"`
	Scheme      string           `json:"scheme"`
	HttpHeaders []HeadersGRPType `json:"httpHeaders"`
}

type HeadersGRPType struct { //Deployment/Spec/Template/Spec/Containers/ReadnessProbe/HttpGet/HttpHeaders
	Name  string `json:"name"`
	Value string `json:"value"`
}

type LifecycleContainer struct { //Deployment/Spec/Template/Spec/Containers/Livecycle
	PostStart *PostStartLCType `json:"postStart,omitempty"`
	PreStop   *PreStopLCType   `json:"preStop,omitempty"`
}

type PostStartLCType struct { //Deployment/Spec/Template/Spec/Containers/Livecycle/PostStart
	Exec *ExecPSLCType `json:"exec,omitempty"`
}

type ExecPSLCType struct { //Deployment/Spec/Template/Spec/Containers/Livecycle/PostStart/Exec
	Command []string `json:"command"`
}

type PreStopLCType struct { //Deployment/Spec/Template/Spec/Containers/Livecycle/PreStop
	Exec *ExecPrSLCType `json:"exec,omitempty"`
}

type ExecPrSLCType struct { //Deployment/Spec/Template/Spec/Containers/Livecycle/PreStop/Exec
	Command []string `json:"command"`
}

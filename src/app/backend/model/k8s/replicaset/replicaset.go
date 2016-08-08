package replicaset

// ReplicaSetList for Get @/apis/extensions/v1beta1/namespaces/{namespace}/replicasets
type ReplicaSetList struct {
	Kind       string           `json:"kind"`
	ApiVersion string           `json:"apiVersion"`
	Metadata   MetadataListType `json:"metadata"`
	Items      []ItemsListType  `json:"items"`
}

type MetadataListType struct {
	SelfLink        string `json:"selfLink"`
	ResourceVersion string `json:"resourceVersion"`
}

// ReplicaSet for Post @/apis/extensions/v1beta1/namespaces/{namespace}/replicasets
type ItemsListType struct {
	Kind       string       `json:"kind"`
	ApiVersion string       `json:"apiVersion"`
	Metadata   MetadataType `json:"metadata"`
	Spec       SpecType     `json:"spec"`
	Status     StatusType   `json:"status,omitempty"`
}

type ReplicaSet ItemsListType

// Replicaset/Metadata
type MetadataType struct {
	Name                      string            `json:"name"`
	GenerateName              string            `json:"generateName,omitempty"`
	Namespace                 string            `json:"namespace,omitempty"`
	SelfLink                  string            `json:"selfLink,omitempty"`
	Uid                       string            `json:"uid,omitempty"`
	ResourceVersion           string            `json:"resourceVersion,omitempty"`
	Generation                float64           `json:"generation,omitempty"`
	CreationTimestamp         string            `json:"creationTimestamp,omitempty"`
	DeletionTimestamp         string            `json:"deletionTimestamp,omitempty"`
	DeletionGracePeriodSecond float64           `json:"deletionGracePeriodSeconds,omitempty"`
	Labels                    map[string]string `json:"labels"`
	Annotations               map[string]string `json:"annotations,omitempty"`
}

// Replicaset/Spec
type SpecType struct {
	Replicas float64    `json:"replicas"`
	Selector *SelectorS `json:"selector,omitempty"`
	Template TemplateS  `json:"template"`
}

// Replicaset/Spec/Selector
type SelectorS struct {
	MatchLabels      map[string]string `json:"matchLabels,omitempty"`
	MatchExpressions []MatchExpSS      `json:"matchExpressions,omitempty"`
}

// Replicaset/Spec/Selector/MatchExpressions
type MatchExpSS struct {
	Key      string   `json:"key,omitempty"`
	Operator string   `json:"operator,omitempty"`
	Values   []string `json:"values,omitempty"`
}

// Replicaset/Spec/Template
type TemplateS struct {
	Metadata MetadataTS `json:"metadata"`
	Spec     SpecTS     `json:"spec"`
}

// Replicaset/Spec/Template/Metadata
type MetadataTS struct {
	Name                       string            `json:"name,omitempty"`
	GenerateName               string            `json:"generateName,omitempty"`
	Namespace                  string            `json:"namespace,omitempty"`
	SelfLink                   string            `json:"selfLink,omitempty"`
	Uid                        string            `json:"uid,omitempty"`
	ResourceVersion            string            `json:"resourceVersion,omitempty"`
	Generation                 float64           `json:"generation,omitempty"`
	CreationTimestamp          string            `json:"creationTimestamp,omitempty"`
	DeletionTimestamp          string            `json:"deletionTimestamp,omitempty"`
	DeletionGracePeriodSeconds float64           `json:"deletionGracePeriodSeconds,omitempty"`
	Labels                     map[string]string `json:"labels,omitempty"`
	Annotations                map[string]string `json:"annotations,omitempty"`
}

// Replicaset/Spec/Template/Spec
type SpecTS struct {
	Volumes                       []VolumesSTS      `json:"volumes,omitempty"`
	Containers                    []ContainersSTS   `json:"containers"`
	RestartPolicy                 string            `json:"restartPolicy,omitempty"`
	TerminationGracePeriodSeconds float64           `json:"terminationGracePeriodSeconds,omitempty"`
	ActiveDeadlineSeconds         float64           `json:"activeDeadlineSeconds,omitempty"`
	DnsPolicy                     string            `json:"dnsPolicy,omitempty"`
	NodeSelector                  map[string]string `json:"nodeSelector,omitempty"`
	NodeName                      string            `json:"nodeName,omitempty"`
	HostNetwork                   bool              `json:"hostNetwork,omitempty"`
	HostPID                       bool              `json:"hostPID,omitempty"`
	HostIPC                       bool              `json:"hostIPC,omitempty"`
}

// Replicaset/Spec/Template/Spec/Volumes
type VolumesSTS struct {
	Name                  string         `json:"name,omitempty"`
	HostPath              *HostPathVSTS  `json:"hostPath,omitempty"`
	EmptyDir              *EmptyDirVSTS  `json:"emptyDir,omitempty"`
	PersistentVolumeClaim *PVClaimVSTS   `json:"persistentVolumeClaim,omitempty"`
	Rbd                   *RbdVSTS       `json:"rbd,omitempty"`
	ConfigMap             *ConfigMapVSTS `json:"configMap,omitempty"`
}

// Replicaset/Spec/Template/Spec/Volumes/HostPath
type HostPathVSTS struct {
	Path string `json:"path,omitempty"`
}

// Replicaset/Spec/Template/Spec/Volumes/EmptyDir
type EmptyDirVSTS struct {
	Medium string `json:"medium,omitempty"`
}

// Replicaset/Spec/Template/Spec/Volumes/PersistentVolumeClaim
type PVClaimVSTS struct {
	ClaimName string `json:"claimName,omitempty"`
	ReadOnly  bool   `json:"readOnly,omitempty"`
}

// Replicaset/Spec/Template/Spec/Volumes/Rbd
type RbdVSTS struct {
	Monitors  []string        `json:"monitors,omitempty"`
	Image     string          `json:"image,omitempty"`
	FsType    string          `json:"fsType,omitempty"`
	Pool      string          `json:"pool,omitempty"`
	User      string          `json:"user,omitempty"`
	Keyring   string          `json:"keyring,omitempty"`
	SecretRef *SecretRefRVSTS `json:"secretRef,omitempty"`
	ReadOnly  bool            `json:"readOnly,omitempty"`
}

// Replicaset/Spec/Template/Spec/Volumes/Rbd/SecretRef
type SecretRefRVSTS struct {
	Name string `json:"name,omitempty"`
}

// Replicaset/Spec/Template/Spec/Volumes/ConfigMap
type ConfigMapVSTS struct {
	Name  string       `json:"name,omitempty"`
	Items []ItemsCVSTS `json:"items,omitempty"`
}

// Replicaset/Spec/Template/Spec/Volumes/ConfigMap/Items
type ItemsCVSTS struct {
	Key  string `json:"key,omitempty"`
	Path string `json:"path,omitempty"`
}

// Replicaset/Spec/Template/Spec/Container
type ContainerType struct {
	Name                   string              `json:"name,omitempty"`
	Image                  string              `json:"image"`
	Command                []string            `json:"command,omitempty"`
	Args                   []string            `json:"args,omitempty"`
	WorkingDir             string              `json:"workingDir,omitempty"`
	Ports                  []PortsCSTS         `json:"ports,omitempty"`
	Env                    []EnvCSTS           `json:"env,omitempty"`
	Resources              *ResourceCSTS       `json:"resource,omitempty"`
	VolumeMounts           []VolumeMountsCSTS  `json:"volumeMounts,omitempty"`
	LivenessProbe          *LivenessProbeCSTS  `json:"livenessProbe,omitempty"`
	ReadinessProbe         *ReadinessProbeCSTS `json:"readinessProbe,omitempty"`
	LifeCycle              *LifeCycleCSTS      `json:"lifecycle,omitempty"`
	TerminationMessagePath string              `json:"terminationMessagePath,omitempty"`
	ImagePullPolicy        string              `json:"imagePullPolicy,omitempty"`
	Stdin                  bool                `json:"stdin,omitempty"`
	StdinOnce              bool                `json:"stdinOnce,omitempty"`
	Tty                    bool                `json:"tty,omitempty"`
}

type ContainersSTS ContainerType

// Replicaset/Spec/Template/Spec/Container/Ports
type PortsCSTS struct {
	Name          string  `json:"name,omitempty"`
	HostPort      float64 `json:"hostPort,omitempty"`
	ContainerPort float64 `json:"containerPort,omitempty"`
	Protocol      string  `json:"protocol,omitempty"`
	HostIP        string  `json:"hostIP,omitempty"`
}

// Replicaset/Spec/Template/Spec/Container/Env
type EnvCSTS struct {
	Name      string          `json:"name,omitempty"`
	Value     string          `json:"value,omitempty"`
	ValueFrom *ValueFromECSTS `json:"valueFrom,omitempty"`
}

// Replicaset/Spec/Template/Spec/Container/Env/ValueFrom
type ValueFromECSTS struct {
	FieldRef        *FieldRefVECSTS        `json:"fieldRef,omitempty"`
	ConfigMapKeyRef *ConfigMapKeyRefVECSTS `json:"configMapKeyRef,omitempty"`
	SecretKeyRef    *SecretKeyRefVECSTS    `json:"secretKeyRef,omitempty"`
}

// Replicaset/Spec/Template/Spec/Container/Env/ValueFrom/FieldRef
type FieldRefVECSTS struct {
	ApiVersion string `json:"apiVersion,omitempty"`
	FieldPath  string `json:"fieldPath,omitempty"`
}

// Replicaset/Spec/Template/Spec/Container/Env/ValueFrom/ConfigMapKeyRef
type ConfigMapKeyRefVECSTS struct {
	Name string `json:"name,omitempty"`
	Key  string `json:"key,omitempty"`
}

// Replicaset/Spec/Template/Spec/Container/Env/ValueFrom/SecretKeyRef
type SecretKeyRefVECSTS struct {
	Name string `json:"name,omitempty"`
	Key  string `json:"key,omitempty"`
}

// Replicaset/Spec/Template/Spec/Container/Resources
type ResourceCSTS struct {
	Limits   map[string]string `json:"limits,omitempty"`
	Requests map[string]string `json:"requests,omitempty"`
}

// Replicaset/Spec/Template/Spec/Container/VolumeMounts
type VolumeMountsCSTS struct {
	Name      string `json:"name,omitempty"`
	ReadOnly  bool   `json:"readOnly,omitempty"`
	MountPath string `json:"mountPath,omitempty"`
}

// Replicaset/Spec/Template/Spec/Container/LivenessProbe
type LivenessProbeCSTS struct {
	Exec                *ExecLCSTS      `json:"exec,omitempty"`
	HttpGet             *HttpGetLCSTS   `json:"httpGet,omitempty"`
	TcpSocket           *TcpSocketLCSTS `json:"tcpSocket,omitempty"`
	InitialDelaySeConds float64         `json:"initialDelaySeconds,omitempty"`
	TimeoutSeconds      float64         `json:"timeoutSeconds,omitempty"`
	PeriodSeconds       float64         `json:"periodSeconds,omitempty"`
	SuccessThreshold    float64         `json:"successThreshold,omitempty"`
	FailureThreshold    float64         `json:"failureThreshold,omitempty"`
}

// Replicaset/Spec/Template/Spec/Container/LivenessProbe/Exec
type ExecLCSTS struct {
	Command []string `json:"command,omitempty"`
}

// Replicaset/Spec/Template/Spec/Container/LivenessProbe/HttpGet
type HttpGetLCSTS struct {
	Path        string              `json:"path,omitempty"`
	Port        string              `json:"port,omitempty"`
	Host        string              `json:"host,omitempty"`
	Scheme      string              `json:"scheme,omitempty"`
	HttpHeaders []HttpHeadersHLCSTS `json:"httpHeaders,omitempty"`
}

// Replicaset/Spec/Template/Spec/Container/LivenessProbe/HttpGet/HttpHeaders
type HttpHeadersHLCSTS struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// Replicaset/Spec/Template/Spec/Container/LivenessProbe/TcpSocket
type TcpSocketLCSTS struct {
	Port string `json:"port,omitempty"`
}

// Replicaset/Spec/Template/Spec/Container/ReadinessProbe
type ReadinessProbeCSTT struct {
	Exec                *ExecRCSTS      `json:"exec,omitempty"`
	HttpGet             *HttpGetRCSTS   `json:"httpGet,omitempty"`
	TcpSocket           *TcpSocketRCSTS `json:"tcpSocket,omitempty"`
	InitialDelaySeConds float64         `json:"initialDelaySeconds,omitempty"`
	TimeoutSeconds      float64         `json:"timeoutSeconds,omitempty"`
	PeriodSeconds       float64         `json:"periodSeconds,omitempty"`
	SuccessThreshold    float64         `json:"successThreshold,omitempty"`
	FailureThreshold    float64         `json:"failureThreshold,omitempty"`
}

// Replicaset/Spec/Template/Spec/Container/ReadinessProbe/Exec
type ExecRCSTS struct {
	Command []string `json:"command,omitempty"`
}

// Replicaset/Spec/Template/Spec/Container/ReadinessProbe/HttpGet
type HttpGetRCSTS struct {
	Path        string              `json:"path,omitempty"`
	Port        string              `json:"port,omitempty"`
	Host        string              `json:"host,omitempty"`
	Scheme      string              `json:"scheme,omitempty"`
	HttpHeaders []HttpHeadersHRCSTS `json:"httpHeaders,omitempty"`
}

// Replicaset/Spec/Template/Spec/Container/ReadinessProbe/HttpGet/HttpHeaders
type HttpHeadersHRCSTS struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// Replicaset/Spec/Template/Spec/Container/ReadinessProbe/TcpSocket
type TcpSocketRCSTS struct {
	Port string `json:"port,omitempty"`
}

// Replicaset/Spec/Template/Spec/Container/Lifecycle
type LifeCycleCSTS struct {
	PostStart *PostStartLCSTS `json:"postStart",omitempty`
	PreStop   *PreStopLCSTS   `json:"preStop,omitempty"`
}

// Replicaset/Spec/Template/Spec/Container/Lifecycle/PostStart
type PostStartLCSTS struct {
	Exec *ExecPostLCSTS `json:"exec,omitempty"`
}

// Replicaset/Spec/Template/Spec/Container/Lifecycle/PostStart/Exec
type ExecPostLCSTS struct {
	Command []string `json:"command,omitempty"`
}

// Replicaset/Spec/Template/Spec/Container/Lifecycle/PreStop
type PreStopLCSTS struct {
	Exec *ExecPreLCSTS `json:"exec,omitempty"`
}

// Replicaset/Spec/Template/Spec/Container/Lifecycle/PreStop/Exec
type ExecPreLCSTS struct {
	Command []string `json:"command,omitempty"`
}

// Replicaset/Status
type StatusType struct {
	Replicas             float64 `json:"replicas,omitempty"`
	FullyLabeledReplicas float64 `json:"fullyLabeledReplicas,omitempty"`
	ObservedGeneration   float64 `json:"observedGeneration,omitempty"`
}

package pod

// AppInfoType is only used for testing.
type AppInfoType struct {
	Healthz    AppHealthzType    `json:"appHealthz"`
	Name       string            `json:"appName"`
	Label      map[string]string `json:"appLabel"`
	Datacenter []string          `json:"appDatacenter"`
	Replicas   float64           `json:"appReplicas"`
	Worktime   string            `json:"appWorktime"`
}

type AppHealthzType struct {
	PodsAvailable string `json:"podsAvailable"`
}

type AppListType []AppInfoType

// Podlist for Get @/api/v1/namespaces/{namespaces}/pods
type PodList struct {
	Kind       string         `json:"kind"`
	ApiVersion string         `json:"apiVersion"`
	Metadata   MetadataPLType `json:"metadata"`
	Items      []ItemsPLType  `json:"items"`
}

type MetadataPLType struct {
	ResourceVersion string `json:"resourceVersion,omitempty"`
}

// Pod for Post @/api/v1/namespaces/{namespaces}/pods

type ItemsPLType struct {
	Kind       string        `json:"kind,omitempty"`
	ApiVersion string        `json:"apiVersion,omitempty"`
	Metadata   MetadataIType `json:"metadata"`
	Spec       SpecType      `json:"spec"`
	Status     *StatusType   `json:"status,omitempty"`
}

type Pod ItemsPLType

// Pod/Metadata
type MetadataIType struct {
	Name            string `json:"name"`
	GenerateName    string `json:"generateName,omitempty"`
	Namespace       string `json:"namespace,omitempty"`
	SelfLink        string `json:"selfLink,omitempty"`
	Uid             string `json:"uid,omitempty"`
	ResourceVersion string `json:"resourceVersion,omitempty"`
	//Generation float64 `json:"generation"` //float64 string
	//Generation string `json:"generation"` //float64 string
	CreationTimeStamp string            `json:"creationTimeStamp,omitempty"`
	DeletionTimeStamp string            `json:"deletionTimeStamp,omitempty"`
	Labels            map[string]string `json:"labels"`
	Annotations       map[string]string `json:"annotations,omitempty"`
}

// Pod/Spec
type SpecType struct {
	Volumes                       []VolumesS          `json:"volumes,omitempty"`
	Containers                    []ContainerS        `json:"containers"`
	RestartPolicy                 string              `json:"restartPolicy,omitempty"`
	TerminationGracePeriodSeconds float64             `json:"terminationGracePeriodSeconds,omitempty"`
	ActiveDeadlineSeconds         float64             `json:"activeDeadlineSeconds,omitempty"`
	DnsPolicy                     string              `json:"dnsPolicy,omitempty"`
	NodeSelector                  map[string]string   `json:"nodeSelector,omitempty"`
	NodeName                      string              `json:"nodeName,omitempty"`
	HostNetwork                   bool                `json:"hostNetwork,omitempty"`
	HostPID                       bool                `json:"hostPID,omitempty"`
	HostIPC                       bool                `json:"hostIPC,omitempty"`
	ImagePullSecrets              []ImagePullSecretsS `json:"imagePullSecrets,omitempty"`
}

// Pod/Spec/Volumes
type VolumesS struct {
	Name                  string       `json:"name"`
	HostPath              *HostPathVS  `json:"hostPath",omitempty`
	EmptyDir              *EmptyDirVS  `json:"emptyDir",omitempty`
	PersistentVolumeClaim *PvClaimVS   `json:"persistentVolumeClaim,omitempty"`
	Rbd                   *RbdVS       `json:"rbd,omitempty"`
	ConfigMap             *ConfigMapVS `json:"configMap,omitempty"`
}

// Pod/Spec/Volumes/HostPath
type HostPathVS struct {
	Path string `json:"path,omitempty"`
}

// Pod/Spec/Volumes/EmptyDir
type EmptyDirVS struct {
	Medium string `json:"medium,omitempty"`
}

// Pod/Spec/Volumes/PersistentVolumeClaim
type PvClaimVS struct {
	ClaimName string `json:"claimName,omitempty"`
	ReadOnly  string `json:"readOnly,omitempty"`
}

// Pod/Spec/Volumes/Rbd
type RbdVS struct {
	Monitors  []string      `json:"monitors,omitempty"`
	Image     string        `json:"image,omitempty"`
	FsType    string        `json:"fsType,omitempty"`
	Pool      string        `json:"pool,omitempty"`
	User      string        `json:"user,omitempty"`
	Keyring   string        `json:"keyring,omitempty"`
	SecretRef *SecretRefRVS `json:"secretRef,omitempty"`
	ReadOnly  bool          `json:"readOnly,omitempty"`
}

// Pod/Spec/Volumes/SecretRef
type SecretRefRVS struct {
	Name string `json:"name,omitempty"`
}

// Pod/Spec/Volumes/ConfigMap
type ConfigMapVS struct {
	Name  string             `json:"name,omitempty"`
	Items []ItemsConfigMapVS `json:"items,omitempty"`
}

// Pod/Spec/Volumes/ConfigMap/Items
type ItemsConfigMapVS struct {
	Key  string `json:"key",omitempty`
	Path string `json:"path,omitempty"`
}

// Pod/Spec/Containers
type ContainerS struct {
	Name                   string            `json:"name"`
	Image                  string            `json:"image"`
	Command                []string          `json:"command,omitempty"`
	Args                   []string          `json:"args,omitempty"`
	WorkingDir             string            `json:"workingDir,omitempty"`
	Ports                  []PortsCS         `json:"ports,omitempty"`
	Env                    []EnvCS           `json:"env,omitempty"`
	Resources              *ResourcesCS      `json:"resources,omitempty"`
	VolumeMounts           []VolumeMountsCS  `json:"volumeMounts,omitempty"`
	LivenessProbe          *LivenessProbeCS  `json:"livenessProbe,omitempty"`
	ReadinessProbe         *ReadinessProbeCS `json:"readinessProbe,omitempty"`
	Lifecycle              *LifecycleCS      `json:"lifecycle,omitempty"`
	TerminationMesasgePath string            `json:"terminationMessagePath,omitempty"`
	ImagePullPolicy        string            `json:"imagePullPolicy,omitempty"`
	Stdin                  bool              `json:"stdin,omitempty"`
	StdinOnce              bool              `json:"stdinOnce,omitempty"`
	Tty                    bool              `json:"tty,omitempty"`
}

// Pod/Spec/Containers/Ports
type PortsCS struct {
	Name          string       `json:"name,omitempty"`
	HostPort      *HostPortPCS `json:"hostPort,omitempty"`
	ContainerPort float64      `json:"containerPort,omitempty"`
	Protocol      string       `json:"protocol,omitempty"`
	HostIP        string       `json:"hostIP,omitempty"`
}

// Pod/Spec/Containers/Ports/HostPort
type HostPortPCS struct {
}

// Pod/Spec/Contaienrs/Env
type EnvCS struct {
	Name      string        `json:"name",omitempty`
	Value     string        `json:"value",omitempty`
	ValueFrom *ValueFromECS `json:"valueFrom",omitempty`
}

// Pod/Spec/Containers/Env/ValueFrom
type ValueFromECS struct {
	FieldRef        *FieldRefVECS        `json:"fieldRef,omitempty"`
	ConfigMapKeyRef *ConfigMapKeyRefVECS `json:"configMapKeyRef,omitempty"`
	SecretKeyRef    *SecretKeyRefVECS    `json:"secretKeyRef,omitempty"`
}

// Pod/Spec/Containers/Env/ValueFrom/FieldRef
type FieldRefVECS struct {
	ApiVersion string `json:"apiVersion,omitempty"`
	FieldPath  string `json:"fieldPath,omitempty"`
}

// Pod/Spec/Containers/Env/ValueFrom/ConfigMapKeyRef
type ConfigMapKeyRefVECS struct {
	Name string `json:"name,omitempty"`
	Key  string `json:"key,omitempty"`
}

// Pod/Spec/Containers/Env/ValueFrom/SecretKeyRef
type SecretKeyRefVECS struct {
	Name string `json:"name,omitempty"`
	Key  string `json:"key,omitempty"`
}

// Pod/Spec/Containers/Resources
type ResourcesCS struct {
	Limits   map[string]string `json:"limits,omitempty"`
	Requests map[string]string `json:"requests,omitempty"`
}

// Pod/Spec/Containers/VolumeMounts
type VolumeMountsCS struct {
	Name      string `json:"name,omitempty"`
	ReadOnly  bool   `json:"readOnly,omitempty"`
	MountPath string `json:"mountPath,omitempty"`
}

// Pod/Spec/Containers/LivenessProbe
type LivenessProbeCS struct {
	Exec                *ExecLCS      `json:"exec,omitempty"`
	HttpGet             *HttpGetLCS   `json:"httpGet,omitempty"`
	TcpSocket           *TcpSocketLCS `json:"tcpSocket,omitempty"`
	InitialDelaySeconds float64       `json:"initialDelaySeconds,omitempty"`
	TimeoutSeconds      float64       `json:"timeoutSeconds,omitempty"`
	PeriodSeconds       float64       `json:"periodSeconds,omitempty"`
	SuccessThreshold    float64       `json:"successThreshold,omitempty"`
	FailureThreshold    float64       `json:"failureThreshold,omitempty"`
}

// Pod/Spec/Containers/LivenessProbe/Exec
type ExecLCS struct {
	Command []string `json:"command,omitempty"`
}

// Pod/Spec/Containers/LivenessProbe/HttpGet
type HttpGetLCS struct {
	Path        string        `json:"path,omitempty"`
	Port        float64       `json:"port,omitempty"`
	Host        string        `json:"host,omitempty"`
	Scheme      string        `json:"scheme,omitempty"`
	HttpHeaders []HeadersGLCS `json:"httpHeaders,omitempty"`
}

// Pod/Spec/Containers/LivenessProbe/HttpGet/HttpHeaders
type HeadersGLCS struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// Pod/Spec/Containers/LivenessProbe/TcpSocket
type TcpSocketLCS struct {
	Port float64 `json:"port,omitempty"`
}

// Pod/Spec/Containers/ReadinessProbe
type ReadinessProbeCS struct {
	Exec                *ExecRCS      `json:"exec,omitempty"`
	HttpGet             *HttpGetRCS   `json:"httpGet,omitempty"`
	TcpSocket           *TcpSocketRCS `json:"tcpSocket,omitempty"`
	InitialDelaySeconds float64       `json:"initialDelaySeconds,omitempty"`
	TimeoutSeconds      float64       `json:"timeoutSeconds,omitempty"`
	PeriodSeconds       float64       `json:"periodSeconds,omitempty"`
	SuccessThreshold    float64       `json:"successThreshold,omitempty"`
	FailureThreshold    float64       `json:"failureThreshold,omitempty"`
}

// Pod/Spec/Containers/ReadinessProbe/Exec
type ExecRCS struct {
	Command []string `json:"command,omitempty"`
}

// Pod/Spec/Containers/Readiness/HttpGet
type HttpGetRCS struct {
	Path        string        `json:"path,omitempty"`
	Port        float64       `json:"port,omitempty"`
	Host        string        `json:"host,omitempty"`
	Scheme      string        `json:"scheme,omitempty"`
	HttpHeaders []HeadersGRCS `json:"httpHeaders,omitempty"`
}

// Pod/Spec/Containers/Readiness/HttpGet/HttpHeaders
type HeadersGRCS struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// Pod/Spec/Containers/Readiness/TcpSocket
type TcpSocketRCS struct {
	Port float64 `json:"port,omitempty"`
}

// Pod/Spec/Containers/Lifecycles
type LifecycleCS struct {
	PostStart *PostStartLCS `json:"postStart,omitempty"`
	PreStop   *PreStopLCS   `json:"preStop,omitempty"`
}

// Pod/Spec/Containers/Lifecycles/PostStart
type PostStartLCS struct {
	Exec      *ExecPostStartLCS      `json:"exec,omitempty"`
	HttpGet   *HttpGetPostStartLCS   `json:"httpGet,omitempty"`
	TcpSocket *TcpSocketPostStartLCS `json:"tcpPSLCType,omitempty"`
}

// Pod/Spec/Containers/Lifecycles/PostStart/Exec
type ExecPostStartLCS struct {
	Command []string `json:"command,omitempty"`
}

// Pod/Spec/Containers/Lifecycles/PostStart/HttpGet
type HttpGetPostStartLCS struct {
	Path        string                `json:"path,omitempty"`
	Port        float64               `json:"port,omitempty"`
	Host        string                `json:"host,omitempty"`
	Scheme      string                `json:"scheme,omitempty"`
	HttpHeaders []HeadersPostStartLCS `json:"httpHeaders,omitempty"`
}

// Pod/Spec/Containers/Lifecycles/PostStart/HttpGet/HttpHeaders
type HeadersPostStartLCS struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// Pod/Spec/Containers/Lifecycles/PostStart/TcpSocket
type TcpSocketPostStartLCS struct {
	Port float64 `json:"port,omitempty"`
}

// Pod/Spec/Containers/Lifecycles/PreStop
type PreStopLCS struct {
	Exec      *ExecPreStopLCS      `json:"exec,omitempty"`
	HttpGet   *HttpGetPreStopLCS   `json:"httpGet,omitempty"`
	TcpSocket *TcpSocketPreStopLCS `json:"tcpPSLCType,omitempty"`
}

// Pod/Spec/Containers/Lifecycles/PreStop/Exec
type ExecPreStopLCS struct {
	Command []string `json:"command,omitempty"`
}

// Pod/Spec/Containers/Lifecycles/PreStop/HttpGet
type HttpGetPreStopLCS struct {
	Path        string              `json:"path,omitempty"`
	Port        float64             `json:"port,omitempty"`
	Host        string              `json:"host,omitempty"`
	Scheme      string              `json:"scheme,omitempty"`
	HttpHeaders []HeadersPreStopLCS `json:"httpHeaders,omitempty"`
}

// Pod/Spec/Containers/Lifecycles/PreStop/HttpGet/HttpHeaders
type HeadersPreStopLCS struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// Pod/Spec/Containers/Lifecycles/PreStop/TcpSocket
type TcpSocketPreStopLCS struct {
	Port float64 `json:"port,omitempty"`
}

// Pod/Spec/ImagePullSecrets
type ImagePullSecretsS struct {
	Name string `json:"name,omitempty"`
}

// Pod/Status
type StatusType struct {
	Phase             string               `json:"phase,omitempty"`
	Conditions        []ConditionsSt       `json:"conditions,omitempty"`
	Message           string               `json:"message,omitempty"`
	Reason            string               `json:"reason,omitempty"`
	HostIP            string               `json:"hostIP,omitempty"`
	StartTime         string               `json:"startTime,omitempty"`
	ContainerStatuses []ContainerStatuesSt `json:"containerStatuses,omitempty"`
}

// Pod/Status/Conditions
type ConditionsSt struct {
	Type           string `json:"type,omitempty"`
	Status         string `json:"status,omitempty"`
	LastProbeTime  string `json:"lastProbeTime,omitempty"`
	LastTransition string `json:"lastTransition,omitempty"`
	Reason         string `json:"reason,omitempty"`
	Message        string `json:"mesage,omitempty"`
}

// Pod/Status/ContainerStatuses
type ContainerStatuesSt struct {
	Name         string        `json:"name,omitempty"`
	State        *StateCSt     `json:"state,omitempty"`
	LastState    *LastStateCSt `json:"lastState,omitempty"`
	Ready        bool          `json:"ready,omitempty"`
	RestartCount float64       `json:"restartCount,omitempty"`
	Image        string        `json:"image,omitempty"`
	ImageID      string        `json:"imageID,omitempty"`
	ContainerID  string        `json:"containerID,omitempty"`
}

// Pod/Status/ContainerStatuses/State
type StateCSt struct {
	Waiting    *WaitSCSt `json:"waiting,omitempty"`
	Running    *RunSCSt  `json:"running,omitempty"`
	Terminated *TermSCSt `json:"terminated,omitempty"`
}

// Pod/Status/ContainerStatuses/State/Waiting
type WaitSCSt struct {
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

// Pod/Status/ContainerStatuses/Running
type RunSCSt struct {
	StartedAt string `json:"startedAt,omitempty"`
}

// Pod/Status/ContainerStatuses/Terminated
type TermSCSt struct {
	ExitCode    float64 `json:"exitCode,omitempty"`
	Signal      float64 `json:"signal,omitempty"`
	Reason      string  `json:"reason,omitempty"`
	Message     string  `json:"message,omitempty"`
	StartedAt   string  `json:"startedAt,omitempty"`
	FinishedAt  string  `json:"finishedAt,omitempty"`
	ContainerID string  `json:"containerID,omitempty"`
}

// Pod/Status/ContainerStatuses/LastState
type LastStateCSt struct {
	Waiting    *WaitLastStateCSt `json:"waiting,omitempty"`
	Running    *RunLastStateCSt  `json:"running,omitempty"`
	Terminated *TermLastStateCSt `json:"terminated,omitempty"`
}

// Pod/Status/ContainerStatuses/Waiting
type WaitLastStateCSt struct {
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

// Pod/Status/ContainerStatuses/Running
type RunLastStateCSt struct {
	StartedAt string `json:"startedAt,omitempty"`
}

// Pod/Status/ContainerStatuses/Terminated
type TermLastStateCSt struct {
	ExitCode    float64 `json:"exitCode,omitempty"`
	Signal      float64 `json:"signal,omitempty"`
	Reason      string  `json:"reason,omitempty"`
	Message     string  `json:"message,omitempty"`
	StartedAt   string  `json:"startedAt,omitempty"`
	FinishedAt  string  `json:"finishedAt,omitempty"`
	ContainerID string  `json:"containerID,omitempty"`
}

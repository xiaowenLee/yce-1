package pod

type podList struct {
	Kind string `json: "kind"`
	ApiVersion string `json: "apiVersion"`
	Metadata metadataPLType `json: "metadata"`
	Items []itemsPLType	`json: "items"`
}

type metadataPLType struct {
	ResourceVersion string `json: "resourceVersion"`
}

type itemsPLType struct {
	Kind string `json: "kind"`
	ApiVersion string `json: "apiVersion"`
	Metadata metadataIType `json: "metadata"`
	Spec specIType `json: "spec"`
	Status statusIType `json: "status"`	
}

type metadataIType struct {
	Name string `json: "name"`
	GenerateName string `json: "generateName"`
	Namespace string `json: "namespace"`
	SelfLink string `json: "selfLink"`
	Uid string `json: "uid"`
	ResourceVersion string `json: "resourceVersion"`
	Generation genMIType `json: "generation"`
	CreationTimeStamp string `json: "creationTimeStamp"`
	DeletionTimeStamp string `json: "deletionTimeStamp"`
	Labels string `json: "labels"`
	Annotations string `json: "annotations"`
}

type genMIType struct {

}


type specType struct {
	Volumes []volumesSType `json: "volumes"`
	Containers []containerSType `json: "containers"`
	RestartPolicy string `json: "restartPolicy"`
	TerminationGracePeriodSeconds termGPSType `json: "terminationGracePeriodSeconds"`
	ActiveDeadlineSeconds actDSType `json: "activeDeadlineSeconds"`
	DnsPolicy string `json: "dnsPolicy"`
	NodeSelector string `json: "nodeSelector"`
	NodeName string `json: "nodeName"`
	HostNetwork bool `json: "hostNetwork"`
	HostPID bool `json: "hostPID"`
	HostIPC bool `json: "hostIPC"`
	ImagePullSecrets []imagePullSecretsType `json: "imagePullSecrets"`

}

type volumesSType struct {
	Name string `json: "name"`
	HostPath hostPathVSType `json: "hostPath"`
	EmptyDir emptyDirVSType `json: "emptyDir"`
	PersistentVolumeClaim pvClaimVSType `json: "persistentVolumeClaim"`
	Rbd rbdVSType `json: "rbd"`
	ConfigMap configMapVSType `json: "configMap"`
}

type hostPathVSType struct {
	Path string `json: "path"`
}

type emptyDirVSType struct {
	Medium string `json: "medium"`
}

type pvClaimVSType struct {
	ClaimName string `json: "claimName"`
	ReadOnly string `json: "readOnly"`
}

type rbdVSType struct {
	Monitors []string `json: "monitors"`
	Image string `json: "image"`
	FsType string `json: "fsType"`
	Pool string `json: "pool"`
	User string `json: "user"`
	Keyring string `json: "keyring"`
	SecretRef secretRefrbdType `json: "secretRef"`
	ReadOnly bool `json: "readOnly"`
}

type secretRefrbdType struct {
	Name string `json: "name"`
}

type configMapVSType struct {
	Name string `json: "name"`
	Items []itemsConfigMapType `json: "items"`
}

type itemsConfigMapType struct {
	Key string `json: "key"`
	Path string `json: "path"`
}

type containerSType struct {
	Name string `json: "name"`
	Image string `json: "image"`
	Command []string `json: "command"`
	Args []string `json: "args"`
	WorkingDir string `json: "workingDir"`
	Ports []portsContainerType `json: "ports"
	Env []envContainerType `json: "env"`
	Resources resourcesContainerType `json: "resources"`
	VolumeMounts volumeMountsContainerType `json: "volumeMounts"`
	LivenessProbe livenessProbeContainerType `json: "livenessProbe"`
	ReadinessProbe readinessProbeContainerType `json: "readinessProbe"`
	Lifecycle lifecycleContainerType `json: "lifecycle"`
	TerminationMesasgePath termsmgPathContainerType `json: "terminationMessagePath"`
	ImagePullPolicy imagePullPolicyContainerType `json: "imagePullPolicy"`
	Stdin bool `json: "stdin"`
	StdinOnce bool `json: "stdinOnce"`
	Tty bool `json: "tty"`
}

type portsContainerType struct {
	Name string `json: "name"`
	HostPort hostPortContainerType `json: "hostPort"`
	ContainerPort cPortContainerType `json: "containerPort"`
	Protocol string `json: "protocol"`
	HostIP string `json: "hostIP"`
}

type hostPortContainerType struct {

}

type cPortContainerType struct {

}

type envContainerType struct {
	Name string `json: "name"`
	Value string `json: "value"`
	ValueFrom valueFromEnvContainer `json: "valueFrom"`
}

type valueFromEnvContainer struct {
	FieldRef fieldRefValueFromEnvCon `json: "fieldRef"`
	ConfigMapKeyRef configMapKeyEnvCon `json: "configMapKeyRef"`
	SecretKeyRef secretKeyRefEnvCon  `json: "secretKeyRef"`
}

type fieldRefValueFromEnvCon struct {
	ApiVersion string `json: "apiVersion"`
	FieldPath string `json: "fieldPath"`
}

type configMapKeyEnvCon struct {
	Name string `json: "name"`
	Key string `json: "key"`
}

type secretKeyRefEnvCon struct {
	Name string `json: "name"
	Key string `json: "key"`
}

type resourcesContainerType struct {
    LimitsRsConType string `json: "limits"
    Requests string `json: "requests"`
}

type volumeMountsContainerType struct {
    Name string `json: "name"`
    ReadOnly bool `json: "readOnly"`
    MountPath string `json: "mountPath"`
}

type livenessProbeContainerType struct {
    Exec execLiveProbeType `json: "exec"`
    HttpGet httpGetLiveProbeType `json: "httpGet"`
    TcpSocket tcpLiveProbeType `json: "tcpSocket"`
    InitialDelaySeconds initLiveProbeType `json: "initialDelaySeconds"`
    TimeoutSeconds timeoutLiveProbeType `json: "timeoutSeconds"`
    PeriodSeconds prdLiveProbeType `json: "periodSeconds"`
    SuccessThreshold succthLiveProbeType `json: "successThreshold"`
    FailureThreshold failthLiveProbeType `json: "failureThreshold"`
}

type execLiveProbeType struct {
    Command []string `json: "command"` 
}

type httpGetLiveProbeType struct {
    Path string `json: "path"`
    Port string `json: "port"`
    Host string `json: "host"`
    Scheme string `json: "scheme"`
    HttpHeaders headersGLPType `json: "httpHeaders"`    
}

type headersGLPType struct {
    Name string `json: "name"`
    Value string `json: "value"`
}

type tcpLiveProbeType struct {
    Port string `json: "port"`
}

type initLiveProbeType struct {

}

type timeoutLiveProbeType struct {

}

type prdLiveProbeType struct {

}

type succthLiveProbeType struct {

}

type failthLiveProbeType struct {

}


type readinessProbeContainerType struct {
    Exec execReadProbeType `json: "exec"`
    HttpGet httpGetReadProbeType `json: "httpGet"`
    TcpSocket tcpReadProbeType `json: "tcpSocket"`
    InitialDelaySeconds initReadProbeType `json: "initialDelaySeconds"`
    TimeoutSeconds timeoutReadProbeType `json: "timeoutSeconds"`
    PeriodSeconds prdReadProbeType `json: "periodSeconds"`
    SuccessThreshold succthReadProbeType `json: "successThreshold"`
    FailureThreshold failthReadProbeType `json: "failureThreshold"`

}

type execReadProbeType struct {
    Command []string `json: "command"` 
}

type httpGetReadProbeType struct {
    Path string `json: "path"`
    Port string `json: "port"`
    Host string `json: "host"`
    Scheme string `json: "scheme"`
    HttpHeaders headersGRPType `json: "httpHeaders"`    
}

type headersGRPType struct {
    Name string `json: "name"`
    Value string `json: "value"`
}

type tcpReadProbeType struct {
    Port string `json: "port"`
}

type initReadProbeType struct {

}

type timeoutReadProbeType struct {

}

type prdReadProbeType struct {

}

type succthReadProbeType struct {

}

type failthReadProbeType struct {

}

type termsmgPathContainerType struct {
    
}

type imagePullPolicyContainerType struct {

}


type termGPSType struct {

}

type actDSType struct {

}


type imagePullSecretsType struct {
    Name string `json: "name"`
}

type statusType struct {
	Phase string `json: "phase"`
	Conditions []conditionsStatType `json: "conditions"`
	Message string `json: "message"`
	Reason string `json: "reason"`
	HostIP string `json: "hostIP"`
	StartTime string `json: "startTime"`
	ContainerStatuses []containerStatuesStatType `json: "containerStatuses"`
}

type conditionsStatType struct {
    Type string `json: "type"`
    Status string `json: "status"`
    LastProbeTime string `json: "lastProbeTime"`
    LastTransition string `json: "lastTransition"`
    Reason string `json: "reason"`
    Message string `json: "mesage"`
} 

type containerStatuesStatType struct {
    Name string `json: "name"`
    State statCSSType `json: "state"`
    LastState lastStateCSSType `json: "lastState"`
    Ready bool `json: "ready"`
    RestartCount restartCountCSSType `json: "restartCount"`
    Image string `json: "image"`
    ImageID string `json: "imageID"`
    ContainerID string `json: "containerID"`
}

type statCSSType struct {
    Waiting waitSCSSType `json: "waiting"`
    Running runSCSSType `json: "running"`
    Terminated termCSSType `json: "terminated"`
}

type waitSCSSType struct {
    Reason string `json: "reason"`
    Message string `json: "message"` 
}

type runSCSSType struct {
    StartedAt string `json: "startedAt"`
}

type termCSSType struct {
    ExitCode exitcodeTCSSType `json: "exitCode"`
    Signal sigTCSSType `json: "signal"`
    Reason string `json: "reason"`
    Message string `json: "message"`
    StartedAt string `json: "startedAt"`
    FinishedAt string `json: "finishedAt"`
    ContainerID string `json: "containerID"`
}

type exitcodeTCSSType struct {

}

type sigTCSSType struct {

}

type lastStateCSSType struct {
    Waiting waitLSCSSType `json: "waiting"`
    Running runLSCSSType `json: "running"`
    Terminated termLSCSSType `json: "terminated"`
}

type waitLSCSSType struct {
    Reason string `json: "reason"`
    Message string `json: "message"`
}

type runLSCSSType struct {
    StartedAt string `json: "startedAt"`
}

type termLSCSSType struct {
    ExitCode exitcodeTLSCSSType `json: "exitCode"`
    Signal sigTLSCSSType `json: "signal"`
    Reason string `json: "reason"`
    Message string `json: "message"`
    StartedAt string `json: "startedAt"`
    FinishedAt string `json: "finishedAt"`
    ContainerID string `json: "containerID"`
}

type exitcodeTLSCSSType struct {
    
}

type sigTLSCSSType struct {
      
}
   
type restartCountCSSType struct {

}









func podlist() {

}

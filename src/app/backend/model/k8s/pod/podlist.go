package pod

type PodList struct {
	Kind string `json: "kind"`
	ApiVersion string `json: "apiVersion"`
	Metadata metadataPLType `json: "metadata"`
	Items []itemsPLType	`json: "items"`
}

type metadataPLType struct {
	SelfLink string `json: "selfLink"`
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
	Generation generationMI `json: "generation"`
	CreationTimeStamp string `json: "creationTimeStamp"`
	DeletionTimeStamp string `json: "deletionTimeStamp"`
	DeletionGracePeriodSeconds deletionGPSMIType `json: "deletionGracePeriodSeconds"` //this long name how it is named ?
	Labels string `json: "labels"`
	Annotations string `json: "annotations"`
}

type generationMI struct {

}

type deletionGPSMIType struct {

}

type specType struct {
	Volumes []volumesSType `json: "volumes"`
	Containers []containerSType `json: "containers"`
	RestartPolicy string `json: "restartPolicy"`
	TerminationGracePeriodSeconds terminationGPSType `json: "terminationGracePeriodSeconds"`
	ActiveDeadlineSeconds activeDeadlineSType `json: "activeDeadlineSeconds"`
	DnsPolicy string `json: "dnsPolicy"`
	NodeSelector string `json: "nodeSelector"`
	ServiceAccountName string `json: "serviceAccountName"`
	ServiceAccount string `json: "serviceAccount"`
	NodeName string `json: "nodeName"`
	HostNetwork bool `json: "hostNetwork"`
	HostPID bool `json: "hostPID"`
	HostIPC bool `json: "hostIPC"`
	SecurityContext securityContextSType `json: "securityContext"`
	ImagePullSecrets []imagePullSecretsType `json: "imagePullSecrets"`

}

type volumesSType struct {
	Name string `json: "name"`
	HostPath hostPathVSType `json: "hostPath"`
	EmptyDir emptyDirVSType `json: "emptyDir"`
	GcePersistentDisk gceDiskVSType `json: "gcePersistentDisk"`
	AwsElasticBlockStore awsEBVSType `json: "awsElasticBlockStore"`
	GitRepo gitRepoVSType `json: "gitRepo"`
	Secret secretVSType `json: "secret"`
	Nfs nfsVSType `json: "nfs"`
	Iscsi iscsiVSType `json: "iscsi"`
	GlusterFS glusterfsVSType `json: "glusterfs"`
	PersistentVolumeClaim pvClaimVSType `json: "persistentVolumeClaim"`
	Rbd rbdVSType `json: "rbd"`
	FlexVolume flexVolumeVSType `json: "flexVolume"`
	Cinder cinderVSType `json: "cinder"`
	Cephfs cephfsVSType `json: "cephfs"`
	Flocker flockerVSType `json: "flocker"`
	DownwardAPI downwardAPIVSType `json: "downwardAPI"`
	Fc fcVSType `json: "fc"`
	AzureFile azureFileVSType `json: "azurefile"`
	ConfigMap configMapVSType `json: "configMap"`
}

type hostPathVSType struct {
	Path string `json: "path"`
}

type emptyDirVSType struct {
	Medium string `json: "medium"`
}

type GcePersistentDisk struct {
	PdName string `json: "pdName"`
	FsType string `json: "fsType"`
	Partition partitionGCEType `json: "partition"`
	ReadOnly bool `json: "readOnly"`
}

type partitionGCEType struct {

}

type awsEBVSType struct {
	VolumeID string `json: "volumeID"`
	FsType string `json: "fsType"`
	Partition partitionAWSType `json: "partition"`
	ReadOnly bool `json: "readOnly"`
}

type partitionAWSType struct {

}

type gitRepoVSType struct {
	Repository string `json: "repository"`
	Revision string `json: "revision"`
	Directory string `json: "directory"`
}

type secretVSType struct {
	SecretName string `json: "secretName"`
}

type nfsVSType struct {
	Server string `json: "server"`
	Path string `json: "path"`
	ReadOnly bool `json: "readOnly"`
} 

type iscsiVSType struct {
	TargetPortal string `json: "targetPortal"`
	Iqn string `json: "iqn"`
	Lun lunISCSIType `json: "lun"`
	IscsiInterface string `json: "iscsiInterface"`
	FsType string `json: "fsType"`
	ReadOnly bool `json: "readOnly"`
}

type glusterfsVSType struct {
	Endpoints string `json: "endppints"`
	Path string `json: "path"`
	ReadOnly bool `json: "readOnly"`
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

type flexVolumeVSType struct {
	Driver string `json: "driver"`
	FsType string `json: "fsType"`
	SecretRef secretRefflexVolumeType `json: "secretRef"`
	ReadOnly bool `json: "readOnly"`
	Options string `json: "options"`
}

type secretRefflexVolumeType struct {
	Name string `json: "name"`
}

type cinderVSType struct {
	VolumeID string `json: "volumeID"`
	FsType string `json: "fsType"`
	ReadOnly bool `json: "readOnly"`
}

type cephfsVSType struct {
	Monitors []string `json: "monitors"
	Path string `json: "path"`
	User string `json: "user"`
	SecretFile string `json: "secretFile"`
	SecretRef secretRefcephfsType `json: "secretRef"`
	ReadOnly bool `json: "readOnly"`
}

type secretRefcephfs struct {
	Name string `json: "name"`
}

type flockerVSType struct {
	DatasetName string `json: "datasetName"`
}

type downwardAPIVSType struct {
	Items []itemsDownTYpe `json: "items"`
}

type itemsDownTYpe struct {
	Path string `json: "path"`
	FieldRef fieldRefItemDownwardType `json: "fieldRef"`
}

type fieldRefItemDownwardType struct {
	ApiVersion string `json: "apiVersion"`
	FieldPath string `json: "fieldPath"`
}

type fcVSType struct {
	TargetWWNs []targetfcType `json: "targetWWNs"`
	Lun lunfcType `json: "lun"`
	FsType string `json: "fsType"`
	ReadOnly bool `json: "readOnly"`
}

type targetfcType struct {
	
}

type azureFileVSType struct {
	SecretName string `json: "secretName"`
	ShareName string `json: "shareName"`
	ReadOnly bool `json: "readOnly"`
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
	SecurityContext secontextContainerType `json: "securityContext"`
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
	SecretKeyRef secretKeyRefEnvCon  `json: "secretKeyRefEnvCon"`
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

-------------------------------------------------------------------

type resourcesContainerType struct {

}

type terminationGPSType struct {

}

type activeDeadlineSType struct {

}

type securityContextSType struct {

}

type imagePullSecretsType struct {

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

func podlist() {

}

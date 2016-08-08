package service

// ServiceList for Get @/api/v1/namespaces/{namespace}/services
type ServiceList struct {
	Kind       string           `json:"kind"`
	ApiVersion string           `json:"apiVersion"`
	Metadata   MetadataListType `json:"metadata"`
	Items      []ItemsListType  `json:"items"`
}

type MetadataListType struct {
	SelfLink        string `json:"selfLink"`
	ResourceVersion string `json:"resourceVersion"`
}

// Service for Post @/api/v1/namespaces/{namespace}/services
type ItemsListType struct {
	Kind       string       `json:"kind"`
	ApiVersion string       `json:"apiVersion"`
	Metadata   MetadataType `json:"metadata"`
	Spec       SpecType     `json:"spec"`
	Status     StatusType   `json:"status,omitempty"`
}

type Service ItemsListType

// Service/Metadata
type MetadataType struct {
	Name                       string            `json:"name"`
	GenerateName               string            `json:generateName,omitempty"`
	Namespace                  string            `json:"namespace,omitempty"`
	SelfLink                   string            `json:"selfLink,omitempty"`
	Uid                        string            `json:uid,omitempty"`
	ResourceVersion            string            `json:"resourceVersion,omitempty"`
	Generation                 float64           `json:"generation,omitempty"`
	CreationTimestamp          string            `json:"creationTimestamp,omitempty"`
	DeletionTimestamp          string            `json:"deletionTimestamp,omitempty"`
	DeletionGracePeriodSeconds float64           `json:"deletionGracePeriodSeconds,omitempty"`
	Labels                     map[string]string `json:"labels"`
	Annotations                map[string]string `json:"annotations,omitempty"`
}

// Service/Spec
type SpecType struct {
	Ports               []PortsS          `json:"ports"`
	Selector            map[string]string `json:"selector,omitempty"`
	ClusterIP           string            `json:"clusterIP,omitempty"`
	Type                string            `json:"type,omitempty"`
	ExternalIPs         []string          `json:"externalIPs,omitempty"`
	DeprecatedPublicIPs []string          `json:"deprecatedPublicIPs,omitempty"`
	SessionAffinity     string            `json:"sessionAffinity,omitempty"`
	LoadBalancerIP      string            `json:"loadBalancerIP,omitempty"`
}

// Service/Spec/Ports
type PortsS struct {
	Name       string  `json:"name,omitempty"`
	Protocol   string  `json:"protocol"`
	Port       float64 `json:"port"`
	TargetPort float64 `json:"targetPort,omitempty"`
	NodePort   float64 `json:"nodePort,omitempty"`
}

type PortsType PortsS

// Service/Status
type StatusType struct {
	LoadBalancer *LBStatusSt `json:"loadBalancer,omitempty"`
}

// Service/Status/LoadBalancer
type LBStatusSt struct {
	Ingress []IngressLBStatusSt `json:"ingress,omitempty"`
}

// Service/Status/LoadBalancer/Ingress
type IngressLBStatusSt struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
}

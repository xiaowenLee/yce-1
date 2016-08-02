package namespace

// NamespaceList for Get /api/v1/namespaces
type NamespaceList struct {
    Kind string `json:"kind"`
    ApiVersion string `json:"apiVersion"`
    Metadata MetadataListType `json:"metadata"`
    Items []ItemsListType `json:"items"`
}

type MetadataListType struct {
    SelfLink string `json:"selfLink"`
    ResourceVersion string `json:"resourceVersion"`
}

type ItemsListType Namespace 

// Namespace for Post /api/v1/namespace
type Namespace struct {
    Kind string `json:"kind"`
    ApiVersion string `json:"apiVersion"`
    Metadata MetadataType `json:"metadata"`
    Spec SpecType `json:"spec"`
    Status StatusType `json:"status"`
}

type MetadataType struct {
    Name string `json:"name"`
    GenerateName string `json:"generateName"`
    Namespace string `json:"namespace"`
    SelfLink string `json:"selfLink"`
    Uid string `json:"uid"`
    ResourceVersion string `json:"resourceVersion"`
    Generation float64 `json:"generation"`
    CreationTimestamp string `json:"creationTimestamp"`
    DeletionTimestamp string `json:"deletionTimestamp"`
    DeletionGracePeriodSeconds float64 `json:"deletionGracePeriodSeconds"`
    Labels map[string] string `json:"labels"`
    Annotations map[string] string `json:"annotations"`
}

type SpecType struct {
//   Finalizers []FinalizersType `json:finalizers"`
    Finalizers []string `json:"finalizers"`
}
/*
type FinalizersType struct {

}
*/
type StatusType struct {
    Phase string `json:"phase"`
}

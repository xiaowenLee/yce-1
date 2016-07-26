package main

type deploymentType struct {
    Kind string `json: "kind"`
    ApiVersion string `json: "apiVersion"`
    Metadata metadataDPType `json: "metadata"`
    Spec specDPType `json: "spec"`
    Status statusDPType `json: "status"`
}

type metadataDPType struct {
    Name string `json: "name"`
    GenerateName string `json: "generateName"`
    Namespace string `json: "namespace"`
    SelfLink string `json: "selfLink"`
    Uid string `json: "uid"`
    ResourceVersion string `json: "resourceVersion"`
    Generation genMDPType `json: "generation"`
    CreationTimestamp string `json: "creationTimestamp"`
    DeletionTimestamp string `json: "deletionTimestamp"`
    DeletionGracePeriodSeconds delGPSMDPType `json: "deletionGracePeriodSeconds"`
    Labels string `json: "labels"`
    Annotations: string `json: "annotations"` 
}

type genMDPType struct {
    
}

type delGPSMDPType struct {

} 

type specDPType struct {
    Replicas replicasType `json: "replicas"`
    Selector selectorType `json: "selector"`
    Template templateType `json: "template"`
    Strategy strategyType `json: "strategy"`
    MinReadySeconds mrsSDPType `json: "minReadySeconds"`
    RevisionHistoryLimit rhlSDPType `json: "revisionHistoryLimit"`
    Paused bool `json: "paused"`
    RollbackTo rbSDPType `json: "rollbackTo"`
}

----------------------------------------------------------

type statusDPType struct {
    
}


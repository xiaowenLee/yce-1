package deploy

// 1. 先找到Deployment: 如何通过Pod/RS找到他的deployment或者deployment的名字
// 2. client.Deployments(namespace).Get(name)得到这个Deployment
// 3. deployment.Spec.RollbackTo RollbackConfig类型

// 4. client.Deployments(namespace).Rollback() 参数是*extensions.DeploymentRollback
/*
	type DeploymentRollback struct {
		unversioned.TypeMeta `json:",inline"`
		// Required: This must match the Name of a deployment.
		Name string `json:"name"`
		// The annotations to be updated to a deployment
		UpdatedAnnotations map[string]string `json:"updatedAnnotations,omitempty"`
		// The config of this deployment rollback.
		RollbackTo RollbackConfig `json:"rollbackTo"`
	}

	type RollbackConfig struct {
	    // The revision to rollback to. If set to 0, rollbck to the last revision.
	    Revision int64 `json:"revision,omitempty"`
	}
*/


// API: https://godoc.org/k8s.io/kubernetes/pkg/client/unversioned#DeploymentInterface


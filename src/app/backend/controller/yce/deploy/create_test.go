package deploy_test

import (
	testclient "app/backend/common/yce/testclient"
	"testing"
	"k8s.io/kubernetes/pkg/api"
	"app/backend/model/yce/deploy"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"encoding/json"
	"fmt"
	"k8s.io/kubernetes/pkg/api/unversioned"
	mylog "app/backend/common/util/log"
)

func TestDeploymentCreate(t *testing.T) {
	// the testing of Init API is a must
	//initUrl := "/api/v1/organizations/1/users/1/deployments/new"

	// solution 1:
	deployment := &deploy.CreateDeployment{
		AppName: "foo-nginx",
		OrgName: "ops",
		DcIdList: []int32{1},
		Deployment: extensions.Deployment{
			TypeMeta: unversioned.TypeMeta{
				Kind: "Deployment",
				APIVersion: "extensions/v1beta1",
			},
			ObjectMeta: api.ObjectMeta{
				Name: "foo-nginx",
				Namespace: "ops",
			},
			Spec: extensions.DeploymentSpec{
				Replicas: 1,
				Template: api.PodTemplateSpec{
					ObjectMeta: api.ObjectMeta{
						Labels: map[string]string{
							"app": "foo-nginx",
						},
					},
					Spec: api.PodSpec{
						Containers: []api.Container{
							api.Container{
								Name: "foo-nginx",
								Image: "nginx:1.7.9",
							},
						},
					},
				},
			},
		},
	}

	// solution 2:
	// expectString := string()

	headers := make(map[string]string, 1)
	headers["Authorization"] = "0e45e2b2-c195-4518-b1f0-46aa895a7aa0"

	//mylog.Log.Debugf("Headers: %v", headers)

	newUrl := "http://localhost:8080/api/v1/organizations/1/users/1/deployments/new"
	body, err := json.Marshal(deployment)

	//mylog.Log.Debugf("Json Marshal: body=%s", string(body))

	if err != nil {
		mylog.Log.Debugf("Json Marshal: error=%s", err)
		return
	}

	c := testclient.TestClient{
		Request: testclient.Request{
			Header: headers,
			Path: newUrl,
			Body: body,
		},
	}
	c.Post()

	actual := string(c.Response.Body)
	fmt.Printf("Deployment Create actual resule: \n%s\n", actual)

	//TODO: using validate() without human interference
}



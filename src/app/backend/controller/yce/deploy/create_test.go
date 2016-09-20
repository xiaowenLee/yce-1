package deploy_test

import (
	testclient "app/backend/common/yce/testclient"
	"testing"
	"k8s.io/kubernetes/pkg/api"
	"app/backend/model/yce/deploy"
	"github.com/kubernetes/kubernetes/pkg/apis/extensions"
	"encoding/json"
	"net/http"
	"fmt"
	"github.com/kubernetes/kubernetes/pkg/api/unversioned"
)

func TestDeploymentCreate(t *testing.T) {
	// the testing of Init API is a must
	//initUrl := "/api/v1/organizations/1/users/1/deployments/new"

	// solution 1:
	deployment := &deploy.CreateDeployment{
		AppName: "foo-nginx",
		OrgName: "ops",
		DcIdList: []int32{1},
		Deployment: &extensions.Deployment{
			TypeMeta: unversioned.TypeMeta{

			},
			ObjectMeta: api.ObjectMeta{
				Name: "foo-nginx",
				Namespace: "ops",
			},
			Spec: extensions.DeploymentSpec{

			},
			//TODO: the other must
		},
	}

	// solution 2:
	// expectString := string()

	headers := make(map[string]string, 0)
	headers["Authorization"] = "0e45e2b2-c195-4518-b1f0-46aa895a7aa0"

	newUrl := "/api/v1/organizations/1/users/1/deployments/new"
	body, _ := json.Marshal(deployment)

	c := testclient.TestClient{
		Request: testclient.Request{
			Header: headers,
			Path: newUrl,
			Body: body,
		},
	}
	c.Get()

	actual := string(c.Response.Body)
	fmt.Println(actual)

	//TODO: using validate() without human interference
}



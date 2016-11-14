package template

import (
	"testing"
	testclient "app/backend/common/yce/testclient"
	"fmt"
)


func Test_List(t *testing.T) {
	header := make(map[string]string)
	header["Authorization"] = testclient.SessionId

	req := &testclient.Request{
		Header: header,
		Path: testclient.LocalServer + "/api/v1/organizations/" + testclient.OrgId + "/users/" + testclient.UserId +
			"/templates",
	}

	resp := new(testclient.Response)

	testClient := &testclient.TestClient{
		Request: *req,
		Response: *resp,
	}

	testClient.Get()

	respString := string([]byte(testClient.Response.Body))

	fmt.Println(respString)

}

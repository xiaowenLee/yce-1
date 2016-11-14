package operationstat

import (
	"testing"
	testclient "app/backend/common/yce/testclient"
	"fmt"
)


func Test_Operationstat(t *testing.T) {
	header := make(map[string]string)
	header["Authorization"] = testclient.SessionId

	req := &testclient.Request{
		Header: header,
		Path: testclient.LocalServer + "/api/v1/organizations/" + testclient.OrgId + "/operationstat",
	}

	resp := new(testclient.Response)

	testClient := &testclient.TestClient{
		Request: *req,
		Response: *resp,
	}
	/*
	testClient.Get()

	respString := string([]byte(testClient.Response.Body))

	if !testClient.Validate(respString, testclient.OperationStat) {
		t.Errorf("Operationstat test failed")
		t.Failed()
	}
	*/

	result := testClient.Get()

	if result.GetCode() == 0 {
		fmt.Printf("OK\nData: %s\n", result.GetData())
	} else {
		fmt.Printf("\nCode: %d\nMessage: %s\nData: %s\n", result.GetCode(), result.GetMessage(), result.GetData())
		t.Fail()
	}
}
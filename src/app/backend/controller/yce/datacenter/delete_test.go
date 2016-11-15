package datacenter

import (
	"testing"
	testclient "app/backend/common/yce/testclient"
	"fmt"
)


func Test_Delete(t *testing.T) {
	header := make(map[string]string)
	header["Authorization"] = testclient.SessionId

	req := &testclient.Request{
		Header: header,
		Path: testclient.LocalServer + "/api/v1/datacenter/delete",
		Body: []byte(testclient.DELETE_DATACENTER),
	}

	resp := new(testclient.Response)

	testClient := &testclient.TestClient{
		Request: *req,
		Response: *resp,
	}

	result := testClient.Post()

	if result.GetCode() == 0 {
		fmt.Println("OK")
	} else {
		fmt.Printf("\nCode: %d\nMessage: %s\nData: %s\n", result.GetCode(), result.GetMessage(), result.GetData())
		t.Fail()
	}


}

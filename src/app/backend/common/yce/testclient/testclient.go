package testclient

import (
	//myhttpclient "app/backend/common/util/http"
	//"k8s.io/kubernetes/pkg/runtime"
	mylog "app/backend/common/util/log"
	myerror "app/backend/common/yce/error"
	mysql "app/backend/common/util/mysql"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	config "app/backend/common/yce/config"
	"github.com/kubernetes/kubernetes/pkg/util/json"
)

type Request struct {
	Header map[string]string
	Path   string
	//Body io.Reader
	Body []byte
	r    io.Reader
	//Body runtime.Object
	//RawBody *string
}

type Response struct {
	StatusCode int
	Body       []byte
	//Body       runtime.Object
	//RawBody    *string
}



type TestClient struct {
	//myhttpclient.HttpClient
	http.Client
	Request  Request
	Response Response
}

func Instance() *TestClient {
	return new(TestClient)
}

func (t *TestClient) Validate(expect, actual string) bool {
	return strings.EqualFold(expect, actual)
}

func (t *TestClient) Get() *myerror.YceError {

	t.Request.r = strings.NewReader(string(t.Request.Body))

	mylog.Log.Debugf("Get t.Request.r: %s\n", t.Request.r)

	req, err := http.NewRequest("GET", t.Request.Path, t.Request.r)
	if err != nil {
		mylog.Log.Errorf("Get error=%s", err)
		return nil
	}

	for k, v := range t.Request.Header {
		req.Header.Add(k, v)
	}

	resp, err := t.Do(req)
	if err != nil {
		mylog.Log.Errorf("Get error=%s", err)
		return nil
	}

	defer resp.Body.Close()

	//t.Response.StatusCode = http.StatusOK
	t.Response.Body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		mylog.Log.Errorf("Get error=%s", err)
		return nil
	}

	ye := new(myerror.YceError)
	err = json.Unmarshal(t.Response.Body, ye)
	if err != nil {
		mylog.Log.Errorf("Get error=%s", err)
		return nil
	}

	return ye

}

func (t *TestClient) Post() *myerror.YceError {
	t.Request.r = strings.NewReader(string(t.Request.Body))
	req, err := http.NewRequest("POST", t.Request.Path, t.Request.r)
	if err != nil {
		mylog.Log.Errorf("Post error=%s", err)
		return nil
	}

	for k, v := range t.Request.Header {
		req.Header.Add(k, v)
	}

	resp, err := t.Do(req)
	if err != nil {
		mylog.Log.Errorf("Post error=%s", err)
		return nil
	}

	defer resp.Body.Close()

	t.Response.Body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		mylog.Log.Errorf("Post error=%s", err)
		return nil
	}

	ye := new(myerror.YceError)
	err = json.Unmarshal(t.Response.Body, ye)
	if err != nil {
		mylog.Log.Errorf("Get error=%s", err)
		return nil
	}

	return ye
}

func (t *TestClient) Delete() *myerror.YceError {
	t.Request.r = strings.NewReader(string(t.Request.Body))
	req, err := http.NewRequest("DELETE", t.Request.Path, t.Request.r)
	if err != nil {
		mylog.Log.Errorf("Delete error=%s", err)
		return nil
	}

	for k, v := range t.Request.Header {
		req.Header.Add(k, v)
	}

	resp, err := t.Do(req)
	if err != nil {
		mylog.Log.Errorf("Delete error=%s", err)
		return nil
	}
	defer resp.Body.Close()

	t.Response.Body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		mylog.Log.Errorf("Delete error=%s", err)
		return nil
	}

	ye := new(myerror.YceError)
	err = json.Unmarshal(t.Response.Body, ye)
	if err != nil {
		mylog.Log.Errorf("Get error=%s", err)
		return nil
	}

	return ye
}


func (t *TestClient)ConnectDB() *myerror.YceError {
	config.Instance().Load()
	mysql.MysqlInstance().Open()
	return nil
}


/*
import (
	"testing"
	simple "k8s.io/kubernetes/pkg/client/unversioned/testclient/simple"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/api/testapi"
)

type TestClient struct {
	simple.Client
}
*/

/*
func getDeploymentsResourceName() string {
	return "deployments"
}

func TestDeploymentCreate(t *testing.T) {
	ns := api.NamespaceDefault
	deployment := extensions.Deployment{
		ObjectMeta: api.ObjectMeta{
			Name:      "abc",
			Namespace: ns,
		},
	}
	c := &simple.Client{
		Request: simple.Request{
			Method: "POST",
			Path:   testapi.Extensions.ResourcePath(getDeploymentsResourceName(), ns, ""),
			Query:  simple.BuildQueryValues(nil),
			Body:   &deployment,
		},
		Response: simple.Response{StatusCode: 200, Body: &deployment},
	}

	response, err := c.Setup(t).Deployments(ns).Create(&deployment)
	defer c.Close()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c.Validate(t, response, err)

}
*/

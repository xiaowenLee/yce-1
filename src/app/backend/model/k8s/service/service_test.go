package service

import (
	hc "app/backend/common/util/http"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
)

func Test_Get(t *testing.T) {

	url := "http://master:8080/api/v1/namespaces/default/services"
	var svclist ServiceList
	client := hc.NewHttpClient("", "")

	response, err := client.Get(url)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(response, &svclist)
	if err != nil {
		log.Println(err)
	}

	num := len(svclist.Items)
	fmt.Println(num)

	for i := 0; i < num; i++ {
		fmt.Println(svclist.Items[i].Metadata.Name)
	}

}

func Test_Post(t *testing.T) {

	url := "http://master:8080/api/v1/namespaces/default/services"
	client := hc.NewHttpClient("", "")

	labels := make(map[string]string, 1)
	labels["name"] = "nginx-svc-test"

	svc := new(Service)
	svc.ApiVersion = "v1"
	svc.Kind = "Service"
	svc.Metadata.Name = "nginx-svc-test"
	svc.Metadata.Labels = labels
	svc.Spec.Ports = make([]PortsS, 1)
	svc.Spec.Ports[0].Protocol = "TCP"
	svc.Spec.Ports[0].Port = 8080
	svc.Spec.Ports[0].TargetPort = 3360
	svc.Spec.ExternalIPs = make([]string, 1)
	svc.Spec.ExternalIPs[0] = "10.112.21.45"

	result, _ := json.MarshalIndent(svc, "", " ")
	fmt.Println(string(result))
	rep, err := client.Post(url, strings.NewReader(string(result)))
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(string(rep))
	}

}
func Test_Delete(t *testing.T) {
	url := "http://master:8080/api/v1/namespaces/default/services/nginx-svc-test"
	client := hc.NewHttpClient("", "")
	resp, err := client.Delete(url)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(string(resp))
	}
}

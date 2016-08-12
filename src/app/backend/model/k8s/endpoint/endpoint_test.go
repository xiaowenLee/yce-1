package endpoint

import (
	hc "app/backend/common/util/http"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
)

func Test_Get(t *testing.T) {
	url := "http://master:8080/api/v1/namespaces/default/endpoints"
	var eplist EndpointList
	client := hc.NewHttpClient("", "")

	response, err := client.Get(url)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(response, &eplist)
	if err != nil {
		log.Println(err)
	}

	num := len(eplist.Items)
	fmt.Println(num)

	for i := 0; i < num; i++ {
		fmt.Println(eplist.Items[i].Metadata.Name)
	}

}

func Test_Post(t *testing.T) {
	url := "http://master:8080/api/v1/namespaces/default/endpoints"
	client := hc.NewHttpClient("", "")

	labels := make(map[string]string, 1)
	labels["name"] = "nginx-ep-test"

	ep := new(Endpoint)
	ep.ApiVersion = "v1"
	ep.Kind = "Endpoints"
	ep.Metadata.Name = "nginx-ep-test"
	ep.Metadata.Labels = labels
	ep.SubSets = make([]SubSetsType, 1)
	ep.SubSets[0].Addresses = make([]AddressesType, 1)
	ep.SubSets[0].Addresses[0].IP = "192.168.1.100"
	ep.SubSets[0].Ports = make([]PortsType, 1)
	ep.SubSets[0].Ports[0].Port = 8080

	result, _ := json.MarshalIndent(ep, "", " ")
	fmt.Println(string(result))
	rep, err := client.Post(url, strings.NewReader(string(result)))
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(string(rep))
	}

}
func Test_Delete(t *testing.T) {
	url := "http://master:8080/api/v1/namespaces/default/endpoints/nginx-ep-test"
	client := hc.NewHttpClient("", "")
	resp, err := client.Delete(url)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(string(resp))
	}
}

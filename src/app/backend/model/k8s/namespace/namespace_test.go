package namespace

import (
	hc "app/backend/common/util/http"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
)

func Test_Get(t *testing.T) {
	url := "http://master:8080/api/v1/namespaces"
	var nslist NamespaceList
	client := hc.NewHttpClient("", "")

	response, err := client.Get(url)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(response, &nslist)
	if err != nil {
		log.Println(err)
	}

	num := len(nslist.Items)
	fmt.Println(num)

	for i := 0; i < num; i++ {
		fmt.Println(nslist.Items[i].Metadata.Name)
	}
}

func Test_Post(t *testing.T) {
	url := "http://master:8080/api/v1/namespaces"
	client := hc.NewHttpClient("", "")

	labels := make(map[string]string, 1)
	labels["name"] = "nginx-ns-test"

	ns := new(Namespace)
	ns.ApiVersion = "v1"
	ns.Kind = "Namespace"
	ns.Metadata.Name = "nginx-ns-test"
	ns.Metadata.Labels = labels

	result, _ := json.MarshalIndent(ns, "", "")
	fmt.Println(string(result))
	rep, err := client.Post(url, strings.NewReader(string(result)))
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(string(rep))
	}
}

func Test_Delete(t *testing.T) {
	url := "http://master:8080/api/v1/namespaces/nginx-ns-test"
	client := hc.NewHttpClient("", "")
	resp, err := client.Delete(url)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(string(resp))
	}
}

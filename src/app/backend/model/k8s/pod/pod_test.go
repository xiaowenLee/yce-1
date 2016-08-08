package pod

import (
	hc "app/backend/common/util/http"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
)

func Test_Get(t *testing.T) {
	url := "http://master:8080/api/v1/namespaces/default/pods"
	client := hc.NewHttpClient("", "")
	var plist PodList
	response, err := client.Get(url)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(response, &plist)
	if err != nil {
		log.Println(err)
	}

	num := len(plist.Items)
	fmt.Println(num)

	for i := 0; i < num; i++ {
		fmt.Println(plist.Items[i].Metadata.Name)
	}
}

func Test_Post(t *testing.T) {
	url := "http://master:8080/api/v1/namespaces/default/pods"

	client := hc.NewHttpClient("", "")

	labels := make(map[string]string, 1)
	labels["name"] = "nginx-pd-test"

	pd := new(Pod)
	pd.ApiVersion = "v1"
	pd.Kind = "Pod"
	pd.Metadata.Name = "nginx-pd-test"
	pd.Metadata.Labels = labels
	pd.Spec.Containers = make([]ContainerS, 1)
	pd.Spec.Containers[0].Name = "nginx-pd-test"
	pd.Spec.Containers[0].Image = "nginx:1.7.9"

	result, _ := json.MarshalIndent(pd, "", " ")
	fmt.Println(string(result))
	rep, err := client.Post(url, strings.NewReader(string(result)))
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(string(rep))
	}
}

func Test_Delete(t *testing.T) {
	url := "http://master:8080/api/v1/namespaces/default/pods/nginx-pd-test"
	client := hc.NewHttpClient("", "")
	resp, err := client.Delete(url)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(string(resp))
	}
}

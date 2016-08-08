package deployment

import (
	hc "app/backend/common/util/http"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
)

func Test_Get(t *testing.T) {
	url := "http://master:8080/apis/extensions/v1beta1/namespaces/default/deployments"
	var dplist DeploymentList
	client := hc.NewHttpClient("", "")

	response, err := client.Get(url)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(response, &dplist)
	if err != nil {
		log.Println(err)
	}

	num := len(dplist.Items)
	fmt.Println(num)

	for i := 0; i < num; i++ {
		fmt.Println(dplist.Items[i].Metadata.Name)
		fmt.Println(dplist.Items[i].Metadata.CreationTimestamp)
	}
}

func Test_Post(t *testing.T) {
	url := "http://master:8080/apis/extensions/v1beta1/namespaces/default/deployments"
	client := hc.NewHttpClient("", "")

	labels := make(map[string]string, 1)
	labels["name"] = "nginx-test"

	dp := new(Deployment)
	dp.ApiVersion = "extensions/v1beta1"
	dp.Kind = "Deployment"
	dp.Metadata.Name = "nginx-test"
	dp.Spec.Replicas = 3
	dp.Spec.Template.Metadata.Labels = labels
	dp.Spec.Template.Spec.Containers = make([]ContainerType, 1)
	dp.Spec.Template.Spec.Containers[0].Name = "nginx-test"
	dp.Spec.Template.Spec.Containers[0].Image = "nginx:1.7.9"

	result, _ := json.MarshalIndent(dp, "", " ")
	rep, err := client.Post(url, strings.NewReader(string(result)))
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(string(rep))
	}
}

func Test_Delete(t *testing.T) {
	url := "http://master:8080/apis/extensions/v1beta1/namespaces/default/deployments/nginx-test"
	client := hc.NewHttpClient("", "")
	resp, err := client.Delete(url)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(string(resp))
	}
}

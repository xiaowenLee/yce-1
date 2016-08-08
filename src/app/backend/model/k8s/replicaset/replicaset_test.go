package replicaset

import (
	hc "app/backend/common/util/http"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"testing"
)

func Test_Get(t *testing.T) {
	url := "http://master:8080/apis/extensions/v1beta1/namespaces/default/replicasets"

	var rclist ReplicaSetList
	client := hc.NewHttpClient("", "")

	response, err := client.Get(url)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal(response, &rclist)
	if err != nil {
		log.Println(err)
	}

	num := len(rclist.Items)
	fmt.Println(num)

	for i := 0; i < num; i++ {
		fmt.Println(rclist.Items[i].Metadata.Name)
	}
}

func Test_Post(t *testing.T) {
	url := "http://master:8080/apis/extensions/v1beta1/namespaces/default/replicasets"

	client := hc.NewHttpClient("", "")

	labels := make(map[string]string, 1)
	labels["name"] = "nginx-rs-test"

	rp := new(ReplicaSet)
	rp.ApiVersion = "extensions/v1beta1"
	rp.Kind = "ReplicaSet"
	rp.Metadata.Name = "nginx-rs-test"
	rp.Spec.Replicas = 3
	rp.Spec.Template.Metadata.Labels = labels
	rp.Spec.Template.Spec.Containers = make([]ContainerType, 1)
	rp.Spec.Template.Spec.Containers[0].Name = "nginx-rs-test"
	rp.Spec.Template.Spec.Containers[0].Image = "nginx:1.7.9"

	result, _ := json.MarshalIndent(rp, "", " ")
	rep, err := client.Post(url, strings.NewReader(string(result)))
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(string(rep))
	}
}

func Test_Delete(t *testing.T) {

	url := "http://master:8080/apis/extensions/v1beta1/namespaces/default/replicasets/nginx-rs-test"

	client := hc.NewHttpClient("", "")
	resp, err := client.Delete(url)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(string(resp))
	}
}

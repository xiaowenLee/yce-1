package registry

import (
	"log"
	"io/ioutil"
	"encoding/json"
	"github.com/kataras/iris"

	myhttps "app/backend/common/util/https"
	myregistry "app/backend/model/yce/registry"
	myerror "app/backend/common/yce/error"
)

type ListRegistryController struct {
	*iris.Context
	c *myhttps.HttpsClient
	BaseUrl string
	Registry *myregistry.Registry
}

func NewListRegistryController() *ListRegistryController {

	lrc := new(ListRegistryController)

	r := myregistry.NewRegistry(myregistry.REGISTRY_HOST, myregistry.REGISTRY_PORT, myregistry.REGISTRY_CERT)

	lrc.c = myhttps.NewHttpsClient(r.Host, r.Port, r.Cert)

	lrc.BaseUrl = "http://" + r.Host + ":" + r.Port

	lrc.Registry = r

	return lrc
}

// curl --cacert /etc/docker/certs.d/registry.test.com\:5000/domain.crt
// 		-X GET https://registry.test.com:5000/v2/_catalog | jq .
// {
//   "repositories": [
//   "busybox",
//   "capttofu/mysql_master_kubernetes",
// 	 "capttofu/mysql_slave_kubernetes",
//   "ceph/daemon",
//   "ceph/mds"
//   ]
// }
func (lrc *ListRegistryController) getRepositories() ([]string, error) {
	url := lrc.BaseUrl + "/v2/_catalog"

	client := lrc.c.Client
	resp, err := client.Get(url)

	if err != nil {
		log.Printf("ListRegistryController getRespositories Error: err=%s\n", err)
		return []string{}, nil
	}

	body, err := ioutil.ReadAll(resp.Body)

	repository := new(myregistry.Repository)

	err = json.Unmarshal(body, repository)

	if err != nil {
		log.Printf("ListRegistryController getRepositories Error: err=%s\n", err)
		return []string{}, nil
	}

	return repository.Repositories, nil

}

// curl --cacert /etc/docker/certs.d/registry.test.com\:5000/domain.crt -X
//	 GET https://registry.test.com:5000/v2/yeepay/nginx/tags/list
// {"name":"yeepay/nginx","tags":["latest"]}
func (lrc *ListRegistryController) getTagsList(name string) (*myregistry.Image, error) {

	// foreech repositories
	url := lrc.BaseUrl + name + "/tags/list"

	client := lrc.c.Client
	resp, err := client.Get(url)

	if err != nil {
		log.Printf("ListRegistryController getTagsList client.Get Error: err=%s\n", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Printf("ListRegistryController getTagsList ioutil.ReadAll Error: err=%s\n", err)
		return nil, err
	}

	image := new(myregistry.Image)

	err = json.Unmarshal(body, image)
	if err != nil {
		log.Printf("ListRegistryController getTagsList json.Unmarshal Error: err=%s\n", err)
		return nil, err
	}

	return image, nil
}

// GET /api/v1/registry/images
func (lrc ListRegistryController) Get() {

	log.Printf("ListRegistryController Get: baseUrl=%s\n", lrc.BaseUrl)

	var ye *myerror.YceError

	// Get repositories in the registry
	list, err := lrc.getRepositories()
	if err != nil {
		log.Printf("ListRegistryController getRepositories Error: err=%s\n", err)
		ye = myerror.NewYceError(1301, "ListRegistryController getRepositories Error!", "")
		js, _ := ye.EncodeJson()
		lrc.Write(js)
		return
	}

	if 0 == len(list) {
		log.Printf("ListRegistryController Repositories is empty!")
		ye = myerror.NewYceError(1302, "ListRegistryController Repositories is empty!", "")
		js, _ := ye.EncodeJson()
		lrc.Write(js)
		return
	}

	// Get detail info for ervery repository
	for _, repo := range list {
		image, err := lrc.getTagsList(repo)

		if err != nil {
			log.Printf("ListRegistryController getTagsList Error: err=%a\n", err)
			continue
		}

		lrc.Registry.Images = append(lrc.Registry.Images, *image)
	}

	if 0 == len(lrc.Registry.Images) {
		log.Printf("ListRegistryController Images is empty!")

		ye = myerror.NewYceError(1303, "ListRegistryController Images is empty!", "")
		js, _ := ye.EncodeJson()
		lrc.Write(js)
		return
	}

	images, _ := lrc.Registry.GetImagesList()

	ye = myerror.NewYceError(0, "OK", images)
	js, _ := ye.EncodeJson()
	lrc.Response.Header.Set("Access-Control-Allow-Origin", "*")

	log.Printf("%v\n", js)
	lrc.Write(js)

}

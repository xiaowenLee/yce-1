package registry

import (
	"encoding/json"
	"github.com/kataras/iris"
	"io/ioutil"

	myhttps "app/backend/common/util/https"
	myerror "app/backend/common/yce/error"
	myregistry "app/backend/model/yce/registry"
	mylog "app/backend/common/util/log"
)


type ListRegistryController struct {
	*iris.Context
	c        *myhttps.HttpsClient
	BaseUrl  string
	Registry *myregistry.Registry
}

// curl --cacert /etc/docker/certs.d/registry.test.com\:5000/domain.crt
// 		-X GET https://registry.test.com:5000/v2/_catalog | jq .
// {
//   "repositories": [
//   "busybox",
//   "ceph/mds"
//   ]
// }
func (lrc *ListRegistryController) getRepositories() ([]string, error) {
	url := lrc.BaseUrl + "/v2/_catalog"

	client := lrc.c.Client
	resp, err := client.Get(url)

	if err != nil {
		mylog.Log.Errorf("ListRegistryController getRespositories Error: err=%s", err)
		return []string{}, nil
	}

	body, err := ioutil.ReadAll(resp.Body)

	repository := new(myregistry.Repository)

	err = json.Unmarshal(body, repository)

	if err != nil {
		mylog.Log.Errorf("ListRegistryController getRepositories Error: err=%s", err)
		return []string{}, nil
	}

	mylog.Log.Infoln("ListRegistryController getRepositories over")
	return repository.Repositories, nil

}

// curl --cacert /etc/docker/certs.d/registry.test.com\:5000/domain.crt -X
//	 GET https://registry.test.com:5000/v2/yeepay/nginx/tags/list
// {"name":"yeepay/nginx","tags":["latest"]}
func (lrc *ListRegistryController) getTagsList(name string) (*myregistry.Image, error) {

	// foreech repositories
	url := lrc.BaseUrl + "/v2/" + name + "/tags/list"

	client := lrc.c.Client
	resp, err := client.Get(url)

	if err != nil {
		mylog.Log.Errorf("ListRegistryController getTagsList client.Get Error: err=%s", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		mylog.Log.Errorf("ListRegistryController getTagsList ioutil.ReadAll Error: err=%s", err)
		return nil, err
	}

	image := new(myregistry.Image)

	err = json.Unmarshal(body, image)
	if err != nil {
		mylog.Log.Errorf("ListRegistryController getTagsList json.Unmarshal Error: err=%s", err)
		return nil, err
	}

	return image, nil
}

// GET /api/v1/registry/images
func (lrc ListRegistryController) Get() {
	//TODO: Validate Session ?
	// init
	r := myregistry.NewRegistry(myregistry.REGISTRY_HOST, myregistry.REGISTRY_PORT, myregistry.REGISTRY_CERT)
	lrc.c = myhttps.NewHttpsClient(r.Host, r.Port, r.Cert)
	lrc.BaseUrl = "https://" + r.Host + ":" + r.Port
	lrc.Registry = r

	var ye *myerror.YceError

	// Get repositories in the registry
	list, err := lrc.getRepositories()
	if err != nil {
		mylog.Log.Errorf("ListRegistryController getRepositories Error: err=%s", err)
		// ye = myerror.NewYceError(1301, "ListRegistryController getRepositories Error!", "")
		ye = myerror.NewYceError(1301, "")
		js, _ := ye.EncodeJson()
		lrc.Write(js)
		return
	}

	if 0 == len(list) {
		mylog.Log.Errorf("ListRegistryController Repositories is empty!")
		// ye = myerror.NewYceError(1302, "ListRegistryController Repositories is empty!", "")
		ye = myerror.NewYceError(1302, "")
		js, _ := ye.EncodeJson()
		lrc.Write(js)
		return
	}

	// Get detail info for ervery repository
	for _, repo := range list {
		image, err := lrc.getTagsList(repo)

		if err != nil {
			mylog.Log.Errorf("ListRegistryController getTagsList Error: err=%a", err)
			continue
		}

		lrc.Registry.Images = append(lrc.Registry.Images, *image)
	}

	if 0 == len(lrc.Registry.Images) {
		mylog.Log.Errorf("ListRegistryController Images is empty!")

		// ye = myerror.NewYceError(1303, "ListRegistryController Images is empty!", "")
		ye = myerror.NewYceError(1303, "")
		js, _ := ye.EncodeJson()
		lrc.Write(js)
		return
	}

	images, _ := lrc.Registry.GetImagesList()

	ye = myerror.NewYceError(0, images)
	js, _ := ye.EncodeJson()
	lrc.Response.Header.Set("Access-Control-Allow-Origin", "*")

	lrc.Write(js)

	mylog.Log.Infoln("ListRegistryController Get over!")
}

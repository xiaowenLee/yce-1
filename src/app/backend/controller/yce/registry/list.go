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
	Ye *myerror.YceError
}

func (lrc *ListRegistryController) WriteBack() {
	lrc.Response.Header.Set("Access-Control-Allow-Origin", "*")
	mylog.Log.Infof("ListRegistryController Response YceError: controller=%p, code=%d, note=%s", lrc, lrc.Ye.Code, myerror.Errors[lrc.Ye.Code].LogMsg)
	lrc.Write(lrc.Ye.String())
}

// curl --cacert /etc/docker/certs.d/registry.test.com\:5000/domain.crt
// 		-X GET https://registry.test.com:5000/v2/_catalog | jq .
// {
//   "repositories": [
//   "busybox",
//   "ceph/mds"
//   ]
// }
func (lrc *ListRegistryController) getRepositories() ([]string) {
	url := lrc.BaseUrl + "/v2/_catalog"

	client := lrc.c.Client
	resp, err := client.Get(url)

	if err != nil {
		mylog.Log.Errorf("ListRegistryController getRespositories Error: err=%s", err)
		lrc.Ye = myerror.NewYceError(myerror.EREGISTRY_GET, "")
		return []string{}
	}

	body, err := ioutil.ReadAll(resp.Body)

	repository := new(myregistry.Repository)

	err = json.Unmarshal(body, repository)
	if err != nil {
		mylog.Log.Errorf("ListRegistryController getRepositories Error: err=%s", err)
		lrc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return []string{}
	}

	mylog.Log.Infoln("ListRegistryController getRepositories over")
	return repository.Repositories

}

// curl --cacert /etc/docker/certs.d/registry.test.com\:5000/domain.crt -X
//	 GET https://registry.test.com:5000/v2/yeepay/nginx/tags/list
// {"name":"yeepay/nginx","tags":["latest"]}
func (lrc *ListRegistryController) getTagsList(name string) (*myregistry.Image) {

	// foreech repositories
	url := lrc.BaseUrl + "/v2/" + name + "/tags/list"

	client := lrc.c.Client
	resp, err := client.Get(url)

	if err != nil {
		mylog.Log.Errorf("ListRegistryController getTagsList client.Get Error: err=%s", err)
		lrc.Ye = myerror.NewYceError(myerror.EREGISTRY_GET, "")
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		mylog.Log.Errorf("ListRegistryController getTagsList ioutil.ReadAll Error: err=%s", err)
		lrc.Ye = myerror.NewYceError(myerror.EJSON, "")
		return nil
	}

	image := new(myregistry.Image)

	err = json.Unmarshal(body, image)
	if err != nil {
		mylog.Log.Errorf("ListRegistryController getTagsList json.Unmarshal Error: err=%s", err)
		return nil
	}

	// Add Prefix
	image.Name = myregistry.REGISTRY_HOST + ":" + myregistry.REGISTRY_PORT + "/" + image.Name

	return image
}

// GET /api/v1/registry/images
func (lrc ListRegistryController) Get() {
	//TODO: Validate Session ?
	// init
	r := myregistry.NewRegistry(myregistry.REGISTRY_HOST, myregistry.REGISTRY_PORT, myregistry.REGISTRY_CERT)
	lrc.c = myhttps.NewHttpsClient(r.Host, r.Port, r.Cert)
	lrc.BaseUrl = "https://" + r.Host + ":" + r.Port
	lrc.Registry = r

	// Get repositories in the registry
	list := lrc.getRepositories()
	if lrc.Ye != nil {
		lrc.WriteBack()
		return
	}

	if 0 == len(list) {
		lrc.Ye = myerror.NewYceError(myerror.EREGISTRY, "")
		lrc.WriteBack()
		return
	}

	// Get detail info for ervery repository
	for _, repo := range list {
		image := lrc.getTagsList(repo)
		lrc.Registry.Images = append(lrc.Registry.Images, *image)
	}

	if 0 == len(lrc.Registry.Images) {
		lrc.Ye = myerror.NewYceError(myerror.EREGISTRY, "")
		lrc.WriteBack()
		return
	}

	images, _ := lrc.Registry.GetImagesList()

	lrc.Ye = myerror.NewYceError(myerror.EOK, images)
	lrc.WriteBack()

	mylog.Log.Infoln("ListRegistryController Get over!")
}

package registry

import (
	"log"

	"github.com/kataras/iris"
	myregistry "app/backend/model/yce/registry"
)

type ListRegistryController struct {
	*iris.Context
}

// curl --cacert /etc/docker/certs.d/registry.test.com\:5000/domain.crt -X GET https://registry.test.com:5000/v2/_catalog | jq .
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
	var repositories []string
	return repositories, nil

}

// curl --cacert /etc/docker/certs.d/registry.test.com\:5000/domain.crt -X GET https://registry.test.com:5000/v2/yeepay/nginx/tags/list
// {"name":"yeepay/nginx","tags":["latest"]}
func (lrc *ListRegistryController) getTagsList(name string) {

	// repositories := lrc.getRepositories()
	// foreech repositories
}

// GET /api/v1/registry/images
func (lrc ListRegistryController) Get() {

}
package deployment

import (
	hc "app/backend/common/util/http"
	mc "app/backend/common/util/mysql"
	deploy "app/backend/model/k8s/deployment"
	"encoding/json"
	"error"
	"log"
	"net/http"
)

type DeploymentController struct {
}

func (dc *DeploymentController) Get() {

}
func (dc *DeploymentController) GetById() {

}
func (dc *DeploymentController) Post() {

	// Json
	// foreach $idc in $idc-array
	// Post the json to k8s-master
	// if success, insert db
	// end
}
func (dc *DeploymentController) DeleteById() {

}

func (dc *DeploymentController) EncodeJson(v interface{}) ([]byte, error) {
	prefix := ""
	indent := "  "
	return json.MarshalIndent(v, prefix, indent)
}

func (dc *DeploymentController) DecodeJson(data []byte, v interface{}) (v interface{}, error) {

}

package registry

import (
	"encoding/json"
	mylog "app/backend/common/util/log"
)

var log =  mylog.Log

const (
	REGISTRY_HOST = "img.reg.3g"
	REGISTRY_PORT = "15000"
	REGISTRY_CERT = "domain.crt"
)

type Image struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type Repository struct {
	Repositories []string `json:"repositories"`
}

type Registry struct {
	Host   string  `json:"host"`
	Port   string  `json:"port"`
	Cert   string  `json:"cert"`
	Images []Image `json:"images"`
}

func NewRegistry(host, port, cert string) *Registry {

	return &Registry{
		Host: host,
		Port: port,
		Cert: cert,
	}
}

func (r *Registry) GetImagesList() (string, error) {
	images, err := json.Marshal(r.Images)

	if err != nil {
		log.Errorf("GetImageList Error: err=%s\n", err)
		return "", err
	}

	return string(images), nil
}

func (r *Registry) DecodeJson(data string) error {
	err := json.Unmarshal([]byte(data), r)

	if err != nil {
		log.Errorf("DecodeJson Error: err=%s\n", err)
		return err
	}

	return nil
}

func (r *Registry) EncodeJson() (string, error) {
	data, err := json.Marshal(r)
	if err != nil {
		log.Errorf("EncodeJson Error: err=%s\n", err)
		return "", err
	}
	return string(data), nil
}

package registry

import (
	"encoding/json"
	"log"
)

const (

	REGISTRY_HOST = "registry.docker"
	REGISTRY_PORT = "5000"
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
	Host   string `json:"host"`
	Port   string `json:"port"`
	Cert   string `json:"cert"`
	Images []Image `json:"images"`
}

func NewRegistry(host, port, cert string) *Registry {

	return &Registry {
		Host: host,
		Port: port,
		Cert: cert,
	}
}

func (r *Registry) GetImagesList() (string, error) {
	images, err := json.Marshal(r.Images)

	if err != nil {
		log.Printf("GetImageList Error: err=%s\n", err)
		return "", err
	}

	return string(images), nil
}

func (r *Registry) DecodeJson(data string) error {
	err := json.Unmarshal([]byte(data), r)

	if err != nil {
		log.Printf("DecodeJson Error: err=%s\n", err)
		return err
	}

	return nil
}

func (r *Registry) EncodeJson() (string, error) {
	data, err := json.Marshal(r)
	if err != nil {
		log.Printf("EncodeJson Error: err=%s\n", err)
		return "", err
	}
	return string(data), nil
}

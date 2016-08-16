package registry

import (
	"encoding/json"
	"log"
)

type Image struct {
	Name string   `json:"name"`
	Tags []string `json:"list"`
}

type Registry struct {
	Host   string  `json:"host"`
	Port   int32   `json:"port"`
	cert   string  `json:"cert"`
	Images []Image `json:"images"`
}

func NewRegistry(host, cert string, port int32) *Registry {

	return &Registry{
		Host: host,
		cert: cert,
		Port: port,
	}
}

func (r *Registry) GetImageList() (string, error) {
	images, err := json.MarshalIndent(r.Images, "", " ")
	if err != nil {
		log.Printf("GetImageList marshal error: err=%s\n", err)
		return "", err
	}
	return images, nil
}

func (r *Registry) DecodeJson(data string) error {
	err := json.Unmarshal([]byte(data), r)

	if err != nil {
		log.Printf("DecodeJson Error: err=%s\n", err)
		return err
	}

	return nil
}

func (r *Registry) EncodeJson() (string , error) {
	data, err := json.MarshalIndent(r, "", " ")
	if err != nil {
		log.Printf("EncodeJson Error: err=%s\n", err)
		return "", err
	}
	return string(data), nil
}

package registry

import (
	"encoding/json"
	"log"
)

const CERT = `
-----BEGIN CERTIFICATE-----
MIIFlTCCA32gAwIBAgIJAL/V18aLUarxMA0GCSqGSIb3DQEBCwUAMGExCzAJBgNV
BAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQxGjAYBgNVBAMMEXJlZ2lzdHJ5LnRlc3QuY29tMB4XDTE1
MTIzMTA2MjYwOFoXDTE2MTIzMDA2MjYwOFowYTELMAkGA1UEBhMCQVUxEzARBgNV
BAgMClNvbWUtU3RhdGUxITAfBgNVBAoMGEludGVybmV0IFdpZGdpdHMgUHR5IEx0
ZDEaMBgGA1UEAwwRcmVnaXN0cnkudGVzdC5jb20wggIiMA0GCSqGSIb3DQEBAQUA
A4ICDwAwggIKAoICAQDYg769YT8Y4jDpWaTDALj3Rq3QqSn0fMYyuFGxWzMZIgz4
svjBVcghPl1UMgPVNRQppQmsuMlgmSxBSX2QUCjkyamYdigsD2TYnR/Bue7x13H2
cR5MytBuQUUaNkC38KAMgjIwq5ztqN+8ZK3yRtZRB4VBj5v6s7M+6ZhvNyp8OM4U
h6R2oatmva3PBHqGCVeTQ00oEshOHSXkrchwpRaWqLaJU1pbS10GhmCQAUaLBeZ1
2M5MVHPFMWDfTXcE5gNdMnJxzZ3oQGYJmr7lmHCsfo7WmzCAn1igSH8mmZR3LlCP
tUfbQDT2jlBvFq2pI9J1tdAZXdtizrpFVRPUZesrsBIYLrQ1rCx8jrK/4l9uxmUj
7TcgtGQ+Ziu0Pg2bZJL6Lz7Vl99eZGS3s8WSO/x+sPZlDjPFP32+PE+8LYekG8ie
7bBJoOa9UzHMC8pxgnuBSNAmRXHwwPhENA3LZFt3FunkfGNNR8HBxvqe9pNOYhBq
k3gCoHd94GmgmeKlf8h06uhWrXrzV32CkUh2OGp0UqRrfyazSeCzAvPrbUSx8hv5
PHMCZ1MppU5E59e8cNJaR0PTxqZeR7r3fS8bRQBONe/jontc4y201gxcRmNNKTVa
2vMn+GCxoiAjhHU4xnu2+n51EIm0WPumLFmIejfb4DfeI7I5QjslQVmxB1+HZQID
AQABo1AwTjAdBgNVHQ4EFgQU53+nvnQLQkZc6RBRqILqiAD0AVMwHwYDVR0jBBgw
FoAU53+nvnQLQkZc6RBRqILqiAD0AVMwDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0B
AQsFAAOCAgEAoKX4UrjUKidEp6d19Hl69PyUHFNOuH7FBL+UFJKyWRw4d1b8FLaV
fMZwrHGKGLzCbsfUcLJF6Yw7SV55ZZLYHSpXy/zOGWKevqdWpg8jIMek7t90u7W5
sEN8Pc4iJj0gu4rQhHOht+Ku/9guvC6gjCEIyqAM/8nKbJhW4w/Yl3iDfE70XFKA
9dLyTqSMZ7ntEM22ge9zgLg+y+O4QQfih009+zFPhMpx4JtdgI2rfeQHQbsvBi9P
bGdfTEE2dFplqdHq3y2aLx3C3WeHZjpBdblk1lSnG3nf+H98vRebRTEq29aMPfOP
mHC5wD5TLPSAlrsxVcMgr3OEvb0KcCrQNBapWqwASkiOT6ouZENRAimPg3YPR/TK
Kvk0Bexzg5pOKkNI79IK5sl7+7pkwTAbJvhBgyScZOuSFRbpEOCcTNql18uf8DLl
BL9K5osTbuDwO4YljRA9Ig9OtgYs9p0+XkWhkcOjt7LuzPXSqbi654UG9RyTVNDF
4BaOsRBgXpSoNp1a1n1YrvA9gJaVWBdK8EyQSRXQrT+RW3lEwkS2z+iK14NY2Mfr
TaUKMkzJ7hRhK2CdnJd8A0ud1j65pMtLs6jrMxNAXT7B0v46aAxYN1LDU+CMrP8J
T7rkh2yVNeXLki1gStjJeso+Qo89Yh+7dPUDsiQvh8sEfPGX6a3aA8Y=
-----END CERTIFICATE-----
`

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

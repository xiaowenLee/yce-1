package https

import (
	"crypto/tls"
	"crypto/x509"
	// "log"
	"net/http"
)

type HttpsClient struct {
	Host   string `json:"host"`
	Port   int32  `json:"port"`
	pool   x509.CertPool
	client *http.Client
}

func NewHttpsClient(host string, port int32, cert string) *HttpsClient {

	https := &HttpsClient{
		Host: host,
		Port: port,
	}

	https.pool.AppendCertsFromPEM([]byte(cert))


	tr := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: https.pool},
		DisableCompression: true,
	}

	https.client = &http.Client{Transport: tr}

	return https
}

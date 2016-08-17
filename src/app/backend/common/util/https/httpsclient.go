package https

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"fmt"
	"io/ioutil"
)

type HttpsClient struct {
	Host   string `json:"host"`
	Port   string  `json:"port"`
	pool   *x509.CertPool
	Client *http.Client
}

func NewHttpsClient(host, port, cert string) *HttpsClient {

	https := &HttpsClient{
		Host: host,
		Port: port,
		pool: x509.NewCertPool(),
	}

	caCrt, err := ioutil.ReadFile(cert)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return nil
	}

	https.pool.AppendCertsFromPEM(caCrt)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: https.pool},
		DisableCompression: true,
	}

	https.Client = &http.Client{Transport: tr}

	return https
}

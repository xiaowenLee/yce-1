package https

import (
	"crypto/tls"
	"crypto/x509"
	// "log"
	"net/http"
	"fmt"
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

	fmt.Printf("%v\n", https.pool)
	// https.pool = new(x509.CertPool)
	https.pool.AppendCertsFromPEM([]byte(cert))


	tr := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: https.pool},
		DisableCompression: true,
	}

	https.Client = &http.Client{Transport: tr}

	return https
}

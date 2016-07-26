package main 

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"encoding/json"
	"crypto/tls"
	"io/ioutil"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, This is k8sapi!", html.EscapeString(r.URL.Path))
}

func Get(url string) (body []byte, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},	
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		log.Println(err)
		panic(err)
		return nil, err	
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err	
	}

	return body, nil
}

func resolveToEndpointStruct(response []byte) (s *EndpointType, err error) {
	err = json.Unmarshal(response, &s)
	return s, err 
}

func resolveToServiceStruct(response []byte) (s *ServiceType, err error) {
        err = json.Unmarshal(response, &s)
        return s, err
}

func Endpointlist(w http.ResponseWriter, r *http.Request) {
	var response []byte
        var err error

        response, err = Get("http://172.21.1.11:8080/api/v1/endpoints")
        if err != nil {
                panic(err)
                log.Println(err)
        }

        rs, err := resolveToEndpointStruct(response)
        if err != nil {
                panic(err)
                log.Println(err)
        }

	json.NewEncoder(w).Encode(rs)
}

func Servicelist(w http.ResponseWriter, r *http.Request) {
        var response []byte
        var err error

        response, err = Get("http://172.21.1.11:8080/api/v1/services")
        if err != nil {
                panic(err)
                log.Println(err)
        }

        rs, err := resolveToServiceStruct(response)
        if err != nil {
                panic(err)
                log.Println(err)
        }

        json.NewEncoder(w).Encode(rs)
}

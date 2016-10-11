package main

import (
	"fmt"
	"io"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	//"bytes"
	"io/ioutil"
)

const (
	SERVER = "master:8080"
)

type logType struct {
	Out io.Writer
}

func (l logType) Get() {
	//ns := api.NamespaceDefault
	ns := "ops"
	podName := "nginx-test-818178620-ucju9"

	opts := &api.PodLogOptions{
		//Follow:     true,
		Follow:     false,
		Timestamps: true,
	}

	config := &restclient.Config{
		Host: SERVER,
	}

	c, err := client.New(config)
	if err != nil {
		fmt.Printf("Could not connect to k8s api: err=%s\n", err)
	}

	reader, err := c.Pods(ns).GetLogs(podName, opts).Stream()
	if err != nil {
		fmt.Println(err)
	}
	defer reader.Close()

	if b, err := ioutil.ReadAll(reader); err == nil {
		fmt.Println(string(b))
	} else {
		fmt.Println(err)
	}

	/*
		buf := new(bytes.Buffer)
		len, err := buf.ReadFrom(reader)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(len)
		s := buf.String()
		fmt.Println(s)
	*/

	/*
		body, err := c.Pods(ns).GetLogs(podName, opts).Stream()
		if err != nil {
			fmt.Println(err)
		}


		l.Out = new(io.Writer)
		len, err := io.Copy(l.Out, body)
		if err != nil {
			fmt.Println(err)
		}




		b := make([]byte, 100)
		l.Out.Write(b)

		fmt.Println(len)
		fmt.Println(string(b))
	*/
}

func main() {
	l := new(logType)
	l.Get()
}

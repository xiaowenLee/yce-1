package main

import "fmt"
import "encoding/json"

type SpecType struct {
	Name string `json:"name"`
	Age int32 `json:"age"`
}

type A struct {
	Spec SpecType `json:"spec"`
}

type B struct {
	A
	Kind string `json:"kind"`
	Version string `json:"version"`
}

func main() {

	b := B{
		Kind: "Node",
		Version: "v1beta3",
		A: A{
			Spec: SpecType{
				Name: "litanhua",
				Age: 32,
			},
		},
	}

	data, _ := json.MarshalIndent(b, "", " ")
	fmt.Println(string(data))
}

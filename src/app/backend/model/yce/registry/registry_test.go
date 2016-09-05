package registry

import (
	"fmt"
	"testing"
)

func TestRegistry_DecodeJson(t *testing.T) {
	r := NewRegistry(REGISTRY_HOST, REGISTRY_PORT, REGISTRY_CERT)
	fmt.Printf("%v\n", r)

	images, err := r.GetImagesList()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Images: %s\n", images)

	js, err := r.EncodeJson()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("EncodeJson: %s\n", js)

	re := new(Registry)
	err = re.DecodeJson(js)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%v\n", re)

}

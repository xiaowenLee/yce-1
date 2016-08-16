package registry

import (
	"testing"
	"fmt"
)

func TestRegistry_DecodeJson(t *testing.T) {
	r := NewRegistry("registry.test.com", CERT, 5000)
	fmt.Printf("%v\n", r)

	images, err := r.GetImageList()
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

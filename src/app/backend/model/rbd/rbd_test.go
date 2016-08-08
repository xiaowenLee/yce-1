package rbd

import (
	"fmt"
	"testing"
)

func Test_NewRbdBlock_Decode_Encode(*testing.T) {
	fmt.Println("Test_NewRbdBlock_Decode_Encode")
	rbd := NewRbdBlock("arbd", "rbd", "ext4", 1000)
	fmt.Printf("%s\n", rbd.EncodeJson())

	r := new(RbdBlock)
	r.DecodeJson(rbd.EncodeJson())

	fmt.Printf("%s\n", r.Image)
}

func Test_CreateRbdBlock(t *testing.T) {
	fmt.Println("Test_CreateRbdBlock")

	rbd := NewRbdBlock("arbd", "rbd", "ext4", 1000)

	fmt.Println("RBD Create")
	// Create a rbd block: name=arbd, size=1000MB, fs=ext4, pool=rbd
	err := rbd.Create()
	if err != nil {
		t.Error(err)
	}

	fmt.Println("RBD Map")
	// Map rbd device
	err = rbd.Map()
	if err != nil {
		t.Error(err)
	}

	fmt.Println("RBD ShowMapped")
	// Show mapped point
	err = rbd.ShowMapped()
	if err != nil {
		t.Error(err)
	}

	fmt.Println("RBD UnMap")
	// UnMap the device
	err = rbd.UnMap()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf(rbd.EncodeJson())

}

func Test_RemoveRbdBlock(t *testing.T) {
	fmt.Println("Test_RemoveRbdBlock")
	rbd := NewRbdBlock("arbd", "rbd", "ext4", 1000)
	err := rbd.Remove()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf(rbd.EncodeJson())
}

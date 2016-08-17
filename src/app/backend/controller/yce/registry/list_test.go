package registry

import (
	"fmt"
	"testing"
)

func Test_ListRegistryController_Get(*testing.T) {

	lrc := NewListRegistryController()

	fmt.Printf("%v\n", lrc)
	// fmt.Printf("%s\n", lrc.Get())
	lrc.Get()
}
package registry

import (
	"github.com/kataras/iris"
	myregistry "app/backend/model/yce/registry"
)

type ListRegistryController struct {
	*iris.Context
}
package navlist

import (
	yce "app/backend/controller/yce"
)

type NavListController struct {
	yce.Controller
}

func (nlc NavListController) Get() {
	nlc.WriteOk(navList)
}

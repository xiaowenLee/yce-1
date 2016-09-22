package apis

import (
	yce "app/backend/controller/yce"
)

type ApisController struct {
	yce.Controller
}

// GET /path
func (ac ApisController) Get() {
	ac.Write(APIS)
}

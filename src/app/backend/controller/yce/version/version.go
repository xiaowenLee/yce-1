package version

import (
	yce "app/backend/controller/yce"
)

type VersionController struct {
	yce.Controller
}

// GET /version
func (vc VersionController) Get() {
	vc.Write(VERSION)
}


package healthz

import (
	yce "app/backend/controller/yce"
)

type HealthzController struct {
	yce.Controller
}

// GET /healthz
func (hc HealthzController) Get() {
	hc.Write("OK")
}


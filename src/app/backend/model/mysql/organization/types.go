package organization

type QuotaType struct {
	CpuQuota int32 `json:"cpuQuota"`
	MemQuota int32 `json:"memQuota"`
}


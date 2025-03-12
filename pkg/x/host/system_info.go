package host

// DefaultSystemInfo is the default implementation of SystemInfo
type DefaultSystemInfo struct{}

func (d DefaultSystemInfo) IsOS(os string) error {
	return IsOS(os)
}

func (d DefaultSystemInfo) IsArchIn(archs []string) error {
	return IsArchIn(archs)
}

func (d DefaultSystemInfo) IsRAMEnough(minRAM string) error {
	return IsRAMEnough(minRAM)
}

func (d DefaultSystemInfo) IsCPUEnough(minCPU int) error {
	return IsCPUEnough(minCPU)
}

func (d DefaultSystemInfo) IsDiskSpaceEnough(path, size string) error {
	return IsDiskSpaceEnough(path, size)
}

func (d DefaultSystemInfo) IsCurlOK(urls []string) error {
	return IsCurlOK(urls)
}

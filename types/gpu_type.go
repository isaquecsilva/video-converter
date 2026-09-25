package types

type GpuType = string

const (
	GpuTypeNvidia GpuType = "nvidia"
	GpuTypeAMD    GpuType = "amd"
	GpuTypeNone   GpuType = "none"
)

package codecs

import (
	"fmt"
	"video-converter/types"
)

var ErrNoCodecForGpu error = fmt.Errorf("Nenhum codec disponível para gpu")

type VideoCodec string

func (vc VideoCodec) String() string {
	return string(vc)
}

type VideoCodecKit interface {
	Name() string
	String(gpuType types.GpuType) (VideoCodec, error)
	HasGpuSupport(gpuType types.GpuType) bool
}

package codecs

import (
	"fmt"
	"video-converter/types"
)

type baseCodecKit struct {
	name   string
	codecs map[types.GpuType]VideoCodec
}

func newBaseCodecKit(name string, codecs map[types.GpuType]VideoCodec) baseCodecKit {
	return baseCodecKit{
		name:   name,
		codecs: codecs,
	}
}

func (bc baseCodecKit) Name() string {
	return bc.name
}

func (bc baseCodecKit) String(gpuType types.GpuType) (VideoCodec, error) {
	codec, ok := bc.codecs[gpuType]
	if !ok {
		return "", fmt.Errorf("%w: %v", ErrNoCodecForGpu, gpuType)
	}

	return codec, nil
}

func (bc baseCodecKit) HasGpuSupport(gpuType types.GpuType) bool {
	_, ok := bc.codecs[gpuType]
	return ok
}

package codecs

import "video-converter/types"

type videoCodecKitAV1 struct {
	baseCodecKit
}

func NewVideoCodecKitAV1() videoCodecKitAV1 {
	kit := videoCodecKitAV1{}
	kit.name = "av1"
	kit.codecs = map[types.GpuType]VideoCodec{
		types.GpuTypeNone:   VideoCodecAV1,
		types.GpuTypeAMD:    VideoCodecAV1Amf,
		types.GpuTypeNvidia: VideoCodecAV1NvEnc,
	}

	return kit
}

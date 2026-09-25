package codecs

import "video-converter/types"

type videoCodecKitVP8 struct {
	baseCodecKit
}

func NewVideoCodecKitVP8() videoCodecKitVP8 {
	kit := videoCodecKitVP8{}
	kit.name = "vp8"
	kit.codecs = map[types.GpuType]VideoCodec{
		types.GpuTypeNone: VideoCodecVP8,
	}

	return kit
}

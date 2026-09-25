package codecs

import "video-converter/types"

type videoCodecKitVP9 struct {
	baseCodecKit
}

func NewVideoCodecKitVP9() videoCodecKitVP9 {
	kit := videoCodecKitVP9{}
	kit.name = "vp9"
	kit.codecs = map[types.GpuType]VideoCodec{
		types.GpuTypeNone: VideoCodecVP9,
	}

	return kit
}

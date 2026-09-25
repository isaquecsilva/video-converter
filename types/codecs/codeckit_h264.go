package codecs

import "video-converter/types"

type videoCodecKitH264 struct {
	baseCodecKit
}

func NewVideoCodecKitH264() videoCodecKitH264 {
	kit := videoCodecKitH264{}
	kit.name = "h264"
	kit.codecs = map[types.GpuType]VideoCodec{
		types.GpuTypeNone:   VideoCodecH264,
		types.GpuTypeAMD:    VideoCodecH264Amf,
		types.GpuTypeNvidia: VideoCodecH264NvEnc,
	}

	return kit
}

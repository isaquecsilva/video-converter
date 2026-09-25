package videoformat

import "video-converter/types/codecs"

type videoFormatWebm struct {
	baseVideoFormat
}

func NewVideoFormatWebm() videoFormatWebm {
	return videoFormatWebm{
		baseVideoFormat: baseVideoFormat{
			name: "webm",
			ext:  ".webm",
			codecKitList: []codecs.VideoCodecKit{
				codecs.NewVideoCodecKitVP9(),
				codecs.NewVideoCodecKitVP8(),
				codecs.NewVideoCodecKitAV1(),
			},
		},
	}
}

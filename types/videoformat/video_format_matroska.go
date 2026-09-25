package videoformat

import "video-converter/types/codecs"

type videoFormatMatroska struct {
	baseVideoFormat
}

func NewVideoFormatMatroska() videoFormatMatroska {
	return videoFormatMatroska{
		baseVideoFormat: baseVideoFormat{
			name: "matroska",
			ext:  ".mkv",
			codecKitList: []codecs.VideoCodecKit{
				codecs.NewVideoCodecKitH264(),
				codecs.NewVideoCodecKitAV1(),
				codecs.NewVideoCodecKitVP8(),
				codecs.NewVideoCodecKitVP9(),
			},
		},
	}
}

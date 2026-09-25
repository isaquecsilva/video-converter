package videoformat

import (
	"video-converter/types/codecs"
)

type videoFormatMp4 struct {
	baseVideoFormat
}

func NewVideoFormatMp4() videoFormatMp4 {
	return videoFormatMp4{
		baseVideoFormat: baseVideoFormat{
			name: "mp4",
			ext:  ".mp4",
			codecKitList: []codecs.VideoCodecKit{
				codecs.NewVideoCodecKitH264(),
			},
		},
	}
}

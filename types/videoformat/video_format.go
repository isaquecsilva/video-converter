package videoformat

import (
	"video-converter/types/codecs"
)

type Format = string

type VideoFormat interface {
	Name() string
	Ext() string
	SupportedCodecs() []codecs.VideoCodecKit
	GetCodecKit(codecName string) (codecs.VideoCodecKit, error)
	DefaultCodecKit() codecs.VideoCodecKit
}

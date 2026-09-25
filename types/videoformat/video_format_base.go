package videoformat

import (
	"fmt"
	"video-converter/types/codecs"
)

type baseVideoFormat struct {
	name         string
	ext          string
	codecKitList []codecs.VideoCodecKit
}

func (b baseVideoFormat) Name() string {
	return b.name
}

func (b baseVideoFormat) Ext() string {
	return b.ext
}

func (b baseVideoFormat) SupportedCodecs() []codecs.VideoCodecKit {
	return b.codecKitList
}

func (b baseVideoFormat) GetCodecKit(codecName string) (codecs.VideoCodecKit, error) {
	for _, kit := range b.codecKitList {
		if kit.Name() == codecName {
			return kit, nil
		}
	}

	return nil, fmt.Errorf("%w: %v", ErrVideoCodecKitNotFound, codecName)
}

func (b baseVideoFormat) DefaultCodecKit() codecs.VideoCodecKit {
	if len(b.codecKitList) > 0 {
		return b.codecKitList[0]
	}

	return nil
}

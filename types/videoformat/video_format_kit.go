package videoformat

import "fmt"

var (
	ErrVideoFormatNotFound error = fmt.Errorf("Formato de vídeo não encontrado")
)

type VideoFormatKit interface {
	SupportedFormats() []VideoFormat
	SelectFormat(formatName string) (VideoFormat, error)
	DefaultFormat() VideoFormat
}

var _ VideoFormatKit = videoFormatKitImpl{}

type videoFormatKitImpl struct {
	formats []VideoFormat
}

func NewVideoFormatKitImpl() videoFormatKitImpl {
	return videoFormatKitImpl{
		formats: []VideoFormat{
			NewVideoFormatMp4(),
			NewVideoFormatMatroska(),
			NewVideoFormatWebm(),
		},
	}
}

func (kit videoFormatKitImpl) SupportedFormats() []VideoFormat {
	return kit.formats
}

func (kit videoFormatKitImpl) SelectFormat(formatName string) (VideoFormat, error) {
	for _, f := range kit.formats {
		if f.Name() == formatName {
			return f, nil
		}
	}

	return nil, fmt.Errorf("%w: %v", ErrVideoFormatNotFound, formatName)
}

func (kit videoFormatKitImpl) DefaultFormat() VideoFormat {
	return kit.formats[0]
}

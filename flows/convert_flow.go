package flows

import (
	"bufio"
	"video-converter/services/converter"
	"video-converter/services/discover"
	"video-converter/services/metadata"
	"video-converter/types/videoformat"
	"video-converter/utils"
	"video-converter/utils/printer"
)

type ConvertionFlow struct {
	reader    *bufio.Reader
	metadata  *metadata.VideoMetadataExtractor
	formatKit videoformat.VideoFormatKit
	discover  *discover.GpuDiscover
	conv      *converter.VideoConverter
	printer   printer.Printer
	baseFlow  BaseFlow

	params converter.ConvertParams
}

func NewConvertionFlow(
	reader *bufio.Reader,
	metadata *metadata.VideoMetadataExtractor,
	baseFlow BaseFlow,
	conv *converter.VideoConverter,
) Flow {
	return &ConvertionFlow{
		metadata: metadata,
		conv:     conv,
		reader:   reader,
		baseFlow: baseFlow,
		params:   converter.ConvertParams{},
	}
}

func (cf *ConvertionFlow) Run() (string, error) {
	cf.params.InputFile = cf.baseFlow.getInputFile()

	metadata, err := cf.metadata.ExtractMetadata(cf.params.InputFile)
	if err != nil {
		return "", err
	}

	cf.params.Format = cf.baseFlow.getTargetFormat()

	codec, err := cf.baseFlow.getCodec(cf.params.Format)
	if err != nil {
		return "", err
	}

	cf.params.Codec = codec
	cf.params.VideoBitrate = cf.baseFlow.getBitrate(metadata.VideoBitrateInKbps)

	ctx, stop := utils.SignalContext()
	defer stop()
	return cf.conv.Convert(ctx, cf.params)
}

package flows

import (
	"bufio"
	"fmt"
	"video-converter/services/metadata"
	"video-converter/services/vcut"
	"video-converter/utils"
	"video-converter/utils/printer"
)

type CutFlow struct {
	metadata *metadata.VideoMetadataExtractor
	cutter   vcut.VideoCutterInterface
	reader   *bufio.Reader
	printer  printer.Printer
	baseFlow BaseFlow
	params   vcut.CutParams
}

func NewCutFlow(
	reader *bufio.Reader,
	baseFlow BaseFlow,
	m *metadata.VideoMetadataExtractor,
	cutter vcut.VideoCutterInterface,
	printer printer.Printer,
) Flow {

	return &CutFlow{
		metadata: m,
		cutter:   cutter,
		printer:  printer,
		reader:   reader,
		baseFlow: baseFlow,
		params:   vcut.CutParams{},
	}
}

func (cf *CutFlow) Run() (string, error) {
	cf.params.InputFile = cf.baseFlow.getInputFile()

	metadata, err := cf.metadata.ExtractMetadata(cf.params.InputFile)
	if err != nil {
		return "", err
	}

	cf.params.TotalSeconds = int64(metadata.Duration.Seconds())
	cf.params.Format = cf.baseFlow.getTargetFormat()

	codec, err := cf.baseFlow.getCodec(cf.params.Format)
	if err != nil {
		return "", err
	}

	cf.params.Codec = codec
	cf.params.Bitrate = cf.baseFlow.getBitrate(metadata.VideoBitrateInKbps)

	cf.getTimestamps()

	fmt.Println()
	cf.printer.Info(true, "Iniciando processo de corte do vídeo...")

	ctx, stop := utils.SignalContext()
	defer stop()
	return cf.cutter.Cut(ctx, cf.params)
}

func (cf *CutFlow) getTimestamps() {
	for {
		cf.printer.Info(false, "Digite o timestamp de início do corte (formato deve ser hh:mm:ss, por exemplo 02:14:53): ")
		timestamp := utils.MustGetUserInput(cf.reader)

		err := utils.ValidateTimestampFormat(timestamp)
		if err != nil {
			cf.printer.Error(true, err.Error())
			continue
		}

		cf.params.From = timestamp
		break
	}

	for {
		cf.printer.Info(false, "Digite o timestamp de fim do corte (formato deve ser hh:mm:ss, por exemplo 02:14:53): ")
		timestamp := utils.MustGetUserInput(cf.reader)

		err := utils.ValidateTimestampFormat(timestamp)
		if err != nil {
			cf.printer.Error(true, err.Error())
			continue
		}

		cf.params.To = timestamp
		break
	}
}

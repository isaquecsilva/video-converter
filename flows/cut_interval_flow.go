package flows

import (
	"bufio"
	"fmt"
	"path/filepath"
	"strings"
	"video-converter/services/metadata"
	"video-converter/services/vcut"
	"video-converter/utils"
	"video-converter/utils/printer"
)

type CutIntervalFlow struct {
	reader   *bufio.Reader
	metadata *metadata.VideoMetadataExtractor
	printer  printer.Printer
	baseFlow BaseFlow
	cutter   vcut.VideoCutterInterface
	params   vcut.CutInIntervalParams
}

func NewCutIntervalFlow(
	reader *bufio.Reader,
	m *metadata.VideoMetadataExtractor,
	baseFlow BaseFlow,
	printer printer.Printer,
	cutter vcut.VideoCutterInterface,
) Flow {
	return &CutIntervalFlow{
		reader:   reader,
		metadata: m,
		printer:  printer,
		baseFlow: baseFlow,
		cutter:   cutter,
		params:   vcut.CutInIntervalParams{},
	}
}

func (cif *CutIntervalFlow) Run() (string, error) {
	cif.params.InputFile = cif.baseFlow.getInputFile()

	metadata, err := cif.metadata.ExtractMetadata(cif.params.InputFile)
	if err != nil {
		return "", err
	}

	cif.params.TotalSeconds = int64(metadata.Duration.Seconds())

	cif.params.Format = cif.baseFlow.getTargetFormat()

	codec, err := cif.baseFlow.getCodec(cif.params.Format)
	if err != nil {
		return "", err
	}

	cif.params.Codec = codec
	cif.params.Bitrate = cif.baseFlow.getBitrate(metadata.VideoBitrateInKbps)

	cif.getInterval()

	ctx, stop := utils.SignalContext()
	defer stop()

	files, err := cif.cutter.BatchCut(ctx, cif.params)
	if err != nil {
		return "", err
	}

	fmt.Println("\n")
	cif.printer.Info(true, "Cortes gerados:")
	cif.printer.Success(true, strings.Join(files, "\n"))
	fmt.Println()
	return filepath.Dir(cif.params.InputFile), nil
}

func (cif *CutIntervalFlow) getInterval() {
	for {
		cif.printer.Info(false, "Digite qual será o tempo de cada corte em hh:mm:ss (por exemplo, 00:05:00, cada corte terá 5 min): ")
		timestamp := utils.MustGetUserInput(cif.reader)

		d, err := utils.TimestampParser(timestamp)
		if err != nil {
			cif.printer.Error(true, err.Error())
			continue
		}

		if int64(d.Seconds()) <= 0 {
			cif.printer.Error(true, "Intervalo de corte deve ser maior que 0 segundos...")
			continue
		}

		cif.params.Interval = int64(d.Seconds())
		break
	}
}

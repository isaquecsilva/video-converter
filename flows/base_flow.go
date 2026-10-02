package flows

import (
	"bufio"
	"errors"
	"os"
	"slices"
	"strconv"
	"strings"
	"video-converter/services/discover"
	"video-converter/types"
	"video-converter/types/codecs"
	"video-converter/types/videoformat"
	"video-converter/utils"
	"video-converter/utils/printer"
)

type BaseFlow interface {
	getInputFile() string
	getTargetFormat() videoformat.VideoFormat
	getCodec(format videoformat.VideoFormat) (codecs.VideoCodec, error)
	getBitrate(defaultBitrate int) int
}

type baseFlow struct {
	formatKit videoformat.VideoFormatKit
	discover  *discover.GpuDiscover
	printer   printer.Printer
	stdin     *bufio.Reader
}

func NewBaseFlow(
	formatKit videoformat.VideoFormatKit,
	disc *discover.GpuDiscover,
	printer printer.Printer,
	stdin *bufio.Reader,
) BaseFlow {
	return &baseFlow{
		formatKit: formatKit,
		discover:  disc,
		printer:   printer,
		stdin:     stdin,
	}
}

func (base *baseFlow) getInputFile() string {
	for {
		base.printer.Info(false, "Cole o caminho para o arquivo de vídeo: ")
		filePath := utils.MustGetUserInput(base.stdin)

		if filePath == "" {
			base.printer.Error(true, "Você deve especificar o caminho para o arquivo de vídeo!")
			continue
		}

		filePath = strings.Trim(filePath, "\r\n\"' ")

		_, err := os.Stat(filePath)
		if errors.Is(err, os.ErrNotExist) {
			base.printer.Error(true, "O arquivo especificado não foi encontrado!")
			continue
		}

		return filePath
	}
}

func (base *baseFlow) getTargetFormat() videoformat.VideoFormat {
	supportedVideoFormats := slices.Collect(func(yield func(string) bool) {
		for _, format := range base.formatKit.SupportedFormats() {
			if !yield(format.Name()) {
				return
			}
		}
	})

	for {
		base.printer.Info(
			false,
			"Por favor, escolha o formato do vídeo de saída (%s | se vazio, será usado %s): ",
			strings.Join(supportedVideoFormats, ", "),
			base.formatKit.DefaultFormat().Name(),
		)
		choice := utils.MustGetUserInput(base.stdin)

		if choice == "" {
			return base.formatKit.DefaultFormat()
		}

		format, err := base.formatKit.SelectFormat(choice)
		if err != nil {
			base.printer.Error(true, err.Error())
			continue
		}

		return format
	}
}

func (base *baseFlow) getCodec(format videoformat.VideoFormat) (codecs.VideoCodec, error) {
	supportedCodecs := slices.Collect(func(yield func(string) bool) {
		for _, codec := range format.SupportedCodecs() {
			if !yield(codec.Name()) {
				return
			}
		}
	})

	for {
		base.printer.Info(
			false,
			"Escolha o codec (%s | se vazio, será usado %s): ",
			strings.Join(supportedCodecs, ", "),
			format.DefaultCodecKit().Name(),
		)

		choice := utils.MustGetUserInput(base.stdin)

		var (
			codecKit codecs.VideoCodecKit
			err      error
		)

		if choice == "" {
			codecKit = format.DefaultCodecKit()
		} else {
			codecKit, err = format.GetCodecKit(choice)
			if err != nil {
				base.printer.Error(true, err.Error())
				continue
			}
		}

		gpuType, err := base.gpuToUse()
		if err != nil {
			base.printer.Error(true, err.Error())
		}

		codec, err := codecKit.String(gpuType)
		if errors.Is(err, codecs.ErrNoCodecForGpu) {
		use_cpu:
			base.printer.Warn(false, "O codec %s não possui suporte para uso com GPU. Deseja converter com a CPU (y/n)? ", codecKit.Name())
			choice = utils.MustGetUserInput(base.stdin)

			if choice != "y" && choice != "n" {
				base.printer.Error(true, "Opção inválida!")
				goto use_cpu
			}

			if choice == "n" {
				continue
			}

			return codecKit.String(types.GpuTypeNone)
		}

		if err != nil {
			return "", err
		}

		return codec, nil
	}
}

func (base *baseFlow) gpuToUse() (types.GpuType, error) {
	info, err := base.discover.Discover()

	if errors.Is(err, discover.ErrNoDedicatedGPUFound) {
		base.printer.Warn(true, "Nenhum GPU dedicada encontrada no ambiente! Conversão será feita via CPU...")
		return types.GpuTypeNone, nil
	}

	if err != nil {
		return types.GpuTypeNone, err
	}

	base.printer.Success(true, "GPU dedicada encontra em seu ambiente: ")
	base.printer.Success(true, "Nome: %s", info.Name)
	base.printer.Success(true, "Total de RAM: %.1fMB", float64(info.AdapterRAM/1024/1024))

	base.printer.Info(false, "Deseja converter com GPU s/n (se vazio, será sim)? ")
	choice := strings.ToLower(utils.MustGetUserInput(base.stdin))

	if choice == "s" {
		return info.Type, nil
	}

	return types.GpuTypeNone, nil
}

func (base *baseFlow) getBitrate(defaultBitrate int) int {
	if defaultBitrate <= 0 {
		base.printer.Warn(true, "Bitrate padrão é muito baixo (%dkbps). Desta forma, caso um valor não seja específicado, seram usados 1000kbps...")
		defaultBitrate = 1000
	}

	for {
		base.printer.Info(false, "Especifique o bitrate do vídeo de saída (se vazio, será usado o padrão do vídeo %dkbps): ", defaultBitrate)

		input := utils.MustGetUserInput(base.stdin)

		if input == "" {
			return defaultBitrate
		}

		bitrate, err := strconv.Atoi(input)
		if err != nil {
			base.printer.Info(true, "Sua entrada deve ser apenas númerica. Por exemplo: 3000. Isso equivale à 3000kbps")
			continue
		}

		return bitrate
	}
}
